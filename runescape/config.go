package runescape

import (
	"github.com/jaspersurmont/rs-drop-emulator/logger"
)

var log logger.LoggerWrapper

// ConfigLogger gets called when the global logger has been set up correctly
func init() {
	log = logger.CreateLogger("runescape")
}
