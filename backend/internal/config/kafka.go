package config

import "os"

type KafkaConfig struct {
	Broker string
	Topic  string
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Broker: getEnv("KAFKA_BROKER", "localhost: 9092"),
		Topic:  getEnv("KAFKA_TOPIC", "vocabulary-events"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
