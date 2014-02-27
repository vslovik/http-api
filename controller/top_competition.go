package controller

import (
	"github.com/motain/gorp"
	"github.com/motain/sCoreAdmin/model"
	"github.com/motain/mux"
	"net/http"
	"log/syslog"
	"fmt"
)

func (controller *TopCompetitionController) getLogPrefix() string {
	return "Section controller: "
}

func (controller *TopCompetitionController) GetResponseMessage(msgType string) string {
	response_map := map[string]string{
		"BadRequest": "Bad request",

		"ReadSuccess": "Top competition found",
		"NotFound":    "Top competition not found",
		"ReadError":   "Error while top competition reading",

		"DeleteSuccess": "Top competition deleted",
		"DeleteError":   "Error while top competition deleting",

		"AddCompetitionSuccess": "Competition added to top competition",
		"AddCompetitionError": "Error while adding competition to top competition",

		"RemoveCompetitionSuccess": "Competition removed from top competition",
		"RemoveCompetitionError": "Error while removing competition from top competition",
	}

	return response_map[msgType]
}

type (
	TopCompetitionControllerResponseData struct {
		Competitions *[]model.TopCompetitionCompetition `json:"competitions"`
	}
)

func (controller *TopCompetitionController) GetResponseData(competitions *[]model.TopCompetitionCompetition) TopCompetitionControllerResponseData {
	return TopCompetitionControllerResponseData{Competitions: competitions}
}

// Core routing
func (controller *TopCompetitionController) GetCoreRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "GET" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Read, w, r)
	case r.Method == "DELETE" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Delete, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}

// Routing
func (controller *TopCompetitionController) ManagementRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "POST" && vars["entity_id"] != "" && vars["child"] == "competition" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.AddCompetition, w, r)
	case r.Method == "DELETE" && vars["entity_id"] != "" && vars["child"] == "competition" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.RemoveCompetition, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}

// Read top list
func (controller *TopCompetitionController) Read(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		competitionList *[]model.TopCompetitionCompetition
		err             error
	)
	vars := mux.Vars(r)
	if competitionList, err = model.GetTopCompetitionManager(dbMap).ListCompetitions(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ReadError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	if competitionList == nil {
		return controller.GetErrorControllerResponse(http.StatusNotFound, GeneralServerErrorCode, controller.GetResponseMessage("NotFound"))
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("ReadSuccess"), controller.GetResponseData(competitionList))
}

// Delete top list
func (controller *TopCompetitionController) Delete(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		err error
	)
	vars := mux.Vars(r)
	if err = model.GetTopCompetitionManager(dbMap).Delete(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("DeleteError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("DeleteSuccess"), Empty{})
}

// Add competition to top competition
func (controller *TopCompetitionController) AddCompetition(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetTopCompetitionManager(dbMap).AddCompetition(vars["entity_id"], vars["child_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("AddCompetitionError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("AddCompetitionSuccess"), Empty{})
}


// Delete translation from top competition
func (controller *TopCompetitionController) RemoveCompetition(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetTopCompetitionManager(dbMap).RemoveCompetition(vars["entity_id"], vars["child_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("RemoveCompetitionError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("RemoveCompetitionSuccess"), Empty{})
}
