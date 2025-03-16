package broker

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"go.uber.org/fx"
	"log"
	"sync"
	"time"
)

type Producer struct {
	producer sarama.AsyncProducer
	topic    string
	group    string
}

func NewProducer(config config.Config, lc fx.Lifecycle) (*Producer, error) {
	c := sarama.NewConfig()
	c.Net.TLS.Enable = false
	c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	c.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	producer, err := sarama.NewAsyncProducer([]string{config.KafkaAddress}, c)
	if err != nil {
		return nil, err
	}
	go func() {
		for err := range producer.Errors() {
			log.Println(err)
		}
	}()
	prod := &Producer{
		producer: producer,
		group:    config.KafkaGroupId,
		topic:    config.KafkaTopic,
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return prod.Close()
		},
	})
	return prod, nil
}

func (p *Producer) Send(message *ws.Message) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(p.group),
		Value: sarama.ByteEncoder(jsonData),
	}
	log.Println("Success send message to ", p.topic, p.group, "mes", string(jsonData))
	return nil
}
func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type Consumer struct {
	consumer     sarama.Consumer
	partConsumer sarama.PartitionConsumer
	topic        string
	group        string
	r            bool
	mu           *sync.RWMutex
	wsconn       *ws.WebSocket
}

func (c *Consumer) R() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.r
}
func (c *Consumer) SetR(r bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.r = r
}

type FxOpts struct {
	fx.In
	Config config.Config
	Wsconn *ws.WebSocket
}

func NewConsumer(opts FxOpts, lc fx.Lifecycle) (*Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{opts.Config.KafkaAddress}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	partConsumer, err := consumer.ConsumePartition(opts.Config.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	cons := &Consumer{
		consumer:     consumer,
		partConsumer: partConsumer,
		topic:        opts.Config.KafkaTopic,
		group:        opts.Config.KafkaGroupId,
		r:            true,
		mu:           &sync.RWMutex{},
		wsconn:       opts.Wsconn,
	}
	cons.Run()
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cons.Close()
		},
	})
	return cons, nil
}
func (c *Consumer) Run() {
	go func() {
		for c.R() {
			select {
			case msg, ok := <-c.partConsumer.Messages():
				if !ok {
					log.Println("Consumer closed the message")
					return
				}
				var m ws.Message
				err := json.Unmarshal(msg.Value, &m)
				if err != nil {
					log.Printf("Error unmarshalling message: %v\n", err)
					continue
				}
				log.Println("Received message:", string(msg.Value))
				c.wsconn.WriteMessage(&m)

			}
		}
	}()
}
func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		log.Println(err)
		return err
	}
	if err := c.partConsumer.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
