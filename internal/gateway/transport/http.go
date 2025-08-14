package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	"github.com/barantoraman/microgate/internal/gateway/endpoints"
	context2 "github.com/barantoraman/microgate/pkg/ctx"
	errs "github.com/barantoraman/microgate/pkg/err"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	r := mux.NewRouter()

	r.MethodNotAllowedHandler = http.HandlerFunc(errs.MethodNotAllowedResponse)
	r.NotFoundHandler = http.HandlerFunc(errs.NotFoundResponse)

	r.Use(authentication, rateLimit, secureHeaders, enableCORS, recoverPanic)
	r.Handle("/v1/task", requireAuthenticatedUser(httpTransport.NewServer(
		ep.AddTaskEndpoint,
		decodeHTTPAddTaskRequest,
		encodeResponse))).Methods(http.MethodPost)

	r.Handle("/v1/task", requireAuthenticatedUser(httpTransport.NewServer(
		ep.ListTaskEndpoint,
		decodeHTTPListTaskRequest,
		encodeResponse))).Methods(http.MethodGet)

	r.Handle("/v1/task/{id}", requireAuthenticatedUser(httpTransport.NewServer(
		ep.DeleteTaskEndpoint,
		decodeHTTPDeleteTaskRequest,
		encodeResponse))).Methods(http.MethodDelete)

	r.Handle("/v1/signup", httpTransport.NewServer(
		ep.SignUpEndpoint,
		decodeSignUpRequest,
		encodeResponse)).Methods(http.MethodPost)

	r.Handle("/v1/login", httpTransport.NewServer(
		ep.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse)).Methods(http.MethodPost)

	r.Handle("/v1/logout", httpTransport.NewServer(
		ep.LogoutEndpoint,
		decodeLogoutRequest,
		encodeResponse)).Methods(http.MethodPost)

	return r
}

func decodeHTTPListTaskRequest(_ context.Context, r *http.Request) (any, error) {
	var req endpoints.ListTaskRequest

	req.UserID = context2.GetUser(r).UserID
	return req, nil
}

func decodeHTTPAddTaskRequest(_ context.Context, r *http.Request) (any, error) {
	var req endpoints.AddTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	req.Task.UserID = context2.GetUser(r).UserID
	return req, nil
}

func decodeHTTPDeleteTaskRequest(_ context.Context, r *http.Request) (any, error) {
	var req endpoints.DeleteTaskRequest
	req.TaskID = mux.Vars(r)["id"]
	req.UserID = context2.GetUser(r).UserID
	return req, nil
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (any, error) {
	var req endpoints.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (any, error) {
	var req endpoints.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLogoutRequest(_ context.Context, r *http.Request) (any, error) {
	var req endpoints.LogoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch err {
	case entity.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
	case entity.ErrDuplicateEmail:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]any{"error": err})
}
