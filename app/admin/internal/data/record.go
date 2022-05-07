package data

import (
	"context"
	"encoding/hex"
	"github.com/zldongly/email_server/app/admin/internal/service"
	"github.com/zldongly/email_server/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
)

func NewRecordRepo(data *Data, log *zap.SugaredLogger) service.RecordRepo {
	return &recordRepo{
		data: data,
		log:  log.With("module", "data/recordRepo"),
	}
}

type recordRepo struct {
	data *Data
	log  *zap.SugaredLogger
}

func (r *recordRepo) ListRecord(ctx context.Context, templateId, name string, isSuccess *int, pageNum, pageSize int) ([]*service.Record, error) {
	var (
		do      = new(Record)
		col     = r.data.db.Database(do.Database()).Collection(do.TableName())
		filter  = make(bson.M, 2)
		records = make([]*Record, 0, pageSize)
		opt     = options.Find()
	)

	if templateId != "" {
		filter["template_id"] = templateId
	} else if name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}
	}

	if isSuccess != nil {
		filter["is_success"] = isSuccess
	}

	opt.SetLimit(int64(pageSize))
	opt.SetSkip(int64((pageNum - 1) * pageSize))
	opt.SetSort(bson.M{"_id": -1})

	cursor, err := col.Find(ctx, filter, opt)
	//if errors.Is(err, mongo.ErrNoDocuments) {
	//	return make([]*service.Record, 0), nil
	//}
	if err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mongo")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &records); err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mongo decode")
		return nil, err
	}

	result := make([]*service.Record, 0, len(records))
	for _, re := range records {
		result = append(result,
			&service.Record{
				Id:         hex.EncodeToString(re.Id[:]),
				SendTime:   re.SendTime,
				Receivers:  re.Receivers,
				TemplateId: re.TemplateId,
				Name:       re.Name,
				Content:    re.Content,
				IsSuccess:  re.IsSuccess,
				Message:    re.Message,
			})
	}

	return result, nil
}

type Record struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	SendTime   int64              `bson:"send_time"`
	Receivers  []string           `bson:"receivers"`
	TemplateId string             `bson:"template_id"`
	Name       string             `bson:"name"`
	Content    string             `bson:"content"`
	IsSuccess  int8               `bson:"is_success"`
	Message    string             `bson:"message"`
}

func (r *Record) TableName() string {
	return "record"
}

func (r *Record) Database() string {
	return "test"
}
