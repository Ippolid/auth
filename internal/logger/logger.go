package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func getCore(level zap.AtomicLevel) zapcore.Core {
	// Настройка вывода в консоль
	stdout := zapcore.AddSync(os.Stdout)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Цветной вывод уровня
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	// Настройка вывода в файл с ротацией
	//file := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   "logs/app.log", // Путь к файлу логов
	//	MaxSize:    10,             // Размер файла в мегабайтах
	//	MaxBackups: 3,              // Количество старых файлов для хранения
	//	MaxAge:     7,              // Количество дней для хранения файлов
	//	Compress:   true,           // Сжимать старые файлы
	//})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	//fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	// Объединение выводов в консоль и файл
	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		//zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel(level string) (zap.AtomicLevel, error) {
	var zapLevel zapcore.Level
	// Преобразуем строку в уровень zap
	if err := zapLevel.Set(strings.ToLower(level)); err != nil {
		return zap.NewAtomicLevel(), err
	}
	return zap.NewAtomicLevelAt(zapLevel), nil
}

func InitLocalLogger(level string) {
	zaplevel, err := getAtomicLevel(level)
	if err != nil {
		panic("Invalid log level: " + level)
	}

	globalLogger = zap.New(getCore(zaplevel))

}

var globalLogger *zap.Logger

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func Logger() *zap.Logger {
	return globalLogger
}
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
