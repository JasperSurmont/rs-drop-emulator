package beasts

import (
	"rs-drop-emulator/runescape/util"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger = zap.L().Sugar()

const (
	MAX_AMOUNT_ROLLS = 200
)

// ConfigLogger gets called when the global logger has been set up correctly
func ConfigLogger() {
	log = zap.L().Sugar().Named("beasts")
}

type namedRSPrice struct {
	name  string
	price util.RSPrice
}
