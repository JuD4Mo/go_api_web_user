package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JuD4Mo/go_api_web_meta/meta"
	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		LastName  string `json:"first_name"`
		FirstName string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimitPage string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
			return
		}

		if req.FirstName == "" || req.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "First name and last name must not be empty",
			})
			return
		}

		user, err := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusCreated,
			Data:   user,
		})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "user does not exist",
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		//Filtros
		v := r.URL.Query()
		filters := Filters{
			FirstName: v.Get("first_name"),
			LastName:  v.Get("last_name"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		//Count
		quant, err := s.Count(filters)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusInternalServerError,
				Err:    err.Error(),
			})
			return
		}

		meta, err := meta.New(page, limit, quant, config.LimitPage)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusInternalServerError,
				Err:    err.Error(),
			})
			return
		}

		users, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusOK,
			Data:   users,
			Meta:   meta,
		})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateReq UpdateReq

		err := json.NewDecoder(r.Body).Decode(&updateReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
			return
		}

		//validaciones

		if updateReq.FirstName != nil && *updateReq.FirstName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "first name is required",
			})
			return
		}

		if updateReq.LastName != nil && *updateReq.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "last name is required",
			})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		err = s.Update(id, updateReq.FirstName, updateReq.LastName, updateReq.Email, updateReq.Phone)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    "user does not exist",
			})
			return
		}

		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusOK,
			Data:   map[string]string{"message": "updated!"},
		})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{
				Status: http.StatusBadRequest,
				Err:    err.Error(),
			})
		}
		json.NewEncoder(w).Encode(&Response{
			Status: http.StatusOK,
			Data:   map[string]string{"response": "deleted complete"},
		})

	}
}
