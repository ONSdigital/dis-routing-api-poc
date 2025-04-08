package config

import (
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
)

type MongoConfig = mongodb.MongoDriverConfig

// KafkaTLSProtocol informs service to use TLS protocol for kafka
const KafkaTLSProtocol = "TLS"

// Config represents service configuration for dis-routing-api-poc
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	MongoConfig
	RouterAPIURL string `envconfig:""`
	Kafka        *Kafka
}

// Kafka contains the config required to connect to Kafka
type Kafka struct {
	ContentUpdatedGroup       string   `envconfig:"KAFKA_CONTENT_UPDATED_GROUP"`
	ProducerTopic             string   `envconfig:"KAFKA_PRODUCER_TOPIC"`
	Addr                      []string `envconfig:"KAFKA_ADDR"`
	Version                   string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest              bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	NumWorkers                int      `envconfig:"KAFKA_NUM_WORKERS"`
	SecProtocol               string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts                string   `envconfig:"KAFKA_SEC_CA_CERTS"            json:"-"`
	SecClientCert             string   `envconfig:"KAFKA_SEC_CLIENT_CERT"         json:"-"`
	SecClientKey              string   `envconfig:"KAFKA_SEC_CLIENT_KEY"          json:"-"`
	SecSkipVerify             bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	MaxBytes                  int      `envconfig:"KAFKA_MAX_BYTES"`
	ConsumerMinBrokersHealthy int      `envconfig:"KAFKA_CONSUMER_MIN_BROKERS_HEALTHY"`
	ProducerMinBrokersHealthy int      `envconfig:"KAFKA_PRODUCER_MIN_BROKERS_HEALTHY"`
}

var cfg *Config

const (
	RoutesCollection    = "RoutesCollection"
	RedirectsCollection = "RedirectsCollection"
)

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:29700",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MongoConfig: MongoConfig{
			ClusterEndpoint:               "localhost:27017",
			Username:                      "",
			Password:                      "",
			Database:                      "router",
			Collections:                   map[string]string{RoutesCollection: "routes", RedirectsCollection: "redirects"},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongodb.TLSConnectionConfig{
				IsSSL: false,
			},
		},
		Kafka: &Kafka{
			ContentUpdatedGroup:       "dis-routing-api-poc",
			ProducerTopic:             "routing-updated",
			Addr:                      []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			Version:                   "1.0.2",
			OffsetOldest:              true,
			NumWorkers:                1,
			SecProtocol:               "",
			SecCACerts:                "",
			SecClientCert:             "",
			SecClientKey:              "",
			SecSkipVerify:             false,
			MaxBytes:                  2000000,
			ConsumerMinBrokersHealthy: 1,
			ProducerMinBrokersHealthy: 1,
		},
	}

	return cfg, envconfig.Process("", cfg)
}
