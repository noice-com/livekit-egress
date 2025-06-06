package logging

import (
	"fmt"
	"strings"

	"github.com/livekit/protocol/logger"
	"github.com/livekit/protocol/logger/medialogutils"
)

const (
	logError = iota
	logWarn
	logInfo
	ignore
)

var actions = map[string]int{
	"0:00:": ignore,
	"te_au": ignore,
	"turnc": logInfo,
	"ice E": logInfo,
	"SDK 2": logInfo,
	"(egre": logWarn,
}

func NewHandlerLogger(handlerID, egressID string) *medialogutils.CmdLogger {
	l := logger.GetLogger().WithValues("handlerID", handlerID, "egressID", egressID)
	return medialogutils.NewCmdLogger(func(s string) {
		// glib inserts 2 carriage returns at the end of its warning/error logs
		lines := strings.Split(strings.TrimLeft(strings.TrimRight(s, "\n"), "\n"), "\n")
		for _, line := range lines {
			if strings.HasSuffix(line, "}") {
				fmt.Println(line)
			} else {
				action := logError
				if len(line) > 5 {
					action = actions[line[:5]]
				}
				switch action {
				case ignore:
					continue
				case logInfo:
					l.Infow(line)
				case logWarn:
					l.Warnw(line, nil)
				case logError:
					l.Errorw(line, nil)
				}
			}
		}
	})
}
