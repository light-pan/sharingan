// +build recorder

package sharingan

import (
	"log"
	"os"

	"github.com/light-pan/sharingan/plugins"
	"github.com/light-pan/sharingan/recorder"
	"github.com/light-pan/sharingan/recorder/koala/hook"
	"github.com/light-pan/sharingan/recorder/koala/logger"
	"github.com/light-pan/sharingan/recorder/koala/sut"
)

// GetCurrentGoRoutineID get current goroutineID incase SetDelegatedFromGoRoutineID
func GetCurrentGoRoutineID() int64 {
	return recorder.GetCurrentGoRoutineID()
}

// SetDelegatedFromGoRoutineID should be used when this goroutine is doing work for another goroutine
func SetDelegatedFromGoRoutineID(gID int64) {
	recorder.SetDelegatedFromGoRoutineID(gID)
}

func init() {
	if os.Getenv("RECORDER_ENABLED") != "true" {
		return
	}

	// init logger
	logger.Init()

	// init plugin && start recorder
	plugins.InitRecorderPlugin()
	plugins.StartRecorder()

	// start hook
	hook.Start()

	// start gc
	sut.StartGC()

	// log
	log.Println("mode", "=====recorder=====")
}
