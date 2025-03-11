package broker

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xLeSHka/mentorLinkSchool/internal/pkg/config"
	"github.com/xLeSHka/mentorLinkSchool/internal/transport/http/handler/ws"
	"go.uber.org/fx"
	"log"
	"sync"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
	group    string
}

func NewProducer(config config.Config) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{

		"bootstrap.servers": config.KafkaAddress,
	})
	if err != nil {
		return nil, err
	}
	prod := &Producer{
		producer: p,
		group:    config.KafkaGroupId,
		topic:    config.KafkaTopic,
	}
	prod.Run()
	return prod, nil
}
func (p *Producer) Run() {
	go func() {
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
}
func (p *Producer) Send(message *ws.Message) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, nil)
	return err
}

type Consumer struct {
	consumer *kafka.Consumer
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
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": opts.Config.KafkaAddress,
		"group.id":          opts.Config.KafkaGroupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	c.Subscribe(opts.Config.KafkaTopic, nil)

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
			ev := c.consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				var msg ws.Message
				err := json.Unmarshal(e.Value, &msg)
				if err != nil {
					log.Printf("Error unmarshalling message: %v\n", err)
					continue
				}
				c.wsconn.WriteMessage(&msg)
			case kafka.Error:
				log.Println(e)
				c.SetR(false)
			default:
				log.Println("ignored %v", e)
			}
		}
		err := c.consumer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
