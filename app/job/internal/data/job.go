package data

import (
	"github.com/zldongly/email_server/app/job/internal/service"
	"go.uber.org/zap"
)

type jobRepo struct {
	data *Data
	log  *zap.SugaredLogger
}

func NewJobRepo(data *Data, log *zap.SugaredLogger) service.MessageRepo {
	log = log.With("module", "data/jobRepo")
	return &jobRepo{
		data: data,
		log:  log,
	}
}
