package config

type AppConfig struct {
	ApplicationName string
	Port            string
	System          SystemConfig
	Otlp            OtlpConfig
	GormDsn         string
	Kafka           KafkaConfig
}

type SystemConfig struct {
	Message string
}

type OtlpConfig struct {
	Endpoint string
}

type KafkaConfig struct {
	Brokers []string
}
