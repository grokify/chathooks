package config

const (
	ParamNameAdapters        = "adapters"
	ParamNameActivityDefault = "defaultActivity"
	ParamNameIconDefault     = "defaultIcon"
	ParamNameInputType       = "inputType"
	ParamNameOutputType      = "outputType"
	ParamNameOutputURL       = "outputURL"
	ParamNameURL             = "url" // legacy. deprecated.
	ParamNameToken           = "token"
	EnvPath                  = "ENV_PATH"
	EnvEngine                = "CHATHOOKS_ENGINE" // awslambda, nethttp, fasthttp
	EnvTokens                = "CHATHOOKS_TOKENS"
	EnvWebhookUrl            = "CHATHOOKS_URL"
	EnvHomeUrl               = "CHATHOOKS_HOME_URL"
	ErrRequiredTokenNotFound = "401.01 Required Token Not Found"
	ErrRequiredTokenNotValid = "401.02 Required Token Not Valid"
)
