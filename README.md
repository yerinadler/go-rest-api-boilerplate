# Go Rest API Project
This is a Go project created with Echo V4, Viper, Gorm, etc. This could be used as a foundational project. This project comes with pre-configured observability stack based on LGTM stack by Grafana and is instrumented with Opentelemetry SDK

## Technologies used
- Echo v4 - API Framework
- Viper - Configuration Management
- Opentelemetry - Observability library
- Grafana Tempo - For trace visualisation
- Grafana Loki - Log backend
- Fluentbit - Log forwarder

## How to run the project
The `docker-compose.yml` file provided with this project includes everything to setup and run the projects. Issue the command on the terminal to
launch the project

```bash
docker compose up -d --build
```

## Convention
Since this is an opinionated project, I lay down the foundation based on my own use case

### Project Structure
The project is based on simple configuration of 3-tier architecture with the following main directories serving each layer

- Presentation Layer is represented as `api` directory
- Business Layer (Domain Layer) is represented as `services` directory
- The Data Access Layer is represented as `models` directory holding ORM specific files e.g. GORM definitions or MongoDB models

> These directories rest inside the `internal` directory

#### Shared packages
The shared packages are all located in the `pkg` directory. This directory does not store business specific code, only shared modules e.g. response objects, utility functions, error codes, etc.

### Configuration Management
For local development, the file `config.yaml` will be used to store local configuration. This file is not meant to be committed into the version control software

#### Config unmarshalling
The main configuration object is `config/config.go` where all configurations (both from YAML and environment variables) are unmarshalled into. Thus, whenever the new configuration is added to the YAML file, you can add it to the `config/config.go`

#### How to override the configuration with environment variables
To override nested configuration like

```yaml
otlp:
  endpoint: localhost:4317
```

Use the convention of `<PARENT>_<CHILD>` (underscore_delimited)

For example, to override the above YAML for Opentelemetry endpoint. Use below environment variable format

```bash
OTLP_ENDPOINT=localhost:4317
```