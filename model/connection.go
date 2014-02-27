package model

import (
	"database/sql"
	"fmt"
	"github.com/motain/gorp"
	_ "github.com/motain/mysql"
	"github.com/motain/sCoreAdmin/cfgloader"
)

func GetMySQLConnection(config *cfgloader.MysqlConfig) (*sql.DB, error) {
	connectionDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.Database)

	return sql.Open("mysql", connectionDSN)
}

func GetDbMap(db *sql.DB) *gorp.DbMap {
	dbMap := &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "UTF8",
		},
	}

	dbMap.AddTableWithName(Country{}, "Country").SetKeys(true, "ID")
	dbMap.AddTableWithName(Section{}, "Section").SetKeys(true, "ID")
	dbMap.AddTableWithName(SectionRecord{}, "Section").SetKeys(true, "ID")
	dbMap.AddTableWithName(SectionTranslationRecord{}, "SectionTranslation").SetKeys(true, "ID")
	dbMap.AddTableWithName(Competition{}, "division").SetKeys(true, "ID")
	dbMap.AddTableWithName(TopCompetitionRecord{}, "TopCompetition").SetKeys(true, "ID")

	return dbMap
}
