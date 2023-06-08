package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"metro/dbutils"
	"net/http"

	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

// DB Driver visible to whole program

var DB *sql.DB

//TrainResource is the model for holding rail information

type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

//StationResource contains information about locations

type StationResource struct {
	ID          int
	Name        string
	OpeningTime string
	ClosingTime string
}

//ScheduleResource links both trains ans stations

type ScheduleResource struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime string
}

//Register adds pathes and routes to container

func (t *TrainResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/v1/trains").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("{train-id}").To(t.removeTrain))

	container.Add(ws)
}

//GET http://localhost:8000/v1/trains/1

func (t TrainResource) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	err := DB.QueryRow("SELECT ID, DRIVER_NAME, OPERATING_STATUS FROM train where id=?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Train could not be found")
	} else {
		response.WriteEntity(t)
	}
}

// POST http://localhost:800/v1/trains
func (t TrainResource) createTrain(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b TrainResource
	err := decoder.Decode(&b)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Cannot decode!")
	} else {
		log.Println(b.DriverName, b.OperatingStatus)
		statement, _ := DB.Prepare("insert into train (DRIVER_NAME, OPERATING_STATUS) values (?, ?)")
		result, err := statement.Exec(b.DriverName, b.OperatingStatus)

		if err == nil {
			newID, _ := result.LastInsertId()
			b.ID = int(newID)
			response.WriteHeaderAndEntity(http.StatusCreated, b)
		} else {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
		}
	}

}

// DELETE http://localhost:8000/trains/1
func (t TrainResource) removeTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	statement, _ := DB.Prepare("delete from train where id = ?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (s *StationResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/stations").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{station-id}").To(s.getStation))
	ws.Route(ws.POST("").To(s.createStation))
	ws.Route(ws.DELETE("/{station-id}").To(s.removeStation))

	container.Add(ws)
}
func (s StationResource) getStation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("station-id")
	err := DB.QueryRow("SELECT ID, NAME, OPENING_TIME, CLOSING_TIME FROM station WHERE id=?", id).Scan(&s.ID, &s.Name, &s.OpeningTime, &s.ClosingTime)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Station could not be found!")
	} else {
		response.WriteEntity(s)
	}
}
func (s StationResource) createStation(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b StationResource
	err := decoder.Decode(&b)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error in decoding!")
		return
	}
	statement, _ := DB.Prepare("insert into station (name, opening_time, closing_time) values (?, ?, ?)")
	result, err := statement.Exec(b.Name, b.OpeningTime, b.ClosingTime)
	if err == nil {
		newID, _ := result.LastInsertId()
		b.ID = int(newID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}

}
func (s StationResource) removeStation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("station-id")
	statement, _ := DB.Prepare("delete from station where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// Register adds paths and routes to the container
func (s *ScheduleResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/v1/schedules").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{schedule-id}").To(s.getSchedule))
	ws.Route(ws.POST("").To(s.createSchedule))
	ws.Route(ws.DELETE("/{schedule-id}").To(s.removeSchedule))

	container.Add(ws)
}

// GET http://localhost:8000/v1/schedules/1
func (s ScheduleResource) getSchedule(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("schedule-id")
	err := DB.QueryRow("SELECT ID, TRAIN_ID, STATION_ID, ARRIVAL_TIME FROM schedule WHERE id=?", id).
		Scan(&s.ID, &s.TrainID, &s.StationID, &s.ArrivalTime)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Schedule could not be found!")
	} else {
		response.WriteEntity(s)
	}
}

// POST http://localhost:800/v1/schedules
func (s ScheduleResource) createSchedule(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b ScheduleResource
	err := decoder.Decode(&b)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error in decoding!")
		return
	}
	statement, _ := DB.Prepare("INSERT INTO schedule (TRAIN_ID, STATION_ID, ARRIVAL_TIME) VALUES (?, ?, ?)")
	result, err := statement.Exec(b.TrainID, b.StationID, b.ArrivalTime)
	if err == nil {
		newID, _ := result.LastInsertId()
		b.ID = int(newID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// DELETE http://localhost:8000/schedules/1
func (s ScheduleResource) removeSchedule(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("schedule-id")
	statement, _ := DB.Prepare("DELETE FROM schedule WHERE id=?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}
func main() {
	var err error
	DB, err = sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("Driver creation failed")
	}
	dbutils.Initialize(DB)
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	t := TrainResource{}
	s := StationResource{}
	sch := ScheduleResource{}
	sch.Register(wsContainer)
	t.Register(wsContainer)
	s.Register(wsContainer)
	log.Printf("start listening on localhost:8000")
	server := &http.Server{Addr: ":8000", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
