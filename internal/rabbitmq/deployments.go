package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DeploymentsQueue_Name       = "deployments"
	DeploymentsQueue_Durable    = false
	DeploymentsQueue_AutoDelete = false
	DeploymentsQueue_Exclusive  = false
	DeploymentsQueue_NoWait     = false
	DeploymentsVhost            = "groundctl"
)

func deploymentsQueue(conn *amqp.Connection) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		DeploymentsQueue_Name,
		DeploymentsQueue_Durable,
		DeploymentsQueue_AutoDelete,
		DeploymentsQueue_Exclusive,
		DeploymentsQueue_NoWait,
		nil,
	)

	return ch, q, err
}

type deploymentsListener struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}

func NewDeploymentsListener(ctx context.Context, conn *amqp.Connection) (*deploymentsListener, error) {
	ch, q, err := deploymentsQueue(conn)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.ConsumeWithContext(
		ctx,
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &deploymentsListener{
		ch:   ch,
		msgs: msgs,
	}, nil
}

func (l *deploymentsListener) Close() error {
	return l.ch.Close()
}

func (l *deploymentsListener) Consume(ctx context.Context) (*Deployment, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	case msg := <-l.msgs:
		var deployment Deployment
		err := json.Unmarshal(msg.Body, &deployment)
		if err != nil {
			return nil, err
		}

		return &deployment, nil
	}
}

type deploymentsClient struct {
	ch *amqp.Channel
	q  amqp.Queue
}

func NewDeploymentsClient(conn *amqp.Connection) (*deploymentsClient, error) {
	ch, q, err := deploymentsQueue(conn)
	if err != nil {
		return nil, err
	}

	return &deploymentsClient{
		ch: ch,
		q:  q,
	}, nil
}

func (c *deploymentsClient) Close() error {
	return c.ch.Close()
}

func (c *deploymentsClient) Send(ctx context.Context, obj Deployment) error {
	out, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = c.ch.PublishWithContext(
		ctx,
		"",
		c.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        out,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
