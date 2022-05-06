package data

import (
	"context"
	"github.com/zldongly/email_server/app/admin/internal/conf"
	"github.com/zldongly/email_server/app/admin/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"testing"
	"time"
)

var (
	_recordRepo service.RecordRepo
)

func init() {
	logz, err := newZap()
	if err != nil {
		log.Fatalln(err)
	}
	mdb := NewMongo(&conf.Data{
		Mongo: &conf.Mongo{
			Uri:        "mongodb://127.0.0.1:27017",
			AuthSource: "test",
			Username:   "testuser",
			Password:   "123456",
		},
	}, logz)
	data, _, err := NewData(mdb, logz)
	if err != nil {
		log.Fatalln(err)
	}
	_recordRepo = NewRecordRepo(data, logz)
	logz.Info("init")
}

func TestListRecord(t *testing.T) {
	id := "624992bb361386fcfe1e07bf"
	id = ""
	isSuccess := 0
	list, err := _recordRepo.ListRecord(context.Background(), id, "注册", &isSuccess, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(list)
	for _, i := range list {
		t.Log(i)
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
