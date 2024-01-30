package logs

import (
	//"marketplace/internal/config"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger Logger
	zapLogger     *zap.Logger
)

type LogConfig struct {
	Env        string `yaml:"env"`
	Path       string `yaml:"path"`
	Encoding   string `yaml:"encoding"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

func Init(logConf LogConfig) {

	// 日志输出文件还是控制台
	var writer zapcore.WriteSyncer
	if logConf.Path != "" {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), getLogFileWriter(logConf))
	} else {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		buildEncoder(logConf),
		writer,
		enableLevel(logConf),
	)

	// 增加堆栈打印 （caller跳过日志库工具类的代码行）
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	defaultLogger = zapLogger.Sugar()
}

func Sync() {
	_ = zapLogger.Sync()
}

func GetZapLogger() *zap.Logger {
	return zapLogger
}

func buildEncoder(logConf LogConfig) zapcore.Encoder {
	// 区分环境配置
	var encoderConfig zapcore.EncoderConfig
	if logConf.Env == "prod" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05.000]") // zapcore.ISO8601TimeEncoder

	// 输出格式json还是console
	if logConf.Encoding == "json" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 输出日志等级
func enableLevel(logConf LogConfig) zapcore.Level {
	if logConf.Env == "prod" {
		return zapcore.InfoLevel
	}
	return zapcore.DebugLevel
}

// 轮转存文件
func getLogFileWriter(logConf LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logConf.Path,
		MaxSize:    logConf.MaxSize,    // 单个文件最大尺寸，单位M （默认值100M）
		MaxBackups: logConf.MaxBackups, // 最多保留备份个数 (跟MaxAge都不配，则保留全部)
		MaxAge:     logConf.MaxAge,     // 最大时间，默认单位 day
		LocalTime:  true,               // 使用本地时间
		Compress:   false,
	}
	_ = lumberJackLogger.Rotate()
	return zapcore.AddSync(lumberJackLogger)
}
