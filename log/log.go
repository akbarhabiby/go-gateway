package log

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger

type Field struct {
	Key string
	Val interface{}
}

type TDRModel struct {
	AppName    string `json:"app"`
	AppVersion string `json:"ver"`
	AppPort    int    `json:"port"`
	ThreadID   string `json:"xid"`

	Path     string `json:"path"`
	Method   string `json:"method"`
	SrcIP    string `json:"srcIP"`
	RespTime int64  `json:"rt"`

	Header   interface{} `json:"header"`
	Request  interface{} `json:"req"`
	Response interface{} `json:"resp"`
	Error    string      `json:"error"`
}

const separator = "|"

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "xtime",
		MessageKey: "x",
		EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
			pae.AppendString(t.UTC().Format(time.RFC3339))
		},
		LineEnding: zapcore.DefaultLineEnding,
	}
	zapLogger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zap.InfoLevel))
}

func Info(ctx context.Context, title string, messages ...interface{}) {
	logs := []zap.Field{
		zap.String("level", "info"),
		zap.String("message", title),
	}

	logs = append(logs, formatLogs(ctx, messages...)...)
	zapLogger.Info(separator, logs...)
}

func Error(ctx context.Context, title string, messages ...interface{}) {
	logs := []zap.Field{
		zap.String("level", "error"),
		zap.String("message", title),
	}

	logs = append(logs, formatLogs(ctx, messages...)...)
	zapLogger.Error(separator, logs...)
}

func TDR(ctx context.Context, request, response []byte) {
	rt := time.Since(GetRequestTimeFromContext(ctx)).Nanoseconds() / 1000000
	ctxLogger := ExtractCtx(ctx)
	tdr := TDRModel{
		AppName:    ctxLogger.ServiceName,
		AppVersion: ctxLogger.ServiceVersion,
		AppPort:    ctxLogger.ServicePort,
		ThreadID:   ctxLogger.ThreadID,
		Path:       ctxLogger.ReqURI,
		Method:     ctxLogger.ReqMethod,
		SrcIP:      GetRequestIPFromContext(ctx),
		RespTime:   rt,
		Header:     GetRequestHeaderFromContext(ctx),
		Request:    string(request),
		Response:   string(response),
		Error:      GetErrorMessageFromContext(ctx),
	}

	fields := make([]zap.Field, 0)
	fields = append(fields, zap.String("level", "info"))
	fields = append(fields, zap.String("message", separator))

	fields = append(fields, formatLogs(ctx, separator)...)

	fields = append(fields, zap.String("app", tdr.AppName))
	fields = append(fields, zap.String("ver", tdr.AppVersion))
	fields = append(fields, zap.Int("port", tdr.AppPort))
	fields = append(fields, zap.String("xid", tdr.ThreadID))

	fields = append(fields, zap.Any("path", tdr.Path))
	fields = append(fields, zap.String("method", tdr.Method))
	fields = append(fields, zap.String("srcIP", tdr.SrcIP))
	fields = append(fields, zap.Int64("rt", tdr.RespTime))

	fields = append(fields, zap.Any("header", tdr.Header))
	fields = append(fields, zap.Any("req", tdr.Request))
	fields = append(fields, zap.Any("resp", tdr.Response))
	fields = append(fields, zap.String("error", tdr.Error))

	zapLogger.Info(separator, fields...)
}

func formatLogs(ctx context.Context, messages ...interface{}) (fields []zap.Field) {
	for index, msg := range messages {
		fields = append(fields, zap.Any("_message_"+cast.ToString(index), msg))
	}
	return
}
