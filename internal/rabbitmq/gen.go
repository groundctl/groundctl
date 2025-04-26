//go:build ignore

package main

import (
	"bytes"
	"go/format"
	"os"
	"text/template"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type RabbitMQConfig struct {
	Queues []struct {
		Name       string `yaml:"name"`
		FileName   string `yaml:"filename,omitempty"`
		Durable    bool   `yaml:"durable,omitempty"`
		AutoDelete bool   `yaml:"auto_delete,omitempty"`
		Exclusive  bool   `yaml:"exclusive,omitempty"`
		NoWait     bool   `yaml:"no_wait,omitempty"`
		Vhost      string `yaml:"vhost,omitempty"`
	} `yaml:"queues"`
}

var pluralClient = pluralize.NewClient()

var funcMap = template.FuncMap{
	"camel":      strcase.ToCamel,
	"lowerCamel": strcase.ToLowerCamel,
	"singular":   pluralClient.Singular,
	"plural":     pluralClient.Plural,
}

func main() {
	// Load queue config
	rabbitmqConfigBytes, err := os.ReadFile("config.yml")
	if err != nil {
		logrus.Panicf("failed to read in queues.yml: %v", err)
	}
	var rabbitmqConfig RabbitMQConfig
	err = yaml.Unmarshal(rabbitmqConfigBytes, &rabbitmqConfig)
	if err != nil {
		logrus.Panicf("failed to unmarshal config.yml: %v", err)
	}

	applyDefaults(&rabbitmqConfig)

	generateQueues(rabbitmqConfig)
	generateVhosts(rabbitmqConfig)
}

func applyDefaults(config *RabbitMQConfig) {
	// Fill blank vhosts with "default"
	for i := range config.Queues {
		if config.Queues[i].Vhost == "" {
			config.Queues[i].Vhost = "default"
		}
	}
}

func generateQueues(config RabbitMQConfig) {
	// Load template file
	tmpl, err := template.New("queues.go.tmpl").Funcs(funcMap).ParseFiles("queues.go.tmpl")
	if err != nil {
		panic(err)
	}
	// Create queue file
	file, err := os.Create("queues.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Template content
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config.Queues)
	if err != nil {
		panic(err)
	}
	// Format content
	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = file.Write(p)
	if err != nil {
		panic(err)
	}
}

func generateVhosts(config RabbitMQConfig) {
	// Collect all unique vhost names
	vhosts := map[string]struct{}{}
	for _, q := range config.Queues {
		// Use "default" for name of blank vhost
		vhosts[q.Vhost] = struct{}{}
	}

	// Load template file
	tmpl, err := template.New("vhosts.go.tmpl").Funcs(funcMap).ParseFiles("vhosts.go.tmpl")
	if err != nil {
		panic(err)
	}
	// Create queue file
	file, err := os.Create("vhosts.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Template content
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vhosts)
	if err != nil {
		panic(err)
	}
	// Format content
	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = file.Write(p)
	if err != nil {
		panic(err)
	}
}
