package cfgloader

import (
	"fmt"
)

type MysqlConfig struct {
	Database string `json:"database"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func (config *MysqlConfig) Validate() (messages []string) {

	messageFormat := "Configuration Error! Entry %s should be defined"

	if config.Database == "" {
		messages = append(messages, fmt.Sprintf(messageFormat, "mysql.database"))
	}

	if config.Host == "" {
		messages = append(messages, fmt.Sprintf(messageFormat, "mysql.host"))
	}

	if config.Port == 0 {
		messages = append(messages, fmt.Sprintf(messageFormat, "mysql.port"))
	}

	return messages
}
