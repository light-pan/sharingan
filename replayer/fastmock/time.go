package fastmock

import (
	"time"

	"github.com/light-pan/sharingan/replayer/internal"
	"github.com/light-pan/sharingan/replayer/monkey"
)

// MockTime mock Time
func MockTime() {
	mockTimeNow()
}

// mock Time.Now()
func mockTimeNow() {
	monkey.MockGlobalFunc(time.Now, func() time.Time {
		threadID := internal.GetCurrentGoRoutineID()
		replayTime := int64(0)

		if thread := ReplayerGlobalThreads.Get(threadID); thread != nil {
			replayTime = thread.replayTime
		}

		if replayTime == 0 {
			return internal.TimeNow()
		}

		sec := replayTime / 1000000000
		nsec := replayTime % 1000000000
		return time.Unix(sec, nsec)
	})
}
