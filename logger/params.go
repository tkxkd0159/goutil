package logger

import (
	"time"
)

const (
	LowPriorityLogFileName  = "info.log"
	HighPriorityLogFileName = "error.log"

	DefaultLogPathPerm = 0o766

	DefaultLogTick          = time.Second
	DefaultLogLogFirst      = 10
	DefaultLogLogThereafter = 5
)
