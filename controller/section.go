package controller

import (
	"github.com/motain/gorp"
	"github.com/motain/sCoreAdmin/model"
	"github.com/motain/mux"
	"net/http"
	"log/syslog"
	"fmt"
	"encoding/json"
	"reflect"
)

func (controller *SectionController) getLogPrefix() string {
	return "Section controller: "
}

func (controller *SectionController) GetResponseMessage(msgType string) string {
	response_map := map[string]string{
		"CreateSuccess":   "Section created",
		"CreateError":     "Error while section creating",
		"ValidationError": "Section validation error",

		"BadRequest": "Bad request",

		"ReadSuccess": "Section found",
		"NotFound":    "Section not found",
		"ReadError":   "Error while section reading",

		"UpdateSuccess": "Section updated",
		"UpdateError":   "Error while section deleting",

		"DeleteSuccess": "Section deleted",
		"DeleteError":   "Error while section deleting",

		"AddTranslationSuccess": "Section translated",
		"AddTranslationError":   "Error while section translating",

		"RemoveTranslationSuccess": "Section translation deleted",
		"RemoveTranslationError":   "Error while section translation deleting",

		"AddCompetitionSuccess": "Section competition added",
		"AddCompetitionError": "Error while adding competition to section",

		"RemoveCompetitionSuccess": "Competition removed from section",
		"RemoveCompetitionError": "Error while removing competition from section",

		"ListSuccess":   "Sections listed",
		"ListError":     "Error while section listing",
	}

	return response_map[msgType]
}

type (
	SectionControllerResponseData struct {
		Section *model.Section `json:"section"`
	}
)

func (controller *SectionController) GetResponseData(entity *model.Section) SectionControllerResponseData {
	return SectionControllerResponseData{Section: entity}
}

// Core (CRUD) routing
func (controller *SectionController) GetCoreRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "POST" && vars["id"] == "":
		controller.ControllerMethodWrapper(controller.Create, w, r)
	case r.Method == "GET" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Read, w, r)
	case r.Method == "GET" && vars["id"] == "":
		controller.ControllerMethodWrapper(controller.List, w, r)
	case r.Method == "PUT" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Update, w, r)
	case r.Method == "DELETE" && vars["id"] != "":
		controller.ControllerMethodWrapper(controller.Delete, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}

// Routing
func (controller *SectionController) ManagementRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch {
	case r.Method == "POST" && vars["entity_id"] != "" && vars["child"] == "translation" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.AddTranslation, w, r)
	case r.Method == "DELETE" && vars["entity_id"] != "" && vars["child"] == "translation" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.RemoveTranslation, w, r)
	case r.Method == "POST" && vars["entity_id"] != "" && vars["child"] == "competition" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.AddCompetition, w, r)
	case r.Method == "DELETE" && vars["entity_id"] != "" && vars["child"] == "competition" && vars["child_id"] != "":
		controller.ControllerMethodWrapper(controller.RemoveCompetition, w, r)
	default:
		controller.WriteBadRequestResponse(w)
	}
}


// Create section
func (controller *SectionController) Create(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		entity model.Section
		err    error
	)
	if err = json.NewDecoder(r.Body).Decode(&entity); err != nil {
		return controller.GetErrorControllerResponse(http.StatusBadRequest, GeneralFormValidationErrorCode, controller.GetResponseMessage("BadRequest"))
	}
	if err = entity.IsValid(); err != nil {
		return controller.GetErrorControllerResponse(http.StatusBadRequest, GeneralFormValidationErrorCode, controller.GetResponseMessage("ValidationError"))
	}
	if err = model.GetSectionManager(dbMap).Create(&entity); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("CreateError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("CreateSuccess"), controller.GetResponseData(&entity))
}

// Read section
func (controller *SectionController) Read(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		entity *model.Section
		err    error
	)
	vars := mux.Vars(r)
	if entity, err = model.GetSectionManager(dbMap).Read(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ReadError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}
	if entity == nil {
		return controller.GetErrorControllerResponse(http.StatusNotFound, GeneralServerErrorCode, controller.GetResponseMessage("NotFound"))
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("ReadSuccess"), controller.GetResponseData(entity))
}


