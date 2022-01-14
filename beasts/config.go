package beasts

import "go.uber.org/zap"

var log *zap.SugaredLogger

func ConfigLogger(l *zap.SugaredLogger) {
	log = l
}
