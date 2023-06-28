package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type ZeroLogWrapper struct {
	log zerolog.Logger
}

func NewZeroLog(logWriter io.Writer, logLevel string) *ZeroLogWrapper {
	ioWriter := zerolog.ConsoleWriter{
		Out:        logWriter,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
	}

	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	zeroLog := zerolog.New(ioWriter)
	zeroLog = zeroLog.Level(lvl).With().Timestamp().Logger()

	return &ZeroLogWrapper{
		log: zeroLog,
	}
}

func (z *ZeroLogWrapper) Warn(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Warn().Msg(msg)
}

func (z *ZeroLogWrapper) Error(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Error().Msg(msg)
}

func (z *ZeroLogWrapper) Debug(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Debug().Msg(msg)
}

func (z *ZeroLogWrapper) Info(kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Info().Msg(msg)
}

func (z *ZeroLogWrapper) Warnf(s string, kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Warn().Str(s, msg).Send()
}

func (z *ZeroLogWrapper) Errorf(s string, kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Error().Str(s, msg).Send()
}

func (z *ZeroLogWrapper) Debugf(s string, kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Debug().Str(s, msg).Send()
}

func (z *ZeroLogWrapper) Infof(s string, kv ...interface{}) {
	msg := fmt.Sprint(kv...)
	z.log.Info().Str(s, msg).Send()
}
