// Package config provides ...
package config

var Conf *Config

type Config struct {
	Version string
}

func Init(version string) {
	Conf = &Config{
		Version: version,
	}
}
