package controller

import (
	"github.com/motain/sCoreAdmin/cfgloader"
)

type (
	CountryController struct {
		Controller
	}

	SectionController struct {
		Controller
	}

	CompetitionController struct {
		Controller
	}

	TopCompetitionController struct {
		Controller
	}
)

func GetBaseController(env string, config *cfgloader.Config) *Controller {
	return &Controller {
			Env: env,
			Config: config,
		}
}

func GetCountryController(env string, config *cfgloader.Config) *CountryController {
	return &CountryController{
					Controller {
						Env: env,
						Config: config,
					},
				}
}

func GetSectionController(env string, config *cfgloader.Config) *SectionController {
	return &SectionController{
		Controller {
			Env: env,
			Config: config,
		},
	}
}

func GetCompetitionController(env string, config *cfgloader.Config) *CompetitionController {
	return &CompetitionController{
		Controller {
			Env: env,
			Config: config,
		},
	}
}

func GetTopCompetitionController(env string, config *cfgloader.Config) *TopCompetitionController {
	return &TopCompetitionController{
		Controller {
			Env: env,
			Config: config,
		},
	}
}
