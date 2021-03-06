// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/zldongly/email_server/app/job/internal/api"
	"github.com/zldongly/email_server/app/job/internal/conf"
	"github.com/zldongly/email_server/app/job/internal/data"
	"github.com/zldongly/email_server/app/job/internal/server"
	"github.com/zldongly/email_server/app/job/internal/service"
	"github.com/zldongly/email_server/pkg/app"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func initServer(confServer *conf.Server, mail *conf.Mail, confData *conf.Data, sugaredLogger *zap.SugaredLogger) (app.Server, func(), error) {
	client := server.NewKafkaClient(confServer, sugaredLogger)
	mongoClient := data.NewMongo(confData, sugaredLogger)
	dataData, cleanup, err := data.NewData(mongoClient, sugaredLogger)
	if err != nil {
		return nil, nil, err
	}
	messageRepo := data.NewJobRepo(dataData, sugaredLogger)
	sender := service.NewMailUseCase(mail, sugaredLogger)
	messageUseCase := service.NewMessageUseCase(messageRepo, sender, sugaredLogger)
	job := api.NewJob(messageUseCase, sugaredLogger)
	appServer := server.NewServer(client, job, sugaredLogger)
	return appServer, func() {
		cleanup()
	}, nil
}
