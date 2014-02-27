package controller

import (
	"github.com/motain/gorp"
	"github.com/motain/sCoreAdmin/model"
	"github.com/motain/mux"
	"net/http"
	"log/syslog"
	"fmt"
	"strconv"
	"encoding/json"
)

type (
	CompetitionControllerResponseData struct {
		Competition *model.Competition `json:"competition"`
	}

	CompetitionControllerPaginatedList struct {
		Competitions *[]model.Competition `json:"competition"`
		Total        int            `json:"total"`
	}
)

func (controller *CompetitionController) getLogPrefix() string {
	return "Competition controller: "
}

func (controller *CompetitionController) GetResponseData(entity *model.Competition) CompetitionControllerResponseData {
	return CompetitionControllerResponseData{Competition: entity}
}

func (controller *CompetitionController) GetResponseMessage(msgType string) string {
	response_map := map[string]string{

		"BadRequest": "Bad request",

		"ReadSuccess": "Competition found",
		"NotFound":    "Competition not found",
		"ReadError":   "Error while competition reading",

		"AddTranslationSuccess": "Competition translated",
		"AddTranslationError":   "Error while competition translating",

		"RemoveTranslationSuccess": "Competition translation deleted",
		"RemoveTranslationError":   "Error while competition translation deleting",

		"PublishSuccess": "Competition published",
		"PublishError":   "Error while publishing competition",

		"HideSuccess": "Competition hidden",
		"HideError":   "Error while hiding competition",

		"ListSuccess": "Competitions listed",
		"ListError":   "Competitions list error",
	}

	return response_map[msgType]
}

// Routing
func (controller *CompetitionController) GetCoreRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "POST" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Publish, w, r)
	case r.Method == "GET" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Read, w, r)
	case r.Method == "DELETE" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Hide, w, r)
	case r.Method == "GET" && vars["id"] == "":
		controller.ControllerMethodWrapper(controller.List, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}

// Routing
func (controller *CompetitionController) ManagementRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "POST" && vars["entity_id"] != "" && vars["child"] == "translation" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.AddTranslation, w, r)
	case r.Method == "DELETE" && vars["entity_id"] != "" && vars["child"] == "translation" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.RemoveTranslation, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}

// Read section
func (controller *CompetitionController) Read(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		entity *model.Competition
		err    error
	)
	vars := mux.Vars(r)
	if entity, err = model.GetCompetitionManager(dbMap).Read(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ReadError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}
	if entity == nil {
		return controller.GetErrorControllerResponse(http.StatusNotFound, GeneralServerErrorCode, controller.GetResponseMessage("NotFound"))
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("ReadSuccess"), controller.GetResponseData(entity))
}


// Add translation
func (controller *CompetitionController) AddTranslation(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		err error
		exists bool
		data map[string]string
	)

	vars := mux.Vars(r)

	if exists, err = model.GetCompetitionManager(dbMap).Exists(vars["entity_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ReadError"), err))

		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	if exists == false {
		return controller.GetErrorControllerResponse(http.StatusNotFound, GeneralServerErrorCode, controller.GetResponseMessage("NotFound"))
	}

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return controller.GetErrorControllerResponse(http.StatusBadRequest, GeneralFormValidationErrorCode, controller.GetResponseMessage("BadRequest"))
	}

	if _, ok := data["name"]; !ok {
		return controller.GetErrorControllerResponse(http.StatusBadRequest, GeneralFormValidationErrorCode, controller.GetResponseMessage("BadRequest"))
	}

	if err = model.GetCompetitionManager(dbMap).AddTranslation(vars["entity_id"], vars["child_id"], data["name"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("AddTranslationError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("AddTranslationSuccess"), Empty{})
}


// Delete translation
func (controller *CompetitionController) RemoveTranslation(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetCompetitionManager(dbMap).RemoveTranslation(vars["entity_id"], vars["child_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("RemoveTranslationError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("RemoveTranslationSuccess"), Empty{})
}


// List competitions
func (controller *CompetitionController) List(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		err             error
		competitionList *[]model.Competition
		total           int
		page            int
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

	if competitionList, total, err = model.GetCompetitionManager(dbMap).List(page); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details:records %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ListError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("ListSuccess"), CompetitionControllerPaginatedList{Competitions: competitionList, Total: total})
}


// Publish competition
func (controller *CompetitionController) Publish(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetCompetitionManager(dbMap).Publish(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("PublishError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("PublishSuccess"), Empty{})
}


// Hide competition
func (controller *CompetitionController) Hide(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetCompetitionManager(dbMap).Hide(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("HideError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("HideSuccess"), Empty{})
}
