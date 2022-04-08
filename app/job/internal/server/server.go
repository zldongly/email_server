package server

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/wire"
	"github.com/zldongly/email_server/app/job/internal/api"
	"github.com/zldongly/email_server/app/job/internal/conf"
	"github.com/zldongly/email_server/pkg/app"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewServer, NewKafkaClient)

func NewServer(kc sarama.Client, a api.Job, log *zap.SugaredLogger) app.Server {
	log = log.With("module", "server")
	return &Server{
		ctx:    context.Background(),
		cancel: nil,
		kc:     kc,
		log:    log,
		a:      a,
	}
}

type Server struct {
	ctx    context.Context
	cancel func()
	a      api.Job
	kc     sarama.Client
	log    *zap.SugaredLogger
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	group, err := sarama.NewConsumerGroupFromClient("c1", s.kc)
	if err != nil {
		return err
	}

	return group.Consume(s.ctx, []string{"email"}, &kafkaConsumer{
		a:   s.a,
		log: s.log,
	})
}

func (s *Server) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func NewKafkaClient(cfg *conf.Server,log *zap.SugaredLogger) sarama.Client {
	log = log.With("module", "server/kafka")

	config := sarama.NewConfig()
	client, err := sarama.NewClient(cfg.Kafka.Addrs, config)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

type kafkaConsumer struct {
	a   api.Job
	log *zap.SugaredLogger
}

func (c *kafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *kafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *kafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	for msg := range claim.Messages() {
		c.log.Debug("msg.val:", string(msg.Value))

		var mail = new(api.Email)
		if err := json.Unmarshal(msg.Value, mail); err != nil {
			c.log.Error(err)
			continue
		}

		if err := c.a.SendEmail(context.Background(), mail); err != nil {
			c.log.Error(err)
			continue
		}

		sess.MarkMessage(msg, "")
	}

	return nil
}
