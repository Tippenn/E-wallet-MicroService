package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Tippenn/E-wallet-MicroService/config"
)

type RabbitMQBroker struct {
	conn  *amqp.Connection
	pubCh *amqp.Channel
	mu    sync.Mutex // a channel is not safe for concurrent publishing
	cfg   *config.RabbitMQ
}

var (
	once     sync.Once
	instance *RabbitMQBroker
)

func NewRabbitMQBroker(conf *config.Config) MessageBroker {
	once.Do(func() {
		rc := conf.RabbitMQ
		url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rc.User, rc.Password, rc.Host, rc.Port, rc.VHost)

		conn, err := amqp.Dial(url)
		if err != nil {
			panic(fmt.Errorf("dial rabbitmq: %w", err))
		}

		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			panic(fmt.Errorf("open channel: %w", err))
		}
		b := &RabbitMQBroker{conn: conn, pubCh: ch, cfg: rc}
		err = b.declare()
		if err != nil {
			conn.Close()
			panic(err)
		}
		instance = b
	})
	return instance
}

func (b *RabbitMQBroker) declare() error {
	err := b.pubCh.ExchangeDeclare(b.cfg.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}
	for _, p := range []struct{ queue, key string }{
		{b.cfg.PaymentQueue, b.cfg.PaymentRoutingKey},
		{b.cfg.StatusQueue, b.cfg.StatusRoutingKey},
	} {
		if _, err := b.pubCh.QueueDeclare(p.queue, true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare queue %s: %w", p.queue, err)
		}
		if err := b.pubCh.QueueBind(p.queue, p.key, b.cfg.Exchange, false, nil); err != nil {
			return fmt.Errorf("bind queue %s: %w", p.queue, err)
		}
	}
	return nil
}

func (b *RabbitMQBroker) Publish(ctx context.Context, routingKey string, body []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.pubCh.PublishWithContext(ctx, b.cfg.Exchange, routingKey, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			Body:         body,
		})
}

func (b *RabbitMQBroker) Consume(ctx context.Context, queue string, handle Handler) error {
	ch, err := b.conn.Channel() // each consumer gets its own channel
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.Qos(5, 0, false); err != nil {
		return err
	}
	msgs, err := ch.ConsumeWithContext(ctx, queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for m := range msgs {
		err := handle(ctx, m.Body)
		if err != nil {
			m.Nack(false, false)
			continue
		}
		m.Ack(false)
	}
	return nil
}

func (b *RabbitMQBroker) Close() error {
	b.pubCh.Close()
	return b.conn.Close()
}
