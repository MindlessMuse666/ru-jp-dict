package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MindlessMuse666/ru-jp-dict/backend/internal/models"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(broker string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		topic: topic,
	}
}

func (p *Producer) SendEvent(eventType string, word models.Vocabulary) error {
	event := map[string]interface{}{
		"event_type": eventType,
		"payload": map[string]interface{}{
			"id":         word.ID,
			"russian":    word.Russian,
			"japanese":   word.Japanese,
			"onyomi":     word.Onyomi,
			"kunyomi":    word.Kunyomi,
			"created_at": word.CreatedAt.Format(time.RFC3339),
			"updated_at": word.UpdatedAt.Format(time.RFC3339),
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: eventBytes,
		},
	)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
