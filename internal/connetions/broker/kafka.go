package broker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"go.uber.org/fx"
	"log"
	"sync"
	"time"
)

type Producer struct {
	producer *kafka.Conn
	topic    string
	group    string
}

func NewProducer(config config.Config) (*Producer, error) {
	p, err := kafka.DialLeader(context.Background(), "tcp", config.KafkaAddress, config.KafkaTopic, 0)
	if err != nil {
		return nil, err
	}

	prod := &Producer{
		producer: p,
		group:    config.KafkaGroupId,
		topic:    config.KafkaTopic,
	}
	return prod, nil
}

func (p *Producer) Send(message *ws.Message) error {
	p.producer.SetWriteDeadline(time.Now().Add(10 * time.Second))
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = p.producer.Write(jsonData)
	return err
}

type Consumer struct {
	consumer *kafka.Conn
	topic    string
	group    string
	r        bool
	mu       *sync.RWMutex
	wsconn   *ws.WebSocket
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

func NewConsumer(opts FxOpts) (*Consumer, error) {
	c, err := kafka.DialLeader(context.Background(), "tcp", opts.Config.KafkaAddress, opts.Config.KafkaTopic, 0)
	if err != nil {
		return nil, err
	}
	cons := &Consumer{
		consumer: c,
		topic:    opts.Config.KafkaTopic,
		group:    opts.Config.KafkaGroupId,
		r:        true,
		mu:       &sync.RWMutex{},
		wsconn:   opts.Wsconn,
	}
	cons.Run()
	return cons, nil
}
func (c *Consumer) Run() {
	go func() {
		for c.R() {
			c.consumer.SetReadDeadline(time.Now().Add(10 * time.Second))
			msg, err := c.consumer.ReadMessage(10e3)
			if err == nil {
				var m ws.Message
				err := json.Unmarshal(msg.Value, &m)
				if err != nil {
					log.Printf("Error unmarshalling message: %v\n", err)
					continue
				}
				c.wsconn.WriteMessage(&m)
			} else if !err.(kafka.Error).Timeout() {
				log.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
		err := c.consumer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
