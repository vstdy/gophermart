package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/vstdy/gophermart/api/rest/model"
	"github.com/vstdy/gophermart/pkg"
	"github.com/vstdy/gophermart/service/gophermart"
)

// Handler keeps handler dependencies.
type Handler struct {
	service gophermart.Service
	jwtAuth *jwtauth.JWTAuth
}

// NewHandler returns a new Handler instance.
func NewHandler(service gophermart.Service, jwtAuth *jwtauth.JWTAuth) Handler {
	return Handler{service: service, jwtAuth: jwtAuth}
}

// register registers user.
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

	if err = h.addJWTCookie(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// login authorizes user.
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

	if err = h.addJWTCookie(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// addUsersOrder adds the user's order.
func (h Handler) addUsersOrder(w http.ResponseWriter, r *http.Request) {
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

	rawObj := model.NewOrder(userID, orderID).ToCanonical()
	obj, err := h.service.AddOrder(r.Context(), rawObj)
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

// getUsersOrders returns user's orders.
func (h Handler) getUsersOrders(w http.ResponseWriter, r *http.Request) {
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

// getUsersBalance returns the user's balance.
func (h Handler) getUsersBalance(w http.ResponseWriter, r *http.Request) {
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

	balance := model.NewBalanceResponse(current, used)
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

// addWithdrawal adds user's withdrawal.
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

// getUsersWithdrawals returns user's withdrawals.
func (h Handler) getUsersWithdrawals(w http.ResponseWriter, r *http.Request) {
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
