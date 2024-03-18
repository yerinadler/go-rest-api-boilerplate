package config

type AppConfig struct {
	Port    string
	System  SystemConfig
	Otlp    OtlpConfig
	GormDsn string
}

type SystemConfig struct {
	Message string
}

type OtlpConfig struct {
	Endpoint string
}
