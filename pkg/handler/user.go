package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JuD4Mo/go_api_web_user/internal/user"
	"github.com/JuD4Mo/go_lib_response/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {
	//Instancia de un router de Gorilla Mux
	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodedError),
	}

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAll,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdate,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDelete,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: %v", err.Error()))
	}

	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	path := mux.Vars(r)
	req := user.GetReq{
		ID: path["id"],
	}

	return req, nil
}

func decodeGetAll(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}

func decodeUpdate(_ context.Context, r *http.Request) (interface{}, error) {

	var req user.UpdateReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: %v", err.Error()))
	}

	path := mux.Vars(r)
	id := path["id"]

	req.ID = id

	return req, nil
}

func decodeDelete(_ context.Context, r *http.Request) (interface{}, error) {
	path := mux.Vars(r)
	id := path["id"]

	req := user.DeleteReq{
		ID: id,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodedError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)

}
