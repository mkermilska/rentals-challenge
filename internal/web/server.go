package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
	apiv1 "github.com/mkermilska/rentals-challenge/api/v1"
	"github.com/mkermilska/rentals-challenge/pkg/database"
	"github.com/mkermilska/rentals-challenge/pkg/service"
	"github.com/mkermilska/rentals-challenge/pkg/utils"
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
	a.logger.Info("Starting API Server", zap.Int("port", a.port))
	a.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.port),
		Handler:           a.handler(),
		ReadHeaderTimeout: 2 * time.Second,
	}
	err := a.httpServer.ListenAndServe()
	if err != nil {
		a.logger.Error("Error starting API Sever")
	}
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

func (a *APIServer) getRentals(w http.ResponseWriter, r *http.Request) {
	//reading the input params could be simplified with using a library.
	queryParams := database.RentalParams{}
	if r.URL.Query().Has("price_min") {
		minPrice, err := strconv.Atoi(r.URL.Query().Get("price_min"))
		if err != nil {
			errorMsg := "Invalid value for price_min parameter"
			a.logger.Error(errorMsg, zap.Error(err))
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		queryParams.PriceMin = minPrice
	}

	if r.URL.Query().Has("price_max") {
		maxPrice, err := strconv.Atoi(r.URL.Query().Get("price_max"))
		if err != nil {
			errorMsg := "Invalid value for price_max parameter"
			a.logger.Error(errorMsg, zap.Error(err))
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		queryParams.PriceMax = maxPrice
	}

	if r.URL.Query().Has("ids") {
		IDs := r.URL.Query().Get("ids")
		if IDs == "" {
			errorMsg := "Empty ids parameter"
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		IDsArr := strings.Split(IDs, ",")
		for _, ID := range IDsArr {
			if _, err := strconv.Atoi(ID); err != nil {
				errorMsg := "Invalid id exists in ids parameter"
				a.logger.Error(errorMsg, zap.Error(err))
				http.Error(w, errorMsg, http.StatusBadRequest)
				return
			}
		}
		queryParams.IDs = IDsArr
	}

	if r.URL.Query().Has("near") {
		near := r.URL.Query().Get("near")
		if near == "" {
			errorMsg := "Empty near parameter"
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		nearPoint := strings.Split(near, ",")
		if len(nearPoint) != 2 {
			errorMsg := "Near parameter expects comma separated pair of float numbers"
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		latPoint, err := strconv.ParseFloat(nearPoint[0], 64)
		if err != nil {
			errorMsg := "Invalid latitude value in near parameter"
			a.logger.Error(errorMsg, zap.Error(err))
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		lngPoint, err := strconv.ParseFloat(nearPoint[1], 64)
		if err != nil {
			errorMsg := "Invalid longitude value in near parameter"
			a.logger.Error(errorMsg, zap.Error(err))
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		queryParams.Near = *utils.CalculateNearBox(utils.Point{
			Lat: latPoint,
			Lng: lngPoint,
		}, 100)
	}

	if r.URL.Query().Has("limit") {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			errorMsg := "Invalid value for limit parameter"
			a.logger.Error(errorMsg, zap.Error(err))
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		queryParams.Limit = limit
	}

	if r.URL.Query().Has("offset") {
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			errorMsg := "Invalid value for offset parameter"
			a.logger.Error(errorMsg, zap.Error(err))
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		queryParams.Offset = offset
	}

	if r.URL.Query().Has("sort") {
		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" {
			errorMsg := "Empty sort parameter"
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		if _, exists := apiv1.SortsMap[sortBy]; !exists {
			errorMsg := "Sort by given column is not allowed"
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		queryParams.Sort = apiv1.SortsMap[sortBy]
	}

	rentals, err := a.rentalSvc.GetRentals(queryParams)
	if err != nil {
		errorMsg := "Error getting rentals"
		a.logger.Error(errorMsg, zap.Error(err))
		http.Error(w, errorMsg, http.StatusNotFound)
		return
	}

	out, err := json.Marshal(rentals)
	if err != nil {
		errorMsg := "Error parsing rentals"
		a.logger.Error(errorMsg, zap.Error(err))
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		a.logger.Error("Error writing API response", zap.Error(err))
	}
}
