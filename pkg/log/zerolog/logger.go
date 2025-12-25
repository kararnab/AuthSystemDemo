package zerolog

import (
	"context"
	stdlog "log"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	appLog "github.com/kararnab/authdemo/pkg/log"
)

type Logger struct {
	l zerolog.Logger
}

func New() *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// zerolog.TimeFieldFormat = time.RFC3339
	//log.Logger = log.Output(zerolog.ConsoleWriter{
	//	Out:        os.Stdout,
	//	TimeFormat: time.RFC3339,
	//})
	// log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger() // In Prod
	base := zlog.Output(zerolog.ConsoleWriter{Out: stdlog.Writer()})
	return &Logger{l: base}
}

func (z *Logger) Debug(msg string, fields ...appLog.Field) {
	z.write(z.l.Debug(), msg, fields)
}

func (z *Logger) Info(msg string, fields ...appLog.Field) {
	z.write(z.l.Info(), msg, fields)
}

func (z *Logger) Warn(msg string, fields ...appLog.Field) {
	z.write(z.l.Warn(), msg, fields)
}

func (z *Logger) Error(msg string, fields ...appLog.Field) {
	z.write(z.l.Error(), msg, fields)
}

func (z *Logger) With(fields ...appLog.Field) appLog.Logger {
	ctx := z.l.With()
	for k, v := range appLog.SanitizeFields(fields) {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{l: ctx.Logger()}
}

func (z *Logger) WithContext(ctx context.Context) appLog.Logger {
	return &Logger{l: z.l.With().Ctx(ctx).Logger()}
}

func (z *Logger) Unsafe() appLog.UnsafeLogger {
	return unsafeLogger{l: z.l}
}

func (z *Logger) write(e *zerolog.Event, msg string, fields []appLog.Field) {
	for k, v := range appLog.SanitizeFields(fields) {
		e.Interface(k, v)
	}
	e.Msg(msg)
}

// ---------- unsafe ----------

type unsafeLogger struct {
	l zerolog.Logger
}

func (u unsafeLogger) Debug(msg string, kv ...any) {
	u.l.Debug().Fields(kv).Msg(msg)
}
func (u unsafeLogger) Info(msg string, kv ...any) {
	u.l.Info().Fields(kv).Msg(msg)
}
func (u unsafeLogger) Warn(msg string, kv ...any) {
	u.l.Warn().Fields(kv).Msg(msg)
}
func (u unsafeLogger) Error(msg string, kv ...any) {
	u.l.Error().Fields(kv).Msg(msg)
}
