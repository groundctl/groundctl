package rabbitmq

import (
	"fmt"

	"github.com/groundctl/groundctl/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func openConnection(vhost string, username string, password string) (*amqp.Connection, error) {
	connStr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		username,
		password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
		vhost,
	)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
