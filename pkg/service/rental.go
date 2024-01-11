package service

import (
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"

	apiv1 "github.com/mkermilska/rentals-challenge/api/v1"
	"github.com/mkermilska/rentals-challenge/pkg/database"
	"github.com/mkermilska/rentals-challenge/pkg/mapper"
)

type RentalService struct {
	rentalsRepository *database.RentalsRepository
	logger            zap.Logger
}

func NewRentalService(db *sqlx.DB, logger *zap.Logger) *RentalService {
	rentalsRepository := database.NewRentalsRepository(db, logger)

	return &RentalService{
		rentalsRepository: rentalsRepository,
		logger:            *logger,
	}
}

func (r *RentalService) GetRentalByID(rentalID int) (*apiv1.Rental, error) {
	rental, err := r.rentalsRepository.FindRentalByID(rentalID)
	if err != nil {
		r.logger.Error("Error getting rental by ID", zap.Error(err))
		return nil, err
	}
	apiRental := mapper.RentalToAPIRental(*rental)
	return apiRental, nil
}

