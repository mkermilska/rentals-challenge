package mapper

import (
	apiv1 "github.com/mkermilska/rentals-challenge/api/v1"
	"github.com/mkermilska/rentals-challenge/pkg/database"
)

func RentalToAPIRental(rental database.Rental) *apiv1.Rental {
	return &apiv1.Rental{
		ID:              rental.ID,
		Name:            rental.Name,
		Description:     rental.Description,
		Type:            rental.Type,
		Make:            rental.VehicleMake,
		Model:           rental.VehicleModel,
		Year:            rental.VehicleYear,
		Length:          rental.VehicleLength,
		Sleeps:          rental.Sleeps,
		PrimaryImageURL: rental.PrimaryImageURL,
		Price: apiv1.Price{
			Day: rental.PricePerDay,
		},
		Location: apiv1.Location{
			City:    rental.HomeCity,
			State:   rental.HomeState,
			Zip:     rental.HomeZip,
			Country: rental.HomeCountry,
			Lat:     rental.Lat,
			Lng:     rental.Lng,
		},
		User: apiv1.User{
			ID:        rental.User.ID,
			FirstName: rental.User.FirstName,
			LastName:  rental.User.LastName,
		},
	}
}

func RentalsToAPIRentals(rentals []database.Rental) []apiv1.Rental {
	apiRentals := make([]apiv1.Rental, len(rentals))
	for i, r := range rentals {
		rental := r
		apiRental := RentalToAPIRental(rental)
		apiRentals[i] = *apiRental
	}
	return apiRentals
}
