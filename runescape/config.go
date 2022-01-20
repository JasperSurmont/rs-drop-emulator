package runescape

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger = zap.L().Sugar()

// ConfigLogger gets called when the global logger has been set up correctly
func ConfigLogger() {
	log = zap.L().Sugar().Named("runescape")
}
