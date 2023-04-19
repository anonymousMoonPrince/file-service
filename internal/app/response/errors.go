package response

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mailru/easyjson"
)

//easyjson:json
type ErrorResponseModel struct {
	Error ErrorModel `json:"error"`
}

//easyjson:json
type ErrorModel struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Details easyjson.RawMessage `json:"details"`
}

func ErrorResponse(w io.Writer, code, message string, details json.Marshaler) {
	detailsJSON := make([]byte, 0)
	if details != nil {
		detailsJSON, _ = details.MarshalJSON()
	}
	errorResponse := ErrorResponseModel{
		Error: ErrorModel{
			Code:    code,
			Message: message,
			Details: detailsJSON,
		},
	}
	err, _ := errorResponse.MarshalJSON()
	_, _ = w.Write(err)
}

func InternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	ErrorResponse(w, "INTERNAL_SERVER_ERROR", err.Error(), nil)
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	ErrorResponse(w, "BAD_REQUEST", err.Error(), nil)
}
