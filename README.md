# Rental API

API providing information for available rentals in JSON format. It returns a single rental or list of rentals that can be filtered, sorted, and paginated.

## Endpoints

- `v1/rentals/<RENTAL_ID>` Read one rental endpoint. Return the following status codes:
    - 200 (OK) on successful request
    - 404 (error not found) rental not found
    - 400 (bad request) incorrect rental id, only numbers are accepted
- `v1/rentals` Read many (list) rentals endpoint
    - Supported query parameters
        - price_min (number)
        - price_max (number)
        - limit (number)
        - offset (number)
        - ids (comma separated list of rental ids)
        - near (comma separated pair [lat,lng]) - retrieve all rentals within 100 miles around the given point
        - sort (string) - rentals could be sorted by one of the fields existing in the response structure. Any other string is considered as not valid. 
    - Examples:
        - `rentals?price_min=9000&price_max=75000`
        - `rentals?limit=3&offset=6&sort=price`
        - `rentals?ids=3,4,5`
        - `rentals?near=33.64,-117.93`
        - `rentals?near=33.64,-117.93&price_min=9000&price_max=75000&limit=3&offset=6&sort=price`
    - Status codes:
        - 200 (OK) on successful request
        - 400 (bad request) on incorrect query parameters

The rental object JSON response structure:
```json
{
  "id": "int",
  "name": "string",
  "description": "string",
  "type": "string",
  "make": "string",
  "model": "string",
  "year": "int",
  "length": "decimal",
  "sleeps": "int",
  "primary_image_url": "string",
  "price": {
    "day": "int"
  },
  "location": {
    "city": "string",
    "state": "string",
    "zip": "string",
    "country": "string",
    "lat": "decimal",
    "lng": "decimal"
  },
  "user": {
    "id": "int",
    "first_name": "string",
    "last_name": "string"
  }
}
```

## Usage
### Prerequisits
- Docker
- Docker Compose

### Run the application
- For starting the API along with postgres database, go to the home directory of the project and run:
```
make start
``` 
- Endpoints could be accessed through localhost on the following url:
```
http://localhost:59191/
```

### Tests

#### Manual tests 
Set of API requests is prepared in `tests/test-requests.http`. The requests can be executed directly from the file using `Visual Studio Code` and `REST Client` extension. This is an easy option for manual testing.

#### Unit tests
Run `make unit-tests` command for starting the unit tests.

#### Integration tests
Run `make integration-tests` command for staring Venom integration tests. Integration tests require an already started application (with `make start`).

*Note: These commands could be combined in a single target as future improvement.


### Stop the application
`make stop` stops all running containers and clean up the resources created by run the application.
