package data

import (
	"context"
	"encoding/hex"
	"github.com/zldongly/email_server/app/admin/internal/service"
	"github.com/zldongly/email_server/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
)

type templateRepo struct {
	data *Data
	log  *zap.SugaredLogger
}

func NewTemplateRepo(data *Data, log *zap.SugaredLogger) service.TemplateRepo {
	return &templateRepo{
		data: data,
		log:  log.With("module", "data/templateRepo"),
	}
}

type Template struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Subject string             `bson:"subject"`
	Content string             `bson:"content"`
}

func (t *Template) TableName() string {
	return "template"
}

func (r *templateRepo) CreateTemplate(ctx context.Context, t *service.Template) (*service.Template, error) {
	temp := Template{
		Name:    t.Name,
		Subject: t.Subject,
		Content: t.Content,
	}
	col := r.data.db.Database("test").Collection(temp.TableName())

	res, err := col.InsertOne(ctx, temp)
	if err != nil {
		return nil, errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mongo.insert")
	}

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		temp.Id = id
	} else {
		return nil, errors.WithCode(errors.WithStack(errors.New("insert one result type")), http.StatusInternalServerError, "")
	}

	t.Id = hex.EncodeToString(temp.Id[:])
	return t, nil
}
