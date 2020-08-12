package log

import (
	"os"
	"time"
	"strings"
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	m "github.com/WeBankPartners/wecube-plugins-saltstack/common/models"
)

var (
	Logger *zap.Logger
	levelStringList = []string{"debug","info","warn","error"}
	LogLevel string
)

func InitZapLogger() {
	LogLevel = strings.ToLower(m.Config.Log.Level)
	var level int
	for i,v := range levelStringList {
		if v == LogLevel {
			level = i-1
			break
		}
	}
	zapLevel := zap.NewAtomicLevel()
	zapLevel.SetLevel(zapcore.Level(level))
	hook := lumberjack.Logger{
		Filename:   m.Config.Log.File,
		MaxSize:    m.Config.Log.ArchiveMaxSize,
		MaxBackups: m.Config.Log.ArchiveMaxBackup,
		MaxAge:     m.Config.Log.ArchiveMaxDay,
		Compress:   m.Config.Log.Compress,
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),zapcore.AddSync(&hook)), zapLevel)
	Logger = zap.New(core, zap.AddCaller(), zap.Development())
	Logger.Info("Success init zap log !!")
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func String(k,v string) zap.Field {
	return zap.String(k, v)
}

func Int(k string,v int) zap.Field {
	return zap.Int(k, v)
}

func Int64(k string,v int64) zap.Field {
	return zap.Int64(k, v)
}

func Float64(k string,v float64) zap.Field {
	return zap.Float64(k, v)
}

func JsonObj(k string,v interface{}) zap.Field {
	b,err := json.Marshal(v)
	if err == nil {
		return zap.String(k, string(b))
	}else{
		return zap.Error(err)
	}
}

func StringList(k string,v []string) zap.Field {
	return zap.Strings(k, v)
}