package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	enabled bool
)

// Config RabbitMQ配置
type Config struct {
	Enabled  bool
	URL      string
	Exchange string
	Queue    string
}

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ(config Config) error {
	if !config.Enabled {
		logrus.Info("RabbitMQ is disabled")
		enabled = false
		return nil
	}

	var err error
	conn, err = amqp.Dial(config.URL)
	if err != nil {
		logrus.Warnf("Failed to connect to RabbitMQ: %v, message queue disabled", err)
		enabled = false
		return err
	}

	channel, err = conn.Channel()
	if err != nil {
		logrus.Warnf("Failed to open channel: %v", err)
		enabled = false
		return err
	}

	enabled = true
	logrus.Infof("✅ RabbitMQ connected successfully: %s", config.URL)
	return nil
}

// IsEnabled 检查RabbitMQ是否启用
func IsEnabled() bool {
	return enabled
}

// Publish 发布消息
func Publish(exchange, routingKey string, message interface{}) error {
	if !enabled {
		logrus.Warn("RabbitMQ is not enabled, message dropped")
		return fmt.Errorf("rabbitmq is not enabled")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

// Close 关闭连接
func Close() error {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		return conn.Close()
	}
	return nil
}
