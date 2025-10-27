package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name        string
	Environment string
	Version     string
	Port        int
	HTTPS       HTTPSConfig
}

type HTTPSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
	Redirect bool
}

type MongoConfig struct {
	URI      string
	Database string
}

type CognitoConfig struct {
	UserPoolID string
	ClientID   string
	Domain     string
	Region     string
}

type S3Config struct {
	Bucket string
}

type SQSConfig struct {
	QueueName string
	QueueURL  string
}

type AWSConfig struct {
	Region            string
	AccessKeyID       string
	SecretAccessKey   string
	SessionToken      string
	Endpoint          string
	UseLocalstack     bool
	S3                S3Config
	SQS               SQSConfig
	Cognito           CognitoConfig
	CredentialsSource string
}

type AuthConfig struct {
	Mode string
}

type SecurityConfig struct {
	AllowedOrigins []string
	EncryptionKey  string
}

type StorageConfig struct {
	ReceiptBucket string
}

type QueueConfig struct {
	TransactionQueue string
}

type LocalConfig struct {
	CredentialsFile string
	AuthUsers       []LocalAuthUser
}

type LocalAuthUser struct {
	Username        string
	Password        string
	Email           string
	Name            string
	DefaultCurrency string
	CognitoSub      string
}

type Config struct {
	App      AppConfig
	Mongo    MongoConfig
	AWS      AWSConfig
	Auth     AuthConfig
	Security SecurityConfig
	Queue    QueueConfig
	Storage  StorageConfig
	Local    LocalConfig
}

var (
	cfg  *Config
	once sync.Once
)

// função interna para reiniciar cache de configuração, utilizada somente em testes
func reset() {
	once = sync.Once{}
	cfg = nil
}

func LoadConfig() (*Config, error) {
	var loadErr error
	once.Do(func() {
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		setDefaults()

		file := os.Getenv("CONFIG_FILE")
		if file == "" {
			defaultFile := "config/local_credentials.yaml"
			if _, err := os.Stat(defaultFile); err == nil {
				file = defaultFile
			}
		}
		if file != "" {
			viper.SetConfigFile(file)
			if err := viper.MergeInConfig(); err != nil {
				loadErr = fmt.Errorf("failed to load config file %s: %w", file, err)
				return
			}
		}

		tmp := &Config{
			App: AppConfig{
				Name:        viper.GetString("app.name"),
				Environment: viper.GetString("app.environment"),
				Version:     viper.GetString("app.version"),
				Port:        viper.GetInt("app.port"),
				HTTPS: HTTPSConfig{
					Enabled:  viper.GetBool("app.https.enabled"),
					CertFile: viper.GetString("app.https.certFile"),
					KeyFile:  viper.GetString("app.https.keyFile"),
					Redirect: viper.GetBool("app.https.redirect"),
				},
			},
			Mongo: MongoConfig{
				URI:      viper.GetString("mongo.uri"),
				Database: viper.GetString("mongo.database"),
			},
			AWS: AWSConfig{
				Region:            viper.GetString("aws.region"),
				AccessKeyID:       viper.GetString("aws.accessKeyId"),
				SecretAccessKey:   viper.GetString("aws.secretAccessKey"),
				SessionToken:      viper.GetString("aws.sessionToken"),
				Endpoint:          viper.GetString("aws.endpoint"),
				UseLocalstack:     viper.GetBool("aws.useLocalstack"),
				CredentialsSource: viper.GetString("aws.credentialsSource"),
				S3: S3Config{
					Bucket: viper.GetString("aws.s3.bucket"),
				},
				SQS: SQSConfig{
					QueueName: viper.GetString("aws.sqs.queueName"),
					QueueURL:  viper.GetString("aws.sqs.queueUrl"),
				},
				Cognito: CognitoConfig{
					UserPoolID: viper.GetString("aws.cognito.userPoolId"),
					ClientID:   viper.GetString("aws.cognito.clientId"),
					Domain:     viper.GetString("aws.cognito.domain"),
					Region:     viper.GetString("aws.cognito.region"),
				},
			},
			Auth: AuthConfig{
				Mode: viper.GetString("auth.mode"),
			},
			Security: SecurityConfig{
				AllowedOrigins: viper.GetStringSlice("security.allowedOrigins"),
				EncryptionKey:  viper.GetString("security.encryptionKey"),
			},
			Queue: QueueConfig{
				TransactionQueue: viper.GetString("queue.transactionQueue"),
			},
			Storage: StorageConfig{
				ReceiptBucket: viper.GetString("storage.receiptBucket"),
			},
			Local: LocalConfig{
				CredentialsFile: viper.GetString("local.credentialsFile"),
				AuthUsers:       readLocalAuthUsers(viper.Get("local.authUsers")),
			},
		}

		cfg = tmp
	})

	if loadErr != nil {
		return nil, loadErr
	}
	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("app.name", "financial-control-api")
	viper.SetDefault("app.version", "0.1.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.https.enabled", false)
	viper.SetDefault("app.https.certFile", "")
	viper.SetDefault("app.https.keyFile", "")
	viper.SetDefault("app.https.redirect", true)

	viper.SetDefault("mongo.uri", "mongodb://mongo:27017")
	viper.SetDefault("mongo.database", "financial-control")

	viper.SetDefault("aws.region", "us-east-1")
	viper.SetDefault("aws.useLocalstack", true)
	viper.SetDefault("aws.endpoint", "http://localstack:4566")
	viper.SetDefault("aws.credentialsSource", "env")
	viper.SetDefault("aws.s3.bucket", "financial-control-receipts")
	viper.SetDefault("aws.sqs.queueName", "financial-transactions-queue")
	viper.SetDefault("aws.cognito.region", "us-east-1")
	viper.SetDefault("aws.cognito.domain", "http://localhost:4566")

	viper.SetDefault("auth.mode", "cognito")

	viper.SetDefault("security.allowedOrigins", []string{"*"})
	viper.SetDefault("queue.transactionQueue", "financial-transactions-queue")
	viper.SetDefault("storage.receiptBucket", "financial-control-receipts")
	viper.SetDefault("local.credentialsFile", "config/local_credentials.yaml")
}

func readLocalAuthUsers(value any) []LocalAuthUser {
	users := []LocalAuthUser{}
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			if m, ok := item.(map[string]any); ok {
				users = append(users, LocalAuthUser{
					Username:        getString(m, "username"),
					Password:        getString(m, "password"),
					Email:           getString(m, "email"),
					Name:            getString(m, "name"),
					DefaultCurrency: getString(m, "defaultCurrency"),
					CognitoSub:      getString(m, "cognitoSub"),
				})
			}
		}
	}
	return users
}

func getString(m map[string]any, key string) string {
	if value, ok := m[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}
