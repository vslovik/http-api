package controller

import (
	"github.com/motain/gorp"
	"github.com/motain/sCoreAdmin/cfgloader"
	"github.com/motain/sCoreAdmin/model"
	"encoding/json"
	"net/http"
	"log"
	"os"
)

const (
	SuccessResponseStatus          = "ok"
	ErrorResponseStatus            = "error"
	SuccessCode                    = 0
	GeneralFormValidationErrorCode = 1000
	GeneralServerErrorCode         = 1001
	RequestInvalidErrorCode        = 1101
)

type (
	Controller struct {
		Env          string
		Config       *cfgloader.Config
	}

	ControllerResponseBody struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	ControllerResponse struct {
		Body         ControllerResponseBody
		StatusCode   int
		Error        error
		ErrorMessage string
		Headers      *http.Header
	}

	ControllerCrudTranslateInterface interface {
		ControllerMethodWrapper(method ControllerMethod, w http.ResponseWriter, r *http.Request)
	}

	ControllerMethod func(dbMap *gorp.DbMap, r *http.Request) *ControllerResponse

	Empty struct{}
)

func (controller *Controller) postControllerMethod(response *ControllerResponse, w http.ResponseWriter, r *http.Request) {

	if response.Error != nil || response.ErrorMessage != "" {
		if response.StatusCode == 0 {
			response.StatusCode = http.StatusInternalServerError
		}

		errorString := response.ErrorMessage
		if response.ErrorMessage == "" {
			errorString = response.Error.Error()
		}
		response.Body = ControllerResponseBody{Status: ErrorResponseStatus, Code: GeneralServerErrorCode, Message: errorString}
	}

	if response.Headers != nil && response.Headers.Get("Location") != "" {
		w.Header().Set("Location", response.Headers.Get("Location"))
	}

	w.WriteHeader(response.StatusCode)

	body, err := json.Marshal(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(body)

}

func (controller *Controller) ControllerMethodWrapper(method ControllerMethod, w http.ResponseWriter, r *http.Request) {
	response := &ControllerResponse{Error: nil}

	w.Header().Set("Content-Type", "application/json")

	connection, err := model.GetMySQLConnection(controller.Config.MysqlConfig);

	if err != nil {
		panic("Can not connect to db!..")
	}

	dbmap := model.GetDbMap(connection)

	if controller.Env == "dev" {
		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "sCoreAdmin:", log.Lmicroseconds))
	}

	defer connection.Close()


	if response.Error == nil {
		response = method(dbmap, r)
	}

	controller.postControllerMethod(response, w, r)
}

func (controller *Controller) GetSuccessControllerResponse(httpStatusCode int, msg string, data interface{}) *ControllerResponse {
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json")
	body := ControllerResponseBody{
		Status:  SuccessResponseStatus,
		Code:    SuccessCode,
		Message: msg,
		Data:    data,
	}
	return &ControllerResponse{
		Headers:    headers,
		StatusCode: httpStatusCode,
		Body:       body,
	}
}

func (controller *Controller) GetErrorControllerResponse(httpStatusCode int, errCode int, msg string) *ControllerResponse {
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json")
	return &ControllerResponse{
		Headers:    headers,
		StatusCode: httpStatusCode,
		Body:       controller.GetErrorResponseBody(errCode, msg),
	}
}

func (controller *Controller) GetErrorResponseBody(errCode int, msg string) ControllerResponseBody {
	type Empty struct{}
	errBody := ControllerResponseBody{
		Status:  ErrorResponseStatus,
		Code:    errCode,
		Message: msg,
		Data:    Empty{},
	}

	return errBody
}

func (controller *Controller) WriteBadRequestResponse(w http.ResponseWriter) {
	type Empty struct{}
	w.WriteHeader(http.StatusBadRequest)
}
