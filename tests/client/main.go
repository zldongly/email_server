package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

var _emailAddr = "xxx@gmail.com"

func main() {
	config := sarama.NewConfig()                                             //实例化个sarama的Config
	config.Producer.Return.Successes = true                                  //是否开启消息发送成功后通知 successes channel
	config.Producer.Partitioner = sarama.NewRandomPartitioner                //随机分区器
	client, err := sarama.NewClient([]string{"192.168.124.40:9092"}, config) //初始化客户端
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatalln(err)
	}

	email := &Email{
		Receivers:  []string{_emailAddr},
		TemplateId: "62498eb7ab468d0495b6a3ff",
		Param:      map[string]string{"email": _emailAddr, "code": "123456"},
	}
	bs, err := json.Marshal(email)
	if err != nil {
		log.Println(err)
		return
	}

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{Topic: "email", Key: nil, Value: sarama.ByteEncoder(bs)})
	if err != nil {
		log.Fatalf("unable to produce message: %q", err)
	}
	fmt.Printf("partition: %v, offset: %v\n", partition, offset)
}

type Email struct {
	Receivers  []string          `json:"receivers"`
	TemplateId string            `json:"template_id"`
	Param      map[string]string `json:"param"`
}
