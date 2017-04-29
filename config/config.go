// Package config provides ...
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
)

var (
	// global config variable
	Conf *Config
	// default configuration
	defaultConfig = `{
  "web": {
    "enable_http": true,
    "http_port": 8080,
    "enable_https": false,
    "https_port": 8081,
    "cert_file": "",
    "key_file": ""
  }
}`
)

// config represents configuration
type Config struct {
	Version string `json:"-"`
	Web     *web   `json:"web"`
}

// web config
type web struct {
	EnableHttp  bool   `json:"enable_http"`
	HttpPort    int    `json:"http_port"`
	EnableHttps bool   `json:"enable_https"`
	HttpsPort   int    `json:"https_port"`
	CertFile    string `json:"cert_file"`
	KeyFile     string `json:"key_file"`
}

// global config init
func Init(version string) {
	Conf = &Config{
		Version: version,
	}

	err := readConfig(Conf)
	if err != nil {
		panic(err)
	}
}

// read config from ~/.imgee or defaultConfig
func readConfig(cfg *Config) error {
	f, err := configFile()
	if err != nil {
		return err
	}

	// read data
	data, err := ioutil.ReadFile(f)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			data, err = initConfig(f)
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}

// config file path
func configFile() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.imgee", nil
}

// initConfig creates new config file
func initConfig(f string) ([]byte, error) {
	err := ioutil.WriteFile(f, []byte(defaultConfig), 0600)
	if err != nil {
		return nil, err
	}

	return []byte(defaultConfig), nil
}

// update config when set everytime
func updateConfig(cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	f, err := configFile()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f, data, 0600)
	return err
}

// implement plugin functions
func (cfg *Config) Command() string {
	return "set"
}

// TODO sub command
func (cfg *Config) Cmds() []string {
	return []string{}
}

// TODO executive function
func (cfg *Config) Exec(args string) error {
	return nil
}