// Update section
func (controller *SectionController) Update(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {

	var (
		data    map[string]interface{}
		exists  bool
		entity  *model.Section
		err     error
	)

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return controller.GetErrorControllerResponse(http.StatusBadRequest, GeneralFormValidationErrorCode, controller.GetResponseMessage("BadRequest"))
	}

	if isValid := controller.ValidateData(data); !isValid {
		return controller.GetErrorControllerResponse(http.StatusBadRequest, GeneralFormValidationErrorCode, controller.GetResponseMessage("BadRequest"))
	}

	vars := mux.Vars(r)

	if exists, err = model.GetSectionManager(dbMap).Exists(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ReadError"), err))

		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	if exists == false {
		return controller.GetErrorControllerResponse(http.StatusNotFound, GeneralServerErrorCode, controller.GetResponseMessage("NotFound"))
	}

	if err = model.GetSectionManager(dbMap).UpdateRecord(data); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("UpdateError"), err))

		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	if entity, err = model.GetSectionManager(dbMap).Read(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.GetResponseMessage("ReadError"), err))

		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("UpdateSuccess"), controller.GetResponseData(entity))
}

func (controller *SectionController) ValidateData(data map[string]interface{}) bool {
	sample := map[string]interface{}{
		"name": "",
		"published": 0,
		"priority": 0,
	}

	for key, val := range data {
		flag := false
		if _, ok := sample[key]; !ok {
			flag = true
		}
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Float64 && v.Kind() != reflect.Bool && v.Kind() != reflect.String {
			flag = true
		}
		if flag {
			delete(data, key)
		}
	}

	if n := len(data); n == 0 {

		return false;
	}

	return true;
}

// Delete section
func (controller *SectionController) Delete(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetSectionManager(dbMap).Delete(vars["id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details:records %s\n", controller.getLogPrefix(), controller.GetResponseMessage("DeleteError"), err))

		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("DeleteSuccess"), Empty{})
}


// Translate section
func (controller *SectionController) AddTranslation(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		err error
		exists bool
		data map[string]string
	)

	vars := mux.Vars(r)

	if exists, err = model.GetSectionManager(dbMap).Exists(vars["entity_id"]); err != nil {
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

	if err = model.GetSectionManager(dbMap).AddTranslation(vars["entity_id"], vars["child_id"], data["name"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("AddTranslationError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("AddTranslationSuccess"), Empty{})
}


// Delete translation
func (controller *SectionController) RemoveTranslation(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetSectionManager(dbMap).RemoveTranslation(vars["entity_id"], vars["child_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("RemoveTranslationError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("RemoveTranslationSuccess"), Empty{})
}


// Add competition
func (controller *SectionController) AddCompetition(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetSectionManager(dbMap).AddCompetition(vars["entity_id"], vars["child_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("AddCompetitionError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("AddCompetitionSuccess"), Empty{})
}


// Remove competition
func (controller *SectionController) RemoveCompetition(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var err error
	vars := mux.Vars(r)
	if err = model.GetSectionManager(dbMap).RemoveCompetition(vars["entity_id"], vars["child_id"]); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details: %s\n", controller.getLogPrefix(), controller.GetResponseMessage("RemoveCompetitionError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("RemoveCompetitionSuccess"), Empty{})
}

// List sections
func (controller *SectionController) List(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse {
	var (
		err         error
		sectionList *[]model.Section
	)
	if sectionList, err = model.GetSectionManager(dbMap).List(); err != nil {
		new(syslog.Writer).Err(fmt.Sprintf("%s %s Details:records %s\n", controller.getLogPrefix(), controller.GetResponseMessage("ListError"), err))
		return controller.GetErrorControllerResponse(http.StatusInternalServerError, GeneralServerErrorCode, "DB: "+err.Error())
	}

	return controller.GetSuccessControllerResponse(http.StatusOK, controller.GetResponseMessage("ListSuccess"), sectionList)
}
