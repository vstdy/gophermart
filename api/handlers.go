package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lestrrat-go/jwx/jwa"

	"github.com/vstdy0/go-diploma/api/model"
	"github.com/vstdy0/go-diploma/cmd/gophermart/cmd/common"
	"github.com/vstdy0/go-diploma/pkg"
	"github.com/vstdy0/go-diploma/service/gophermart"
)

// Handler keeps handler dependencies.
type Handler struct {
	service   gophermart.Service
	config    common.Config
	tokenAuth *jwtauth.JWTAuth
}

// NewHandler returns a new Handler instance.
func NewHandler(service gophermart.Service, config common.Config) Handler {
	tokenAuth := jwtauth.New(jwa.HS256.String(), []byte(config.SecretKey), nil)

	return Handler{service: service, config: config, tokenAuth: tokenAuth}
}

func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var bodyObj model.RegisterBody
	err := json.NewDecoder(r.Body).Decode(&bodyObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rawObj := bodyObj.ToCanonical()

	obj, err := h.service.CreateUser(r.Context(), rawObj)
	if err != nil {
		if errors.Is(err, pkg.ErrAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.setAuthCookie(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var bodyObj model.RegisterBody
	err := json.NewDecoder(r.Body).Decode(&bodyObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rawObj := bodyObj.ToCanonical()

	obj, err := h.service.AuthenticateUser(r.Context(), rawObj)
	if err != nil {
		if errors.Is(err, pkg.ErrWrongCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.setAuthCookie(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) addUserOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	orderID := string(body)

	obj, err := h.addOrder(r.Context(), userID, orderID)
	if err != nil {
		if errors.Is(err, pkg.ErrAlreadyExists) && obj.UserID == userID {
			http.Error(w, "order number has been already uploaded", http.StatusOK)
			return
		}
		if errors.Is(err, pkg.ErrAlreadyExists) {
			http.Error(w, "order number has been already uploaded by another user", http.StatusConflict)
			return
		}
		if errors.Is(err, pkg.ErrInvalidInput) {
			http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h Handler) getUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	objs, err := h.service.GetOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if objs == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	orders := model.NewOrdersFromCanonical(objs)

	res, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) getUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	current, used, err := h.service.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance := model.BalanceResponse{
		Current:   current,
		Withdrawn: used,
	}

	res, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) addWithdrawal(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var bodyObj model.AddWithdrawalBody
	err = json.NewDecoder(r.Body).Decode(&bodyObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rawObj := bodyObj.ToCanonical(userID)

	err = h.service.AddWithdrawal(r.Context(), rawObj)
	if err != nil {
		if errors.Is(err, pkg.ErrNonSufficientFunds) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, pkg.ErrInvalidInput) {
			http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, pkg.ErrAlreadyExists) {
			http.Error(w, "withdrawal has been already added", http.StatusOK)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) getUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	objs, err := h.service.GetWithdrawals(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if objs == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	withdrawals := model.NewGetWithdrawalsFromCanonical(objs)

	res, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
