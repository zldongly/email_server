package data

import (
	"context"
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/job/internal/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewJobRepo, NewData, NewMongo)

type Data struct {
	db  *mongo.Client
	log *zap.SugaredLogger
}

// NewData .
func NewData(db *mongo.Client, log *zap.SugaredLogger) (*Data, func(), error) {
	log = log.With("module", "data")
	d := &Data{
		db:  db,
		log: log,
	}
	return d, func() {
		d.db.Disconnect(context.TODO())
	}, nil
}

func NewMongo(cfg *conf.Data, log *zap.SugaredLogger) *mongo.Client {
	log = log.With("module", "data/mongo")

	clientOptions := options.Client().
		ApplyURI(cfg.Mongo.Uri).
		SetAuth(options.Credential{
			AuthSource: cfg.Mongo.AuthSource,
			Username:   cfg.Mongo.Username,
			Password:   cfg.Mongo.Password,
		})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
