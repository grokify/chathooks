package adapters

import (
	"github.com/grokify/commonchat/glip/config"
)

func GlipConfig() *config.ConverterConfig {
	return &config.ConverterConfig{
		UseAttachments:        true,
		UseFieldExtraSpacing:  true,
		ConvertTripleBacktick: true,
	}
}
