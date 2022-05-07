package data

import (
	"context"
	"encoding/hex"
	"github.com/zldongly/email_server/app/job/internal/service"
	"github.com/zldongly/email_server/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Record struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	SendTime   int64              `bson:"send_time"`
	Receivers  []string           `bson:"receivers"`
	TemplateId string             `bson:"template_id"`
	Name       string             `bson:"name"`
	Content    string             `bson:"content"`
	IsSuccess  int8               `bson:"is_success"`
	Message    string             `bson:"message"` // 失败错误原因
}

func (r *Record) TableName() string {
	return "record"
}

func (r *jobRepo) CreateRecord(ctx context.Context, record *service.Record) (*service.Record, error) {
	re := &Record{
		SendTime:   record.SendTime,
		Receivers:  record.Receivers,
		TemplateId: record.TemplateId,
		Name:       record.Name,
		Content:    record.Content,
		IsSuccess:  record.IsSuccess,
		Message:    record.Message,
	}
	col := r.data.db.Database("test").Collection(re.TableName())

	res, err := col.InsertOne(ctx, re)
	if err != nil {
		return nil, errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mongo.insert")
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		re.Id = id
	} else {
		return nil, errors.WithCode(errors.WithStack(errors.New("insert one result type")), http.StatusInternalServerError, "")
	}

	record.Id = hex.EncodeToString(re.Id[:])

	return record, nil
}
