package main

import (
	"flag"
	"github.com/zldongly/email_server/app/admin/internal/conf"
	"github.com/zldongly/email_server/pkg/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"
)

var (
	_name     = "email.admin"
	_version  = "1.0.0"
	_confPath string
)

func init() {
	flag.StringVar(&_confPath, "conf", "../../configs/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// 初始化日志
	zapLog, err := newZap()
	if err != nil {
		log.Fatalln(err)
	}
	zapLog.With("server.name", _name,
		"server.version", _version)

	// 读取配置
	bs, err := ioutil.ReadFile(_confPath)
	if err != nil {
		zapLog.Fatal(err)
	}
	var cfg conf.Conf
	err = yaml.Unmarshal(bs, &cfg)
	if err != nil {
		zapLog.Fatal(err)
	}
	zapLog.Debugf("conf: %+v", cfg.Server.Http.Addr)

	srv, cleanup, err := initServer(cfg.Server, cfg.Data, zapLog)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	a := app.New(_name, _version, []app.Server{srv})
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func newZap() (*zap.SugaredLogger, error) {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2016-01-02T15:04:05"))
	}
	logCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: true,
		//Sampling *SamplingConfig
		Encoding:         "console",
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		//InitialFields map[string]interface{}
	}

	logger, err := logCfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
