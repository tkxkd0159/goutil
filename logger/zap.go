package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type clientCtx interface {
	GetHomeDir() string
}

type Mode int

const (
	DEBUG Mode = iota
	PROD
)

var s *zap.SugaredLogger

type LogSamplerOpts struct {
	Tick       time.Duration
	First      int // samples by logging the first N entries with a given level and message each tick
	Thereafter int // If more Entries with the same level and message are seen during the same interval, every Mth message is logged and the rest are dropped.
}

// SetLogger makes zap-based logger
// Use filepath.Join to make logpath
func SetLogger(m Mode, opts *LogSamplerOpts, logpath string, pathperm ...os.FileMode) {
	switch m {
	case PROD:
		encCfg := zap.NewProductionEncoderConfig()
		enc := zapcore.NewJSONEncoder(encCfg)

		if pathperm == nil {
			MustMkDir(logpath, DefaultLogPathPerm)
		} else {
			MustMkDir(logpath, pathperm[0])
		}

		lowsink, _, err := zap.Open(filepath.Join(logpath, LowPriorityLogFileName))
		if err != nil {
			panic(err)
		}
		highsink, _, err := zap.Open(filepath.Join(logpath, HighPriorityLogFileName))
		if err != nil {
			panic(err)
		}
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})

		core := zapcore.NewTee(
			zapcore.NewCore(enc, lowsink, lowPriority),
			zapcore.NewCore(enc, highsink, highPriority),
		)
		if opts == nil {
			core = zapcore.NewSamplerWithOptions(core, DefaultLogTick, DefaultLogLogFirst, DefaultLogLogThereafter)
		} else {
			core = zapcore.NewSamplerWithOptions(core, opts.Tick, opts.First, opts.Thereafter)
		}
		zapopts := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)}
		logger := zap.New(core, zapopts...)
		s = logger.Sugar()
	case DEBUG:
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		s = logger.Sugar()
	}
}

func S() *zap.SugaredLogger {
	if s == nil {
		panic("you must set logger before run the fetcher")
	} else {
		return s
	}
}
