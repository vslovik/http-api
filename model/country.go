package model

type (
	Country struct {
		ID   int64  `json:"id" db:"id"`
		Name string `json:"name" db:"name"`
		Code string `json:"code" db:"code"`
	}
)

const CountriesPerPage = 200

// Lists sections
func (manager *CountryManager) List(page int) (*[]Country, int, error) {
	var (
		err         error
		countryList []Country
		total       []int
	)
	if _, err = manager.DbMap.Select(&countryList, "SELECT SQL_CALC_FOUND_ROWS `id`, `name`, `code` FROM `Country` ORDER BY `name` LIMIT ?, ?",
			(page-1)*CountriesPerPage, CountriesPerPage); err != nil {
		return &countryList, 0, err
	}
	if _, err = manager.DbMap.Select(&total, "SELECT FOUND_ROWS()"); err != nil {
		return &countryList, 0, err
	}

	return &countryList, total[0], err
}
