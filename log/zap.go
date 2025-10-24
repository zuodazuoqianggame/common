package log

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/akkuman/zaploki"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// gorm的日志
// https://github.com/paul-milne/zap-loki
// https://github.com/moul/zapgorm2

// RotateFileHandler 自定义支持日志轮转的文件日志处理器
func InitZapLogger(dir string, name string, lokiAddress string) (*zap.Logger, error) {
	baseLogPath := path.Join(dir, name)
	// 配置 file-rotatelogs
	writer, err := rotatelogs.New(
		baseLogPath+"-%Y%m%d%H%M.log",
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
	)
	if err != nil {
		return nil, err
	}

	// 配置 zap 编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 创建 zap 核心
	// encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建 cores 列表
	var cores []zapcore.Core

	// 1. 文件输出核心
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zap.InfoLevel,
	)
	cores = append(cores, fileCore)

	// 2. 控制台输出核心 (可选)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)
	cores = append(cores, consoleCore)

	// 3. Loki 输出核心 (如果配置了地址)
	if lokiAddress != "" {
		lokiCore, err := createLokiCore(lokiAddress, name, encoderConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Loki core: %w", err)
		}
		cores = append(cores, lokiCore)
	}

	core := zapcore.NewTee(cores...)
	// 创建 zap 核心
	z := zap.New(core, zap.AddStacktrace(zap.PanicLevel))
	zap.ReplaceGlobals(z)
	return z, nil
}

// 创建 Loki 核心
func createLokiCore(lokiAddress, appName string, encoderConfig zapcore.EncoderConfig) (zapcore.Core, error) {
	// 创建 Loki 客户端
	// loki client config
	cfg := &zaploki.LokiClientConfig{
		// the loki api url
		URL: lokiAddress, //"http://admin:admin@loki.xxx.com/api/prom/push",
		// (optional, default: severity) the label's key to distinguish log's level, it will be added to Labels map
		LevelName: "level",
		// (optional, default: zapcore.InfoLevel) logs beyond this level will be sent
		SendLevel: zapcore.InfoLevel,
		// the labels which will be sent to loki, contains the {levelname: level}
		Labels: map[string]string{
			"app": appName,
		},
	}

	return zaploki.NewLokiCore(cfg)
}
