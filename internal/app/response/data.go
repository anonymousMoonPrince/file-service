package response

import (
	"encoding/json"
	"net/http"

	"github.com/mailru/easyjson"
)

//easyjson:json
type DataResponseModel struct {
	Data easyjson.RawMessage `json:"data"`
}

func Ok(w http.ResponseWriter, data json.Marshaler) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	dataJSON, _ := data.MarshalJSON()
	dataResponse := DataResponseModel{
		Data: dataJSON,
	}
	dataResponseJSON, _ := dataResponse.MarshalJSON()
	_, _ = w.Write(dataResponseJSON)
}
