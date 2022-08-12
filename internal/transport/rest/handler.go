package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/julienschmidt/httprouter"
)

type Bank interface {
	CreateAccount(ctx context.Context, account domain.Account) error
	ReadAccount(ctx context.Context, id int64) (domain.Account, error)
	ReadAllAccounts(ctx context.Context) ([]domain.Account, error)
	UpdateAccount(ctx context.Context, id int64, inp domain.AccountUpdateInput) error
	DeleteAccount(ctx context.Context, id int64) error
}

type Handler struct {
	bankService Bank
}

func NewHandler(bank Bank) *Handler {
	return &Handler{
		bankService: bank,
	}
}

func (h *Handler) InitRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/account", h.CreateAccount)
	router.PUT("/account", h.CreateAccount)
	router.GET("/account", h.ReadAllAccounts)
	router.GET("/account/:id", h.ReadAccount)
	router.POST("/account/:id", h.UpdateAccount)
	router.PUT("/account/:id", h.UpdateAccount)
	router.DELETE("/account/:id", h.DeleteAccount)

	return router
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("[REST] - CreateAccount()")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	var account domain.Account
	if err = json.Unmarshal(b, &account); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	if err = h.bankService.CreateAccount(r.Context(), account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("[REST] %d", http.StatusCreated)
}

func (h *Handler) ReadAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("[REST] - ReadAccount()")
	id, err := parseId(ps)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	account, err := h.bankService.ReadAccount(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	response, err := json.Marshal(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
	log.Printf("[REST] %d", http.StatusOK)
}

func (h *Handler) ReadAllAccounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("[REST] - ReadAllAccounts()")
	accounts, err := h.bankService.ReadAllAccounts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	response, err := json.Marshal(accounts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
	log.Printf("[REST] %d", http.StatusOK)
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("[REST] - UpdateAccount()")
	id, err := parseId(ps)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	var account domain.AccountUpdateInput
	if err = json.Unmarshal(b, &account); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	if err = h.bankService.UpdateAccount(r.Context(), id, account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[REST] %d", http.StatusOK)
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("[REST] - DeleteAccount()")
	id, err := parseId(ps)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	if err := h.bankService.DeleteAccount(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[REST][ERROR] %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[REST] %d", http.StatusOK)
}

func parseId(ps httprouter.Params) (int64, error) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		return id, err
	}

	if id < 1 {
		return id, errors.New("invalid Id")
	}

	return id, err
}
