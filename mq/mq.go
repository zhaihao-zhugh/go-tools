package mq

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type ChannelConfig struct {
	Type     string   `json:"type"`
	Exchange string   `json:"exchange"`
	Queue    string   `json:"queue"`
	Key      []string `json:"key"`
	Durable  bool     `json:"durable"`
}

type MQCLIENT struct {
	Conn *amqp.Connection
	Lock sync.Mutex
}

type MQCHANNEL struct {
	*amqp.Channel
}

type Consumer struct {
	Client      *MQCLIENT
	Channel     *amqp.Channel
	Logger      *log.Logger
	RdData      <-chan amqp.Delivery
	NotifyClose chan *amqp.Error
	Config      *ChannelConfig
}

type Producer struct {
	Client      *MQCLIENT
	Channel     *amqp.Channel
	Logger      *log.Logger
	NotifyClose chan *amqp.Error
	Config      *ChannelConfig
}

type PublishMsg struct {
	Exchange string      `json:"exchange"`
	Key      string      `json:"key"`
	Body     interface{} `json:"body"`
}

type MqMsg struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

const (
	reconnectDelay = 15 * time.Second
)

func NewConnect(path string) *amqp.Connection {
	conn, err := amqp.Dial(path)
	if err != nil {
		log.Printf("MQ Client error: %s", err.Error())
		for {
			time.Sleep(reconnectDelay)
			log.Println("MQ Client reconnet")
			conn, err = amqp.Dial(path)
			if err == nil {
				break
			}
		}
	}
	log.Println("MQ Client seccess")
	return conn
}

func (client *MQCLIENT) NewChannel() *MQCHANNEL {
	ch, err := client.Conn.Channel()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &MQCHANNEL{ch}
}

func (client *MQCLIENT) NewProducer(cfg *ChannelConfig) *Producer {
	p := Producer{
		Config: cfg,
		Client: client,
		Logger: log.New(os.Stdout, "[mq-"+cfg.Exchange+"] ", log.LstdFlags|log.Lshortfile),
	}

	for !p.handelConnect() {
		p.Logger.Printf("Failed to open channel. Retrying...")
		time.Sleep(reconnectDelay)
	}
	p.Logger.Printf("Producer type:%s, exchange:%s, queue:%s, key:%s  \n", cfg.Type, cfg.Exchange, cfg.Queue, cfg.Key)
	return &p
}

func (client *MQCLIENT) NewConsumer(cfg *ChannelConfig) *Consumer {
	c := Consumer{
		Config: cfg,
		Client: client,
		Logger: log.New(os.Stdout, "[mq-"+cfg.Exchange+"] ", log.LstdFlags|log.Lshortfile),
	}

	for !c.handelConnect() {
		c.Logger.Printf("Failed to open channel. Retrying...")
		time.Sleep(reconnectDelay)
	}
	c.Logger.Printf("Consumer type:%s, exchange:%s, queue:%s, key:%s  \n", cfg.Type, cfg.Exchange, cfg.Queue, cfg.Key)
	return &c
}

func (p *Producer) handelConnect() bool {
	ch, err := p.Client.Conn.Channel()
	if err != nil {
		p.Logger.Printf(err.Error())
		return false
	}

	err = ch.ExchangeDeclare(
		p.Config.Exchange, // name
		p.Config.Type,     // type
		p.Config.Durable,  // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		p.Logger.Printf(err.Error())
		return false
	}
	p.NotifyClose = make(chan *amqp.Error)
	p.Channel = ch
	p.Channel.NotifyClose(p.NotifyClose)
	return true
}

func (p *Producer) PublishMsg(data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for _, v := range p.Config.Key {
		err := p.Channel.Publish(
			p.Config.Exchange,
			v,
			false, //mandatory：true：如果exchange根据自身类型和消息routeKey无法找到一个符合条件的queue，那么会调用basic.return方法将消息返还给生产者。false：出现上述情形broker会直接将消息扔掉
			false, //如果exchange在将消息route到queue(s)时发现对应的queue上没有消费者，那么这条消息不会放入队列中。当与消息routeKey关联的所有queue(一个或多个)都没有消费者时，该消息会通过basic.return方法返还给生产者。
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        buf,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Producer) PublishMsgWithKey(key string, data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = p.Channel.Publish(
		p.Config.Exchange,
		key,
		false, //mandatory：true：如果exchange根据自身类型和消息routeKey无法找到一个符合条件的queue，那么会调用basic.return方法将消息返还给生产者。false：出现上述情形broker会直接将消息扔掉
		false, //如果exchange在将消息route到queue(s)时发现对应的queue上没有消费者，那么这条消息不会放入队列中。当与消息routeKey关联的所有queue(一个或多个)都没有消费者时，该消息会通过basic.return方法返还给生产者。
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buf,
		})
	return err
}

func (c *Consumer) handelConnect() bool {
	ch, err := c.Client.Conn.Channel()
	if err != nil {
		c.Logger.Printf(err.Error())
		return false
	}

	err = ch.ExchangeDeclare(
		c.Config.Exchange, // name
		c.Config.Type,     // type
		c.Config.Durable,  // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		c.Logger.Printf(err.Error())
		return false
	}

	q, err := ch.QueueDeclare(
		c.Config.Queue,   // name
		c.Config.Durable, // durable
		true,             // delete when unused
		false,            // exclusive 是否私有
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		c.Logger.Printf(err.Error())
		return false
	}

	for _, v := range c.Config.Key {
		err = ch.QueueBind(
			q.Name,            // queue name
			v,                 // routing key
			c.Config.Exchange, // exchange
			false,             //	noWait
			nil,
		)
		if err != nil {
			c.Logger.Printf(err.Error())
			return false
		}
	}

	//订阅消息，并不是把mq的消息直接写到msgs，不需要死循环订阅，订阅之后mq有消息就会往msgs里写
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		c.Logger.Printf(err.Error())
		return false
	}

	c.NotifyClose = make(chan *amqp.Error)
	c.Channel = ch
	c.RdData = msgs
	c.Channel.NotifyClose(c.NotifyClose)
	return true
}
