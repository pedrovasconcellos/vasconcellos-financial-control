package auth

import (
	"context"
	"fmt"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	appErrors "github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/domain/port"
)

type CognitoAuthProvider struct {
	client   *cognitoidentityprovider.Client
	clientID string
}

var _ port.AuthProvider = (*CognitoAuthProvider)(nil)

func NewCognitoAuthProvider(ctx context.Context, region string, customEndpoint string, clientID string, accessKey string, secretKey string, sessionToken string) (*CognitoAuthProvider, error) {
	loadOptions := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(region),
	}
	if customEndpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: customEndpoint, SigningRegion: region}, nil
		})
		loadOptions = append(loadOptions, awsConfig.WithEndpointResolverWithOptions(resolver))
	}
	if accessKey != "" && secretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)
		loadOptions = append(loadOptions, awsConfig.WithCredentialsProvider(creds))
	}

	cfg, err := awsConfig.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, err
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)
	return &CognitoAuthProvider{client: client, clientID: clientID}, nil
}

func (p *CognitoAuthProvider) Login(ctx context.Context, credentials port.AuthCredentials) (*port.AuthTokens, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(p.clientID),
		AuthParameters: map[string]string{
			"USERNAME": credentials.Username,
			"PASSWORD": credentials.Password,
		},
	}

	output, err := p.client.InitiateAuth(ctx, input)
	if err != nil {
		if isCognitoUserError(err) {
			return nil, appErrors.ErrInvalidInput
		}
		return nil, fmt.Errorf("cognito auth failed: %w", err)
	}

	authResult := output.AuthenticationResult
	if authResult == nil {
		return nil, fmt.Errorf("authentication failed")
	}

	expiresIn := authResult.ExpiresIn
	return &port.AuthTokens{
		AccessToken:  aws.ToString(authResult.AccessToken),
		RefreshToken: aws.ToString(authResult.RefreshToken),
		IDToken:      aws.ToString(authResult.IdToken),
		ExpiresIn:    expiresIn,
		TokenType:    aws.ToString(authResult.TokenType),
	}, nil
}

func (p *CognitoAuthProvider) Validate(ctx context.Context, accessToken string) (map[string]any, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	output, err := p.client.GetUser(ctx, input)
	if err != nil {
		return nil, err
	}

	attributes := map[string]any{}
	for _, attr := range output.UserAttributes {
		attributes[aws.ToString(attr.Name)] = aws.ToString(attr.Value)
	}
	attributes["username"] = aws.ToString(output.Username)

	return attributes, nil
}
