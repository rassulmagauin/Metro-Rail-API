# Go Metro API

This is a Go application that provides a RESTful API for managing train, station, and schedule resources in a metro system.

## Getting Started

These instructions will help you set up and run the Go Metro API on your local machine for development and testing purposes.

### Prerequisites

- Go programming language (version 1.13 or higher)
- SQLite database

### Installing

1. Clone the repository:

   ```shell
   git clone https://github.com/rassulmagauin/metro-rail-api.git
   ```

2. Navigate to the project directory:

   ```shell
   cd go-metro-api
   ```

3. Install the required dependencies:

   ```shell
   go mod download
   ```

4. Create a SQLite database file:

   ```shell
   touch railapi.db
   ```

5. Run the application:

   ```shell
   go run railAPI/main.go
   ```

The API will start running on `http://localhost:8000`.

## API Endpoints

The following API endpoints are available:

- **GET /v1/trains/{train-id}**: Get details of a specific train.
- **POST /v1/trains**: Create a new train.
- **DELETE /v1/trains/{train-id}**: Delete a train.

- **GET /v1/stations/{station-id}**: Get details of a specific station.
- **POST /v1/stations**: Create a new station.
- **DELETE /v1/stations/{station-id}**: Delete a station.

- **GET /v1/schedules/{schedule-id}**: Get details of a specific schedule.
- **POST /v1/schedules**: Create a new schedule.
- **DELETE /v1/schedules/{schedule-id}**: Delete a schedule.

## Built With

- [Go](https://golang.org/) - The programming language used
- [go-restful](https://github.com/emicklei/go-restful) - Go framework for building RESTful APIs
- [SQLite](https://www.sqlite.org/) - Database engine

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments
- [The Go Programming Language](https://golang.org/) - for the powerful Go language and its ecosystem.
- [go-restful](https://github.com/emicklei/go-restful) - for the RESTful API framework used in this project.

```

