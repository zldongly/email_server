//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/job/internal/api"
	"github.com/zldongly/email_server/app/job/internal/conf"
	"github.com/zldongly/email_server/app/job/internal/data"
	"github.com/zldongly/email_server/app/job/internal/server"
	"github.com/zldongly/email_server/app/job/internal/service"
	"github.com/zldongly/email_server/pkg/app"
	"go.uber.org/zap"
)

func initServer(*conf.Server, *conf.Mail, *conf.Data, *zap.SugaredLogger) (app.Server, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service.ProviderSet, api.ProviderSet))
}

