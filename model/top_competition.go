package model

type (
	TopCompetitionRecord struct {
		ID                  int64  `json:"-" db:"id"`
		CompetitionId       int    `json:"-" db:"competitionId"`
		CountryId           int    `json:"-" db:"countryId"`
		CompetitionPriority int    `json:"-" db:"competitionPriority"`
	}

	TopCompetitionCompetition struct {
		CompetitionId       int64  `json:"id" db:"competitionId"`
		Name                string `json:"name" db:"name"`
		CompetitionPriority int    `json:"priority" db:"competitionPriority"`
	}
)

// Add competition to top
func (manager *TopCompetitionManager) AddCompetition(countryId string, competitionId string) (err error) {
	_, err = manager.DbMap.Exec("REPLACE INTO `TopCompetition` (`countryId`, `competitionId`) VALUES(?,?)", countryId, competitionId);

	return err
}

// Remove competition from top
func (manager *TopCompetitionManager) RemoveCompetition(countryId string, competitionId string) (err error) {
	_, err = manager.DbMap.Exec("DELETE FROM `TopCompetition` WHERE `countryId`=? AND `competitionId`=?", countryId, competitionId)

	return err
}

// Read competition top list identified by country
func (manager *TopCompetitionManager) ListCompetitions(countryId string) (*[]TopCompetitionCompetition, error) {
	var competitionList []TopCompetitionCompetition

	_, err := manager.DbMap.Select(&competitionList, "SELECT `t`.`competitionId`, `t`.`competitionPriority`, `c`.`name` FROM `TopCompetition` AS `t` INNER JOIN `Division` AS `c` ON `t`.`competitionId` =  `c`.`id` WHERE `t`.`countryId`=? ORDER BY `t`.`competitionPriority`", countryId)
	if(len(competitionList) == 0) {
		return nil, err
	}

	return &competitionList, err
}

// Delete competition top list identified by country
func (manager *TopCompetitionManager) Delete(countryId string)  error {
	_, err := manager.DbMap.Exec("DELETE FROM `TopCompetition` WHERE `countryId`=?", countryId)

	return err
}
