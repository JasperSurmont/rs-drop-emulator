package beasts

import "go.uber.org/zap"

var log *zap.SugaredLogger

const (
	MAX_AMOUNT_ROLLS = 200
)

func ConfigLogger(l *zap.SugaredLogger) {
	log = l.Desugar().Sugar().Named("beasts")
}
