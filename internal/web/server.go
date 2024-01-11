package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
	"github.com/mkermilska/rentals-challenge/pkg/service"
)

type APIServer struct {
	port       int
	rentalSvc  service.RentalService
	logger     *zap.Logger
	httpServer *http.Server
}

func New(port int, rentalSvc *service.RentalService, logger *zap.Logger) *APIServer {
	return &APIServer{
		port:      port,
		rentalSvc: *rentalSvc,
		logger:    logger,
	}
}

func (a *APIServer) Start() {
	a.logger.Info("Starting APIServer server", zap.Int("port", a.port))
	a.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.port),
		Handler:           a.handler(),
		ReadHeaderTimeout: 2 * time.Second,
	}
	a.httpServer.ListenAndServe()
}

func (a *APIServer) handler() http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/rentals", a.getRentals)
		r.Get("/rentals/{rentalID}", a.getRentalByID)
	})

	return r
}

func (a *APIServer) getRentalByID(w http.ResponseWriter, r *http.Request) {
	rentalID, err := strconv.Atoi(chi.URLParam(r, "rentalID"))
	if err != nil {
		errorMsg := "Incorrect rental ID, please enter a valid number"
		a.logger.Error(errorMsg, zap.String("rentalID", chi.URLParam(r, "rentalID")), zap.Error(err))
		http.Error(w, errorMsg, http.StatusBadRequest)
		return

	}

	rental, err := a.rentalSvc.GetRentalByID(rentalID)
	if err != nil {
		errorMsg := "Rental not found"
		a.logger.Error(errorMsg, zap.Int("rentalID", rentalID), zap.Error(err))
		http.Error(w, errorMsg, http.StatusNotFound)
		return
	}

	out, err := json.Marshal(rental)
	if err != nil {
		errorMsg := "Error parsing rental"
		a.logger.Error(errorMsg, zap.Any("rental", rental), zap.Error(err))
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		a.logger.Error("Error writing API response", zap.Error(err))
	}
}
