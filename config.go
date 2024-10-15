package openai

import (
	"net/http"
	"regexp"
)

const (
	openaiAPIURLv1                 = "https://api.openai.com/v1"
	defaultEmptyMessagesLimit uint = 300

	azureAPIPrefix         = "openai"
	azureDeploymentsPrefix = "deployments"
)

type APIType string

const (
	APITypeOpenAI           APIType = "OPEN_AI"
	APITypeAzure            APIType = "AZURE"
	APITypeAzureAD          APIType = "AZURE_AD"
	APITypeCloudflareAzure  APIType = "CLOUDFLARE_AZURE"
	APITypeAbstractProvider APIType = "ABSTRACT_PROVIDER"
)

const AzureAPIKeyHeader = "api-key"

const defaultAssistantVersion = "v2" // upgrade to v2 to support vector store

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	Authtoken string

	BaseURL              string
	OrgID                string
	APIType              APIType
	APIVersion           string // required when APIType is APITypeAzure or APITypeAzureAD
	AssistantVersion     string
	AzureModelMapperFunc func(model string) string // replace model to azure deployment name func
	HTTPClient           HTTPDoer

	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		Authtoken:        authToken,
		BaseURL:          openaiAPIURLv1,
		APIType:          APITypeOpenAI,
		AssistantVersion: defaultAssistantVersion,
		OrgID:            "",

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func ConfigWithCustomHttpClient(baseUrl, authToken string, httpClient *http.Client) ClientConfig {
	return ClientConfig{
		Authtoken:        authToken,
		BaseURL:          baseUrl,
		APIType:          APITypeAbstractProvider,
		AssistantVersion: defaultAssistantVersion,
		OrgID:            "",

		HTTPClient: httpClient,

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		Authtoken:  apiKey,
		BaseURL:    baseURL,
		OrgID:      "",
		APIType:    APITypeAzure,
		APIVersion: "2023-05-15",
		AzureModelMapperFunc: func(model string) string {
			return regexp.MustCompile(`[.:]`).ReplaceAllString(model, "")
		},

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}

func (c ClientConfig) GetAzureDeploymentByModel(model string) string {
	if c.AzureModelMapperFunc != nil {
		return c.AzureModelMapperFunc(model)
	}

	return model
}
