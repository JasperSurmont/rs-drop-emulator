package runescape

import (
	"github.com/jaspersurmont/rs-drop-emulator/logger"
)

var log logger.LoggerWrapper

func init() {
	log = logger.CreateLogger("runescape")
}
