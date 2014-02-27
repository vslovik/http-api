package model

import (
	"bytes"
	"fmt"
	"strings"
	"strconv"
	"reflect"
)

type (
	Section struct {
		ID             int64                `json:"id" db:"id"`
		Name           string               `json:"name" db:"name"`
		Published      int                  `json:"published" db:"published"`
		Priority       int                  `json:"priority" db:"priority"`
		TranslationMap map[string]string    `json:"translation" db:"-"`
		Competitions   []SectionCompetition `json:"competitions" db:"-"`
	}

	SectionRecord struct {
		ID        int64  `json:"id" db:"id"`
		Name      string `json:"name" db:"name"`
		Published int    `json:"published" db:"published"`
		Priority  int    `json:"priority" db:"priority"`
	}

	SectionTranslationRecord struct {
		ID        int64  `db:"id"`
		SectionId int    `db:"sectionId"`
		Language  string `db:"language"`
		Name      string `db:"name"`
	}

	SectionCompetition struct {
		CompetitionId        int64  `json:"id" db:"competitionId"`
		Name                 string `json:"name" db:"name"`
		CompetitionPriority  int    `json:"priority" db:"competitionPriority"`
		CompetitionPublished int    `json:"published" db:"competitionPublished"`
	}
)

func (section *Section) IsValid() error {
	if section.Name == "" {
		return fmt.Errorf("Field Name cannot be empty")
	}

	return nil
}

// Add or update tranlations for section
func (manager *SectionManager) GetUpdateTranslationMapQuery(sectionId string, translationMap map[string]string) string {
	var buffer bytes.Buffer
	s := ""
	buffer.WriteString("REPLACE INTO `SectionTranslation` (`sectionId`, `language`, `name`) VALUES ")
	for language, name := range translationMap {
		if item := strings.Join([]string{sectionId, language, name}, "','"); s != "" {
			s = strings.Join([]string{s, item}, "'),('")
		} else {
			s = item
		}
	}
	if s == "" {
		return ""
	} else {
		s = strings.Join([]string{"('", s, "')"}, "")
		buffer.WriteString(s)
		return buffer.String()
	}
}

// Create a new section
func (manager *SectionManager) Create(section *Section) error {
	if err := section.IsValid(); err != nil {
		return err
	}
	trans, err := manager.DbMap.Begin()
	if err != nil {
		return err
	}
	if err := trans.Insert(section); err != nil {
		trans.Rollback()
		return err
	}
	query := manager.GetUpdateTranslationMapQuery(strconv.FormatInt(section.ID, 10), section.TranslationMap)
	if query != "" {
		if _, err := trans.Exec(query); err != nil {
			trans.Rollback()
			return err
		}
	}

	return trans.Commit()
}

// List sections
func (manager *SectionManager) List() (*[]Section, error) {
	var sectionList []Section
	_, err := manager.DbMap.Select(&sectionList, "SELECT * FROM `Section` ORDER BY `priority`")

	return &sectionList, err
}

// Read section identificeted by Id
func (manager *SectionManager) ReadRecord(sectionId string) (section *SectionRecord, err error) {
	obj, err := manager.DbMap.Get(SectionRecord{}, sectionId)
	if obj != nil && err == nil {
		section = obj.(*SectionRecord)

		return section, err
	}

	return nil, err
}

// Read section identified by Id
func (manager *SectionManager) Read(sectionId string) (section *Section, err error) {
	obj, err := manager.DbMap.Get(Section{}, sectionId)
	if obj != nil && err == nil {
		section = obj.(*Section)
		section.TranslationMap, err = manager.ReadTranslationMap(sectionId)
		section.Competitions, err = manager.ListCompetitions(sectionId)

		return section, err
	}

	return nil, err
}

