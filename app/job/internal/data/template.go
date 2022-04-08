package data

import (
	"context"
	"encoding/hex"
	"github.com/zldongly/email_server/app/job/internal/service"
	"github.com/zldongly/email_server/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Template struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Subject string             `bson:"subject"`
	Content string             `bson:"content"`
}

func (t *Template) TableName() string {
	return "template"
}

func (r *jobRepo) GetTemplate(ctx context.Context, id string) (*service.Template, error) {
	var (
		t     Template
		col   = r.data.db.Database("test").Collection(t.TableName())
		objId primitive.ObjectID
	)

	bs, err := hex.DecodeString(id)
	if err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusBadRequest, "decode hex")
		return nil, err
	}

	copy(objId[:], bs)

	err = col.FindOne(ctx, &bson.D{{"_id", objId}}).Decode(&t)
	if err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "decode mongo result")
		return nil, err
	}

	return &service.Template{
		Id:      id,
		Name:    t.Name,
		Subject: t.Subject,
		Content: t.Content,
	}, nil
}
