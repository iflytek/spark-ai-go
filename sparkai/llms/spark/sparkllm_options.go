package spark

import (
	"github.com/iflytek/spark-ai-go/sparkai/callbacks"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark/internal/sparkclient"
)

const (
	AppIdEnvVarName        = "SPARKAI_APP_ID"     //nolint:gosec
	ApiKeyEnvVarName       = "SPARKAI_API_KEY"    //nolint:gosec
	ApiSecretEnvVarName    = "SPARKAI_API_SECRET" //nolint:gosec
	SparkDomainEnvVarName  = "SPARKAI_DOMAIN"
	sparkVersionEnvVarName = "SPARKAI_API_VERSION" //nolint:gosec
	BaseURLEnvVarName      = "SPARKAI_URL"         //nolint:gosec
	organizationEnvVarName = "SPARK_ORGANIZATION"  //nolint:gosec
)

const (
	DefaultAPIVersion = "v3.1"
)

type options struct {
	appId        string
	apiKey       string
	apiSecret    string
	domain       string
	baseURL      string
	organization string
	httpClient   sparkclient.Doer

	// required when APIType is APITypeAzure or APITypeAzureAD
	apiVersion     string
	embeddingModel string

	callbackHandler callbacks.Handler
}

type Option func(*options)

// WithToken passes the SPARK API token to the client. If not set, the token
// is read from the SPARK_API_KEY environment variable.
func WithApiKey(api_key string) Option {
	return func(opts *options) {
		opts.apiKey = api_key
	}
}

// WithToken passes the SPARK API token to the client. If not set, the token
// is read from the SPARK_API_KEY environment variable.
func WithApiSecret(app_sec string) Option {
	return func(opts *options) {
		opts.apiSecret = app_sec
	}
}

// WithToken passes the SPARK API token to the client. If not set, the token
// is read from the SPARK_API_KEY environment variable.
func WithAppId(app_id string) Option {
	return func(opts *options) {
		opts.appId = app_id
	}
}

// WithToken passes the SPARK API token to the client. If not set, the token
// is read from the SPARK_API_KEY environment variable.
func WithAPIDomain(domain string) Option {
	return func(opts *options) {
		opts.domain = domain
	}
}

// WithEmbeddingModel passes the SPARK model to the client. Required when ApiType is Azure.
func WithEmbeddingModel(embeddingModel string) Option {
	return func(opts *options) {
		opts.embeddingModel = embeddingModel
	}
}

// WithBaseURL passes the SPARK base url to the client. If not set, the base url
// is read from the SPARK_BASE_URL environment variable. If still not set in ENV
// VAR SPARK_BASE_URL, then the default value is https://api.SPARK.com/v1 is used.
func WithBaseURL(baseURL string) Option {
	return func(opts *options) {
		opts.baseURL = baseURL
	}
}

// WithOrganization passes the SPARK organization to the client. If not set, the
// organization is read from the SPARK_ORGANIZATION.
func WithOrganization(organization string) Option {
	return func(opts *options) {
		opts.organization = organization
	}
}

// WithAPIVersion passes the api version to the client. If not set, the default value
// is DefaultAPIVersion.
func WithAPIVersion(apiVersion string) Option {
	return func(opts *options) {
		opts.apiVersion = apiVersion
	}
}

// WithHTTPClient allows setting a custom HTTP client. If not set, the default value
// is http.DefaultClient.
func WithHTTPClient(client sparkclient.Doer) Option {
	return func(opts *options) {
		opts.httpClient = client
	}
}

//// WithCallback allows setting a custom Callback Handler.
//func WithCallback(callbackHandler callbacks.Handler) Option {
//	return func(opts *options) {
//		opts.callbackHandler = callbackHandler
//	}
//}
