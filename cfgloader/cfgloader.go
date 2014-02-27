package cfgloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

const (
	SCORES_CONFIG_PATH_KEY = "SCORES_CONFIG_PATH"
)

type Config struct {
	MysqlConfig *MysqlConfig `json:"mysql"`
}

func (config *Config) Validate() (err error) {
	var errmessages []string

	errmessages = append(errmessages, config.MysqlConfig.Validate()...)

	if len(errmessages) > 0 {
		output := "Config Loader: Invalid application config!\n"

		for _, message := range errmessages {
			output += fmt.Sprintf("> %s\n", message)
		}

		return errors.New(output)
	}

	return nil
}

func GetConfig(env string) *Config {
	cdir := os.Getenv(SCORES_CONFIG_PATH_KEY)

	if cdir == "" {
		_, currentFilename, _, _ := runtime.Caller(0)
		cdir = path.Join(path.Dir(currentFilename), "/../")
	}

	cpath := path.Join(cdir, fmt.Sprintf("config_%s.json", env))

	return NewConfig(cpath)
}

func NewConfig(cpath string) (config *Config) {
	file, err := ioutil.ReadFile(cpath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	err = config.Validate()
	if err != nil {
		panic(err)
	}

	return config
}
