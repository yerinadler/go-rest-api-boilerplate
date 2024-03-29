package config

type AppConfig struct {
	Application ApplicationConfig
	Otlp        OtlpConfig
	Kafka       KafkaConfig
}

type ApplicationConfig struct {
	Name string
}

type OtlpConfig struct {
	Endpoint string
}

type KafkaConfig struct {
	Brokers []string
}
