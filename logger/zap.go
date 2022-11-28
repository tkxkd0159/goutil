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

var l *zap.Logger
var s *zap.SugaredLogger

type ProdLoggerOpts struct {
	Sample   *LogSamplerOpts
	Logpath  string
	Pathperm os.FileMode
}

type LogSamplerOpts struct {
	Tick       time.Duration
	First      int // samples by logging the first N entries with a given level and message each tick
	Thereafter int // If more Entries with the same level and message are seen during the same interval, every Mth message is logged and the rest are dropped.
}

// SetLogger makes zap-based logger
// Use filepath.Join to make logpath
func SetLogger(m Mode, prodOpts ...ProdLoggerOpts) {
	switch m {
	case PROD:
		if prodOpts == nil {
			panic("you should set logger options for production mode")
		}
		opts := prodOpts[0]
		encCfg := zap.NewProductionEncoderConfig()
		enc := zapcore.NewJSONEncoder(encCfg)

		if opts.Pathperm == 0 {
			MustMkDir(opts.Logpath, DefaultLogPathPerm)
		} else {
			MustMkDir(opts.Logpath, opts.Pathperm)
		}

		lowsink, _, err := zap.Open(filepath.Join(opts.Logpath, LowPriorityLogFileName))
		if err != nil {
			panic(err)
		}
		highsink, _, err := zap.Open(filepath.Join(opts.Logpath, HighPriorityLogFileName))
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
		if opts.Sample == nil {
			core = zapcore.NewSamplerWithOptions(core, DefaultLogTick, DefaultLogLogFirst, DefaultLogLogThereafter)
		} else {
			core = zapcore.NewSamplerWithOptions(core, opts.Sample.Tick, opts.Sample.First, opts.Sample.Thereafter)
		}
		zapopts := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)}
		logger := zap.New(core, zapopts...)

		l = logger
		s = logger.Sugar()
	case DEBUG:
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		l = logger
		s = logger.Sugar()
	case TEST:
		logger := zap.NewNop()
		l = logger
		s = logger.Sugar()
	}
}

func L() *zap.Logger {
	if l == nil {
		panic("you must set logger before run the fetcher")
	} else {
		return l
	}
}

func S() *zap.SugaredLogger {
	if s == nil {
		panic("you must set logger before run the fetcher")
	} else {
		return s
	}
}
