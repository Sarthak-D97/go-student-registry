package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK                  = "OK"
	StatusCreated             = "Created"
	StatusBadRequest          = "Bad Request"
	StatusInternalServerError = "Internal Server Error"
	StatusError               = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errMsgs []string
	for _, e := range err {
		switch e.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, e.Field()+" is required")
		case "email":
			errMsgs = append(errMsgs, e.Field()+" must be a valid email")
		default:
			errMsgs = append(errMsgs, e.Field()+" is not valid")
		}

	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
