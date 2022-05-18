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

func (r *templateRepo) CreateTemplate(ctx context.Context, t *service.Template) (*service.Template, error) {
	var (
		temp = &Template{
			Name:    t.Name,
			Subject: t.Subject,
			Content: t.Content,
		}
		col = r.data.db.Database(temp.Database()).Collection(temp.TableName())
	)

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

func (r *templateRepo) ListTemplate(ctx context.Context, id, name string, pageNum, pageSize int) ([]*service.Template, error) {
	var (
		do        = new(Template)
		col       = r.data.db.Database(do.Database()).Collection(do.TableName())
		filter    = make(bson.M, 1)
		templates = make([]*Template, 0, pageSize)
		opt       = options.Find()
	)

	if id != "" {
		// id 转换 ObjectID
		var objId primitive.ObjectID
		bs, err := hex.DecodeString(id)
		if err != nil {
			err = errors.WithCode(errors.WithStack(err), http.StatusBadRequest, "decode hex")
			return nil, err
		}

		copy(objId[:], bs)

		filter["_id"] = objId
	} else if name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}
	}

	opt.SetLimit(int64(pageSize))
	opt.SetSkip(int64((pageNum - 1) * pageSize))
	opt.SetSort(bson.M{"_id": -1})

	cursor, err := col.Find(ctx, filter, opt)
	if err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mongo")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &templates); err != nil {
		err = errors.WithCode(errors.WithStack(err), http.StatusInternalServerError, "mongo decode")
		return nil, err
	}

	result := make([]*service.Template, 0, len(templates))
	for _, t := range templates {
		result = append(result,
			&service.Template{
				Id:      hex.EncodeToString(t.Id[:]),
				Name:    t.Name,
				Subject: t.Subject,
				Content: t.Content,
			})
	}

	return result, nil
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

func (t *Template) Database() string {
	return "test"
}
