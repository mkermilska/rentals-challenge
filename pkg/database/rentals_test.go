package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRentalsRepository_FindRentals(t *testing.T) {
	tests := map[string]struct {
		params        RentalParams
		expectedCount int
	}{
		"Find all rentals": {
			params:        RentalParams{},
			expectedCount: 30,
		},
		"Filter by ids": {
			params: RentalParams{
				IDs: []string{"1", "2"},
			},
			expectedCount: 2,
		},
		// more tests needs to be added
	}

	rentalsRepository := NewRentalsRepository(db, zap.NewNop())
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rentals, err := rentalsRepository.FindRentals(test.params)
			require.Nil(t, err, "Error getting rentals")
			assert.Len(t, rentals, test.expectedCount)
		})
	}
}

func TestRentalsRepository_FindRentalByID(t *testing.T) {
	tests := map[string]struct {
		ID             int
		expectedExists bool
		expectedError  bool
	}{
		"Existing rental": {
			ID:             30,
			expectedExists: true,
			expectedError:  false,
		},
		"Not found rental": {
			ID:             3000,
			expectedExists: false,
			expectedError:  true,
		},
	}

	rentalsRepository := NewRentalsRepository(db, zap.NewNop())

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rental, err := rentalsRepository.FindRentalByID(test.ID)
			if err != nil && !test.expectedError {
				t.Fatalf("Test case %s failed: %s", name, err)
			}
			if err == nil && test.expectedError {
				t.Fatalf("Test case %s failed: %s", name, err)
			}
			if rental == nil && test.expectedExists {
				t.Fatalf("Test case %s failed: %s", name, err)
			}
			if rental != nil && !test.expectedExists {
				t.Fatalf("Test case %s failed: %s", name, err)
			}
		})
	}
}