// Read translation for section identified by id
func (manager *SectionManager) ReadTranslationMap(sectionID string) (translationMap map[string]string, err error) {
	var records []SectionTranslationRecord
	_, err = manager.DbMap.Select(&records, "SELECT `language`, `name` FROM `SectionTranslation` WHERE `sectionId`=? ORDER BY `language`", sectionID)

	if len(records) == 0 {
		return nil, err
	}

	translationMap = make(map[string]string)
	for _, record := range records {
		translationMap[record.Language] = record.Name
	}

	return translationMap, err
}

// Update record in the section table
func (manager *SectionManager) UpdateRecord(data map[string]interface{}) (err error) {
	var chunk string
	query := "";

	for key, val := range data {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Float64 {
			chunk = fmt.Sprintf("%s=%d", key, int(v.Float()))
		} else if v.Kind() == reflect.Bool {
			if v.Bool() {
				chunk = fmt.Sprintf("%s='1'", key)
			} else {
				chunk = fmt.Sprintf("%s='0'", key)
			}
		} else if v.Kind() == reflect.String {
			chunk = fmt.Sprintf("%s='%s'", key, v)
		}

		if len(query) == 0 {
			query = chunk
		} else {
			query = fmt.Sprintf("%s, %s", query, chunk)
		}
	}

	_, err = manager.DbMap.Exec(fmt.Sprintf("UPDATE `Section` SET %s", query));

	return err
}

// Delete section identified by id
func (manager *SectionManager) Delete(sectionId string) (err error) {
	section, err := manager.Read(sectionId)

	if section != nil && err == nil {
		trans, err := manager.DbMap.Begin()
		if err != nil {
			return err
		}
		if _, err = manager.DbMap.Exec("DELETE FROM `CompetitionSectionBinding` WHERE `sectionId`=?", sectionId); err != nil {
			trans.Rollback()
			return err
		}
		if _, err = manager.DbMap.Delete(section);  err != nil {
			trans.Rollback()
			return err
		}

		return trans.Commit()
	}

	return err
}

// Add translation
func (manager *SectionManager) AddTranslation(sectionId string, language string, translation string) (err error) {
	_, err = manager.DbMap.Exec("REPLACE INTO `SectionTranslation` (`sectionId`, `name`, `language`) VALUES (?, ?, ?)", sectionId, translation, language)

	return err
}

// Remove translation
func (manager *SectionManager) RemoveTranslation(sectionId string, language string) (err error) {
	_, err = manager.DbMap.Exec("DELETE FROM `SectionTranslation` WHERE `sectionId`=? AND `language`=?", sectionId, language)

	return err
}

// List section competitions
func (manager *SectionManager) ListCompetitions(sectionId string) ([]SectionCompetition, error) {
	var competitionList []SectionCompetition
	_, err := manager.DbMap.Select(&competitionList, "SELECT `b`.`competitionId`, `b`.`competitionPublished`, `b`.`competitionPriority`, `c`.`name` FROM `CompetitionSectionBinding` AS `b` INNER JOIN `division` AS `c` ON `b`.`competitionId` =  `c`.`id` WHERE `b`.`sectionId`=? ORDER BY `b`.`competitionPriority`", sectionId)

	return competitionList, err
}

// Add competition to section
func (manager *SectionManager) AddCompetition(sectionId string, competitionId string) (err error) {
	_, err = manager.DbMap.Exec("REPLACE INTO `CompetitionSectionBinding` (`sectionId`, `competitionId`) VALUES(?,?)", sectionId, competitionId);

	return err
}

// Remove competition from section
func (manager *SectionManager) RemoveCompetition(sectionId string, competitionId string) (err error) {
	_, err = manager.DbMap.Exec("DELETE FROM `CompetitionSectionBinding` WHERE `sectionId`=? AND `competitionId`=?", sectionId, competitionId)

	return err
}

// Check if section exists
func (manager *SectionManager) Exists(sectionId string) (res bool, err error) {
	if obj, err := manager.DbMap.Get(Section{}, sectionId); obj != nil {
		return obj.(*Section) != nil, err
	}

	return false, err
}
