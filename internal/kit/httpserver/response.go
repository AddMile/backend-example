package http

import "net/http"

type envelope map[string]any

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := encodeJSON(w, r, status, env)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func successResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := encodeJSON(w, r, status, message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	message := http.StatusText(http.StatusUnauthorized)
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func TooManyRequests(w http.ResponseWriter, r *http.Request) {
	message := http.StatusText(http.StatusTooManyRequests)
	errorResponse(w, r, http.StatusTooManyRequests, message)
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

func OK(w http.ResponseWriter, r *http.Request, message any) {
	successResponse(w, r, http.StatusOK, message)
}

func Created(w http.ResponseWriter, r *http.Request, message any) {
	successResponse(w, r, http.StatusCreated, message)
}
