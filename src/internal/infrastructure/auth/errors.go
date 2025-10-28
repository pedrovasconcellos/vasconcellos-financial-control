package auth

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func isCognitoUserError(err error) bool {
	var authErr *types.NotAuthorizedException
	var userNotFoundErr *types.UserNotFoundException
	var passwordResetRequired *types.PasswordResetRequiredException
	var userNotConfirmed *types.UserNotConfirmedException

	return errors.As(err, &authErr) ||
		errors.As(err, &userNotFoundErr) ||
		errors.As(err, &passwordResetRequired) ||
		errors.As(err, &userNotConfirmed)
}

