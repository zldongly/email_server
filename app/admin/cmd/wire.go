//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/admin/internal/api"
	"github.com/zldongly/email_server/app/admin/internal/conf"
	"github.com/zldongly/email_server/app/admin/internal/data"
	"github.com/zldongly/email_server/app/admin/internal/server"
	"github.com/zldongly/email_server/app/admin/internal/service"
	"github.com/zldongly/email_server/pkg/app"
	"go.uber.org/zap"
)

func initServer(*conf.Server, *conf.Data, *zap.SugaredLogger) (app.Server, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service.ProviderSet, api.ProviderSet))
}
