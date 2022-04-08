package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/admin/internal/api"
	"github.com/zldongly/email_server/app/admin/internal/conf"
	"github.com/zldongly/email_server/pkg/app"
	"github.com/zldongly/email_server/pkg/errors"
	"github.com/zldongly/email_server/pkg/middleware"
	"github.com/zldongly/email_server/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

var ProviderSet = wire.NewSet(NewServer)

func NewServer(cfg *conf.Server, i api.Admin, log *zap.SugaredLogger) app.Server {
	log = log.With("module", "server")

	e := gin.New()
	e.Use(middleware.Log(log), gin.Recovery())

	r := e.Group("/v1/admin")

	r.POST("/template", func(c *gin.Context) {
		var req api.CreateTemplateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(err)
			err = errors.WithCode(err, http.StatusBadRequest, "参数解析错误")
			web.ResponseHttp(c, nil, err)
			return
		}

		reply, err := i.CreateTemplate(c, &req)
		web.ResponseHttp(c, reply, err)
	})

	srv := &http.Server{
		Addr:    cfg.Http.Addr,
		Handler: e,
	}

	return &Server{srv: srv}
}

type Server struct {
	srv *http.Server
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// shutdown 会处理完正在执行的连接（优雅退出）
	return s.srv.Shutdown(ctx)
}
