package api

import (
	"net/http"

	"github.com/vstdy/gophermart/cmd/gophermart/cmd/common"
	"github.com/vstdy/gophermart/service/gophermart/v1"
)

// NewServer returns server.
func NewServer(svc *gophermart.Service, config common.Config) *http.Server {
	router := NewRouter(svc, config)

	return &http.Server{Addr: config.RunAddress, Handler: router}
}
