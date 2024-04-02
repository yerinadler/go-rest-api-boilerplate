package config

type AppConfig struct {
	Application ApplicationConfig
	Port        string
	System      SystemConfig
	Otlp        OtlpConfig
	Gorm        GormConfig
}

type ApplicationConfig struct {
	Name string
}

type GormConfig struct {
	Dsn string
}

type SystemConfig struct {
	Message string
}

type OtlpConfig struct {
	Endpoint string
}
