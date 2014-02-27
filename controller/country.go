package controller

import (
	"github.com/motain/gorp"
	"github.com/motain/sCoreAdmin/model"
	"github.com/motain/mux"
	"net/http"
	"strconv"
	"log/syslog"
	"fmt"
)

type (
	CountryControllerPaginatedList struct {
		Countries *[]model.Country `json:"country"`
		Total     int        `json:"total"`
	}
)

func (controller *CountryController) getLogPrefix() string {
	return "Country controller: "
}

func (controller *CountryController) GetResponseMessage(msgType string) string {
	response_map := map[string]string{
		"ListSuccess": "Countries listed",
	}

	return response_map[msgType]
}

// Routing
func (controller *CountryController) GetCoreRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "GET" && vars["id"] == "":
		controller.ControllerMethodWrapper(controller.List, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}

// List competitions
func (controller *CountryController) List(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		err         error
		countryList *[]model.Country
		total       int
		page        int
	)

	query := r.URL.Query()

	if list, ok := query["page"]; ok {
		page, err = strconv.Atoi(list[0])
		if page == 0 {
			page = 1
		}
	} else {
		page = 1
	}
	if err != nil {
		page = 1
	}

	if countryList, total, err = model.GetCountryManager(dbMap).List(page); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details:records %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ListError"), err))

		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("ListSuccess"), CountryControllerPaginatedList{Countries: countryList, Total: total})
}
