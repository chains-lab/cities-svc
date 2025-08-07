package logger

import (
	"context"
	"strings"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/interceptors"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewLogger(cfg config.Config) Logger {
	base := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(cfg.Logger.Level))
	if err != nil {
		base.Warnf("invalid log level '%s', defaulting to 'info'", cfg.Logger.Level)
		lvl = logrus.InfoLevel
	}
	base.SetLevel(lvl)

	switch strings.ToLower(cfg.Logger.Format) {
	case "json":
		base.SetFormatter(&logrus.JSONFormatter{})
	default:
		base.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	return NewWithBase(base)
}

func UnaryLogInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Вместо context.Background() используем входящий ctx,
		// чтобы не потерять таймауты и другую информацию.
		ctxWithLog := context.WithValue(
			ctx,
			interceptors.LogCtxKey,
			log, // ваш интерфейс Logger
		)

		// Далее передаём новый контекст в реальный хэндлер
		return handler(ctxWithLog, req)
	}
}

func Log(ctx context.Context, requestID uuid.UUID) Logger {
	entry, ok := ctx.Value(interceptors.LogCtxKey).(Logger)
	if !ok {
		logrus.Info("no logger in context")

		entry = NewWithBase(logrus.New())
	}
	return &logger{Entry: entry.WithField("request_id", requestID)}
}

// Logger — это ваш интерфейс: все методы FieldLogger + специальный WithError.
type Logger interface {
	WithError(err error) *logrus.Entry

	logrus.FieldLogger // сюда входят Debug, Info, WithField, WithError и т.д.
}

// logger — реальный тип, который реализует Logger.
type logger struct {
	*logrus.Entry // за счёт встраивания мы уже наследуем все методы FieldLogger
}

// WithError — ваш особый метод.
func (l *logger) WithError(err error) *logrus.Entry {
	ae := ape.Unwrap(err)
	if ae != nil {
		return l.Entry.WithError(ae.Unwrap())
	}

	return l.Entry.WithError(err)
}

func NewWithBase(base *logrus.Logger) Logger {
	log := logger{
		Entry: logrus.NewEntry(base),
	}

	return &log
}
