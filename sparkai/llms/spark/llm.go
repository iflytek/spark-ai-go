package spark

import (
	"errors"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark/internal/sparkclient"
	"net/http"
	"os"
)

var (
	ErrEmptyResponse            = errors.New("no response")
	ErrMissingAPPID             = errors.New("missing the Spark APP id, set it in the SPARK_APP_ID environment variable")         //nolint:lll
	ErrMissingAPIKey            = errors.New("missing the Spark API key, set it in the SPARK_API_KEY environment variable")       //nolint:lll
	ErrMissingAPISecret         = errors.New("missing the Spark API secret, set it in the SPARK_API_SECRET environment variable") //nolint:lll
	ErrMissingAPI               = errors.New("missing the SPARK_BASE_URL set it in the SPARK_BASE_URL environment variable")      //nolint:lll
	ErrUnexpectedResponseLength = errors.New("unexpected length of response")
	DefaultSparkUrl             = "wss://spark-api.xf-yun.com/v3.1/multimodal"
)

// newClient is wrapper for sparkclient internal package.
func newClient(opts ...Option) (*options, *sparkclient.Client, error) {
	sparkUrl := os.Getenv(baseURLEnvVarName)
	if sparkUrl == "" {
		sparkUrl = DefaultSparkUrl
	}
	options := &options{
		apiKey:       os.Getenv(apiKeyEnvVarName),
		apiSecret:    os.Getenv(apiSecretEnvVarName),
		appId:        os.Getenv(appIdEnvVarName),
		baseURL:      sparkUrl,
		organization: os.Getenv(organizationEnvVarName),
		httpClient:   http.DefaultClient,
	}

	for _, opt := range opts {
		opt(options)
	}
	if len(options.baseURL) == 0 {
		return options, nil, ErrMissingAPI
	}
	if len(options.appId) == 0 {
		return options, nil, ErrMissingAPPID
	}
	if len(options.apiSecret) == 0 {
		return options, nil, ErrMissingAPISecret
	}
	if len(options.apiKey) == 0 {
		return options, nil, ErrMissingAPIKey
	}
	cli, err := sparkclient.New(options.apiKey, options.apiSecret, options.appId, options.baseURL, options.organization,
		options.apiVersion, options.embeddingModel)
	return options, cli, err
}

func getEnvs(keys ...string) string {
	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		if ok {
			return val
		}
	}
	return ""
}
