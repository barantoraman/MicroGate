package err

import (
	"encoding/json"
	"net/http"
)

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(map[string]any{"err": message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(js)
}
