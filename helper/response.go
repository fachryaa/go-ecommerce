package helper

import (
	"encoding/json"
	"github.com/fachryaa/project-assignment-synapsis-ecommerce/Responses"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, payload Responses.WebResponse) {
	response, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(payload.Code)
	w.Write(response)
}
