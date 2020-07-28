package shutdown

import (
	"sync"
)

var fShutdownRequested bool = false
var mShutdownMutex sync.Mutex

func RequestShutdown() {
	mShutdownMutex.Lock()
	fShutdownRequested = true
	mShutdownMutex.Unlock()
}

func GetShutdownStatus() bool {
	var ret bool

	mShutdownMutex.Lock()
	ret = fShutdownRequested
	mShutdownMutex.Unlock()

	return ret
}
