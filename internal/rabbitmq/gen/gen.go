package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:generate go run ./gen.go

type QueueConfig struct {
	Name       string `yaml:"name"`
	FileName   string `yaml:"filename,omitempty"`
	Durable    bool   `yaml:"durable,omitempty"`
	AutoDelete bool   `yaml:"auto_delete,omitempty"`
	Exclusive  bool   `yaml:"exclusive,omitempty"`
	NoWait     bool   `yaml:"no_wait,omitempty"`
	Vhost      string `yaml:"vhost,omitempty"`
}

var caser = cases.Title(language.English)

func capitalize(value string) string {
	return strings.ReplaceAll(caser.String(strings.ReplaceAll(value, "_", " ")), " ", "")
}

func singular(value string) string {
	return pluralize.NewClient().Singular(value)
}

func plural(value string) string {
	return pluralize.NewClient().Plural(value)
}

func main() {
	// Load queue config
	queueConfigBytes, err := os.ReadFile("queues.yml")
	if err != nil {
		logrus.Panicf("failed to read in queues.yml: %v", err)
	}
	var queueConfig []QueueConfig
	err = yaml.Unmarshal(queueConfigBytes, &queueConfig)
	if err != nil {
		logrus.Panicf("failed to unmarshal queues.yml: %v", err)
	}

	// Load template file
	tmpl, err := template.New("queue.go.tmpl").Funcs(template.FuncMap{
		"capitalize": capitalize,
		"singular":   singular,
		"plural":     plural,
	}).ParseFiles("./queue.go.tmpl")
	if err != nil {
		panic(err)
	}

	for _, config := range queueConfig {
		// Use key name as filename if not present
		if config.FileName == "" {
			config.FileName = config.Name
		}
		// Create queue file
		file, err := os.Create(fmt.Sprintf("../%s.go", config.FileName))
		if err != nil {
			logrus.Errorf("failed to create %s.go: %v", config.FileName, err)
			continue
		}
		// Template content into it
		err = tmpl.Execute(file, config)
		if err != nil {
			logrus.Errorf("failed to template %s.go: %v", config.FileName, err)
		}
	}
}
