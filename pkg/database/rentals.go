package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	apiv1 "github.com/mkermilska/rentals-challenge/api/v1"
)

type Rental struct {
	ID              int        `db:"id"`
	UserID          int        `db:"user_id"`
	Name            string     `db:"name"`
	Type            string     `db:"type"`
	Description     string     `db:"description"`
	Sleeps          int        `db:"sleeps"`
	PricePerDay     int        `db:"price_per_day"`
	HomeCity        string     `db:"home_city"`
	HomeState       string     `db:"home_state"`
	HomeZip         string     `db:"home_zip"`
	HomeCountry     string     `db:"home_country"`
	VehicleMake     string     `db:"vehicle_make"`
	VehicleModel    string     `db:"vehicle_model"`
	VehicleYear     int        `db:"vehicle_year"`
	VehicleLength   float32    `db:"vehicle_length"`
	Created         time.Time  `db:"created"`
	Updated         time.Time  `db:"updated"`
	Lat             float32    `db:"lat"`
	Lng             float32    `db:"lng"`
	PrimaryImageURL string     `db:"primary_image_url"`
	User            apiv1.User `db:"user"`
}

type RentalsRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRentalsRepository(db *sqlx.DB, logger *zap.Logger) *RentalsRepository {
	return &RentalsRepository{
		db:     db,
		logger: logger,
	}
}

func (rr *RentalsRepository) FindRentalByID(rentalID int) (*Rental, error) {
	rr.logger.Debug("Getting rental by ID", zap.Int("rentalID", rentalID))
	rental := Rental{}
	err := rr.db.Get(&rental,
		`SELECT r.*,
		u.id as "user.id",
		u.first_name as "user.first_name",
		u.last_name as "user.last_name"
		FROM rentals r
		JOIN users u ON r.user_id = u.id
		WHERE r.id = $1`, rentalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, fmt.Sprintf("not found rentals with id %d", rentalID))
			//return nil, fmt.Errorf("not found rentals with ID %d: %w", rentalID, err)
		}
		return nil, errors.Wrap(err, fmt.Sprintf("error getting rental with id %d", rentalID))
	}
	return &rental, nil
}

