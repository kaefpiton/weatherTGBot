package zerolog

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
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

func (z *ZeroLogWrapper) Infof(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Info().Msg(msg)
}

func (z *ZeroLogWrapper) Warnf(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Warn().Msg(msg)
}

func (z *ZeroLogWrapper) Errorf(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Error().Msg(msg)
}

func (z *ZeroLogWrapper) Debugf(str string, kv ...interface{}) {
	msg := fmt.Sprintf(str, kv...)
	z.log.Debug().Msg(msg)
}
