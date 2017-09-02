package logmw

import (
	"io/ioutil"
	"path"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	// SystemField is used in every log statement made through grpc_logrus. Can be overwritten before any initialization code.
	SystemField = "system"

	// KindField describes the log gield used to incicate whether this is a server or a client log statment.
	KindField = "span.kind"
)

// UnaryServerInterceptor returns a new unary server interceptors that sets the values for request tags.
func UnaryServerInterceptor(entry *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// TODO: extract request fields later...

		newCtx := newLoggerForCall(ctx, entry, info.FullMethod)
		startTime := time.Now()

		out, err := handler(newCtx, req)

		// TODO: extract response fields

		code := grpc_logging.DefaultErrorToCode(err)
		level := grpc_logrus.DefaultCodeToLevel(code)
		durField, durVal := grpc_logrus.DefaultDurationToField(time.Now().Sub(startTime))
		fields := logrus.Fields{
			"grpc.code": code.String(),
			durField:    durVal,
		}
		if err != nil {
			fields[logrus.ErrorKey] = err
		}
		levelLogf(
			ExtractEntry(newCtx).WithFields(fields), // re-extract logger from newCtx, as it may have extra fields that changed in the holder.
			level,
			"finished unary call")
		return out, err
	}
}

// TODO: Stream server interceptor

func levelLogf(entry *logrus.Entry, level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.DebugLevel:
		entry.Debugf(format, args...)
	case logrus.InfoLevel:
		entry.Infof(format, args...)
	case logrus.WarnLevel:
		entry.Warningf(format, args...)
	case logrus.ErrorLevel:
		entry.Errorf(format, args...)
	case logrus.FatalLevel:
		entry.Fatalf(format, args...)
	case logrus.PanicLevel:
		entry.Panicf(format, args...)
	}
}

func newLoggerForCall(ctx context.Context, entry *logrus.Entry, fullMethodString string) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	callLog := entry.WithFields(
		logrus.Fields{
			SystemField:    "grpc",
			KindField:      "server",
			"grpc.service": service,
			"grpc.method":  method,
		})
	return EntryToContext(ctx, callLog)
}

var (
	nullLogger = &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.PanicLevel,
	}
)
