package model

type (
	Competition struct {
		ID                int64             `json:"id" db:"id"`
		OptaId            int64             `json:"optaId" db:"id_opta"`
		PlayoffId         int64             `json:"PlayoffId" db:"id_playoff"`
		HeimspielId       int64             `json:"heimspielId" db:"id_heimspiel"`
		Name              string            `json:"name" db:"name"`
		Longname          string            `json:"longname" db:"longname"`
		Description       string            `json:"description" db:"description"`
		Url               string            `json:"url" db:"url"`
		Country           string            `json:"country" db:"country"`
		HasLiveticker     bool              `json:"hasLiveticker" db:"has_liveticker"`
		Grouped           bool              `json:"Grouped" db:"grouped"`
		National          bool              `json:"national" db:"national"`
		Women             bool              `json:"women" db:"women"`
		PushNotifications bool              `json:"pushNotifications" db:"pushNotifications"`
		HasLiveMatches    bool              `json:"HasLiveMatches" db:"HasLiveMatches"`
		TranslationMap    map[string]string `json:"translation" db:"-"`
	}

	CompetitionTranslationRecord struct {
		ID            int64  `db:"id"`
		CompetitionId int64  `db:"competitionId"`
		Language      string `db:"language"`
		Name          string `db:"name"`
	}
)

const CompetitionsPerPage = 10

// Read competition identified by Id
func (manager *CompetitionManager) Read(competitionId string) (competition *Competition, err error) {
	obj, err := manager.DbMap.Get(Competition{}, competitionId)
	if obj != nil && err == nil {
		competition = obj.(*Competition)
		competition.TranslationMap, err = manager.ReadTranslationMap(competitionId)

		return competition, err
	}

	return nil, err
}

// Read translation for competition identified by id
func (manager *CompetitionManager) ReadTranslationMap(competitionID string) (translationMap map[string]string, err error) {
	translationMap = make(map[string]string)
	var records []CompetitionTranslationRecord
	_, err = manager.DbMap.Select(&records, "SELECT `language`, `name` FROM `CompetitionTranslation` WHERE `competitionId`=? ORDER BY `language`", competitionID)

	if len(records) == 0 {
		return nil, err
	}

	translationMap = make(map[string]string)
	for _, record := range records {
		translationMap[record.Language] = record.Name
	}

	return translationMap, err
}

// Add competition
func (manager *CompetitionManager) AddTranslation(competitionId string, language string, translation string) (err error) {
	_, err = manager.DbMap.Exec("REPLACE INTO `CompetitionTranslation` (`competitionId`, `name`, `language`) VALUES (?,?,?)", competitionId, translation, language)

	return err
}

// Remove competition
func (manager *CompetitionManager) RemoveTranslation(competitionId string, language string) (err error) {
	_, err = manager.DbMap.Exec("DELETE FROM `CompetitionTranslation` WHERE `competitionId`=? AND `language`=?", competitionId, language)

	return err
}

// List competitions
func (manager *CompetitionManager) List(page int) (*[]Competition, int, error) {
	var (
		err             error
		competitionList []Competition
		total           []int
	)
	if _, err := manager.DbMap.Select(&competitionList, "SELECT SQL_CALC_FOUND_ROWS * FROM `division` LIMIT ?, ?", (page-1)*CompetitionsPerPage, CompetitionsPerPage); err != nil {
		return &competitionList, 0, err
	}
	if _, err = manager.DbMap.Select(&total, "SELECT FOUND_ROWS()"); err != nil {
		return &competitionList, 0, err
	}

	return &competitionList, total[0], err
}

// Check if competition exists
func (manager *CompetitionManager) Exists(competitionId string) (res bool, err error) {
	if obj, err := manager.DbMap.Get(Competition{}, competitionId); obj != nil {
		return obj.(*Competition) != nil, err
	}

	return false, err
}

// Publish competition
func (manager *CompetitionManager) Publish(competitionId string) (err error) {
	_, err = manager.DbMap.Exec("UPDATE `CompetitionSectionBinding` SET `competitionPublished`=1 WHERE `competitionId`=?", competitionId)

	return err
}

// Hide competition
func (manager *CompetitionManager) Hide(competitionId string) (err error) {
	_, err = manager.DbMap.Exec("UPDATE `CompetitionSectionBinding` SET `competitionPublished`=0 WHERE `competitionId`=?", competitionId)

	return err
}
