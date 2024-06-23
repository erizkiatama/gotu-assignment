package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	Convey("Test get config", t, func() {
		Convey("empty", func() {
			cfg := Get()

			So(cfg.Database, ShouldBeZeroValue)
		})
	})
}

func TestLoad(t *testing.T) {
	Convey("Test load config", t, func() {
		Convey("success", func() {
			cfg, err := Load(
				WithConfigFolder("../config/"),
				WithConfigFile("config"),
				WithConfigType("yaml"),
			)

			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeEmpty)
		})
	})
}
