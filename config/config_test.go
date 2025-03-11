package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var config *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned, and values are as expected", func() {
				config, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)

				So(config.BindAddr, ShouldEqual, "localhost:29700")
				So(config.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(config.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(config.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)

				So(config.MongoConfig.ClusterEndpoint, ShouldEqual, "localhost:27017")
				So(config.MongoConfig.Database, ShouldEqual, "router")
				So(config.MongoConfig.Collections, ShouldResemble, map[string]string{RoutesCollection: "routes", RedirectsCollection: "redirects"})
				So(config.MongoConfig.Username, ShouldEqual, "")
				So(config.MongoConfig.Password, ShouldEqual, "")
				So(config.MongoConfig.IsSSL, ShouldEqual, false)
				So(config.MongoConfig.QueryTimeout, ShouldEqual, 15*time.Second)
				So(config.MongoConfig.ConnectTimeout, ShouldEqual, 5*time.Second)
				So(config.MongoConfig.IsStrongReadConcernEnabled, ShouldEqual, false)
				So(config.MongoConfig.IsWriteConcernMajorityEnabled, ShouldEqual, true)

				So(cfg.Kafka.ContentUpdatedGroup, ShouldEqual, "dis-routing-api-poc")
				So(cfg.Kafka.ProducerTopic, ShouldEqual, "routing-updated")
				So(cfg.Kafka.Addr, ShouldResemble, []string{"localhost:9092", "localhost:9093", "localhost:9094"})
				So(cfg.Kafka.Version, ShouldEqual, "1.0.2")
				So(cfg.Kafka.OffsetOldest, ShouldBeTrue)
				So(cfg.Kafka.NumWorkers, ShouldEqual, 1)
				So(cfg.Kafka.SecProtocol, ShouldEqual, "")
				So(cfg.Kafka.SecCACerts, ShouldEqual, "")
				So(cfg.Kafka.SecClientCert, ShouldEqual, "")
				So(cfg.Kafka.SecClientKey, ShouldEqual, "")
				So(cfg.Kafka.SecSkipVerify, ShouldBeFalse)
				So(cfg.Kafka.MaxBytes, ShouldEqual, 2000000)
				So(cfg.Kafka.ConsumerMinBrokersHealthy, ShouldEqual, 1)
				So(cfg.Kafka.ProducerMinBrokersHealthy, ShouldEqual, 1)
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
