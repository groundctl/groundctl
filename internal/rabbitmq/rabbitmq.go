package rabbitmq

import (
	"fmt"

	"github.com/groundctl/groundctl/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate go run ./gen.go

// Opens an AMQP 0-9-1 connection to a RabbitMQ server on a given vhost using username/password credentials
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
