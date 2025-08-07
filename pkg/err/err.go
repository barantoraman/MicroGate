package err

import (
	"encoding/json"
	"fmt"
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

func ServerErroResponse(w http.ResponseWriter) {
	msg := "The server encountered an error while processing your request"
	errorResponse(w, http.StatusInternalServerError, msg)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "The resource could not be found on the server"
	errorResponse(w, http.StatusNotFound, msg)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("The requested method '%s' is not allowed for this resource", r.Method)
	errorResponse(w, http.StatusMethodNotAllowed, msg)
}

func RateLimitExceededResponse(w http.ResponseWriter) {
	msg := "rate limit exceeded"
	errorResponse(w, http.StatusTooManyRequests, msg)
}

func InvalidAuthenticationToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	msg := "invalid or missing authentication token"
	errorResponse(w, http.StatusUnauthorized, msg)
}

func AuthenticationRequiredResponse(w http.ResponseWriter) {
	msg := "you must be authenticated to access this resource"
	errorResponse(w, http.StatusUnauthorized, msg)
}
