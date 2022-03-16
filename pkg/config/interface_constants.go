package config

import "strings"

const (
	ParamNameAdapters        = "adapters"
	ParamNameActivityDefault = "defaultActivity"
	ParamNameIconDefault     = "defaultIcon"
	ParamNameInputType       = "inputType"
	ParamNameOutputType      = "outputType"
	ParamNameOutputFormat    = "outputFormat" // `card`, `adaptivecard`, `nocard`. Default is `card`.
	ParamNameOutputURL       = "outputURL"
	ParamNameToken           = "token"
	EnvPath                  = "ENV_PATH"
	EnvEngine                = "CHATHOOKS_ENGINE" // awslambda, nethttp, fasthttp
	EnvTokens                = "CHATHOOKS_TOKENS"
	EnvWebhookURL            = "CHATHOOKS_URL"
	EnvHomeURL               = "CHATHOOKS_HOME_URL"
	ErrRequiredTokenNotFound = "401.01 Required Token Not Found"
	ErrRequiredTokenNotValid = "401.02 Required Token Not Valid"
	// ParamNameURL             = "url" // legacy. deprecated.

	ParamNameOutputFormatNocard       = "nocard"
	ParamNameOutputFormatCard         = "card"
	ParamNameOutputFormatAdaptivecard = "adaptivecard"
)

func MustParseOutputFormat(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	switch input {
	case ParamNameOutputFormatNocard:
		return ParamNameOutputFormatNocard
	case ParamNameOutputFormatNocard + "s":
		return ParamNameOutputFormatNocard
	case ParamNameOutputFormatCard:
		return ParamNameOutputFormatCard
	case ParamNameOutputFormatCard + "s":
		return ParamNameOutputFormatCard
	case ParamNameOutputFormatAdaptivecard:
		return ParamNameOutputFormatAdaptivecard
	case ParamNameOutputFormatAdaptivecard + "s":
		return ParamNameOutputFormatAdaptivecard
	default:
		return ""
	}
}
