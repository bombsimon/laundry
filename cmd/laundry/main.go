package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin"
	"github.com/bombsimon/laundry/pkg/api"
	"github.com/bombsimon/laundry/pkg/api/middleware"
	"github.com/bombsimon/laundry/pkg/database"
	"github.com/bombsimon/laundry/pkg/laundry"

	"github.com/gorilla/mux"
)

func main() {
	var (
		httpListen = kingpin.Flag("http-listen", "Where to listen").Envar("HTTP_LISTEN").Default(":3400").String()
		mysqlHost  = kingpin.Flag("mysql-host", "MySQL host").Envar("MYSQL_HOST").Default("localhost").String()
		mysqlPort  = kingpin.Flag("mysql-port", "MySQL port").Envar("MYSQL_TCP_PORT").Default("3401").Int()
	)

	kingpin.Parse()

	l := laundry.New(database.DBConfig{
		Host:          *mysqlHost,
		Port:          *mysqlPort,
		Database:      "laundry",
		Username:      "laundry",
		Password:      "laundry",
		RetryCount:    2,
		RetryInterval: 5,
	})

	api := api.New(l)

	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	// Bookers
	v1.HandleFunc("/bookers", api.GetBookers).Name("get_bookers").Methods("GET", "OPTIONS")
	v1.HandleFunc("/bookers", api.AddBooker).Name("add_booker").Methods("POST", "OPTIONS")
	v1.HandleFunc("/bookers/{id:[0-9]+}", api.GetBooker).Name("get_booker").Methods("GET", "OPTIONS")
	v1.HandleFunc("/bookers/{id:[0-9]+}", api.UpdateBooker).Name("update_booker").Methods("PUT", "OPTIONS")
	v1.HandleFunc("/bookers/{id:[0-9]+}", api.RemoveBooker).Name("remove_booker").Methods("DELETE", "OPTIONS")
	v1.HandleFunc("/bookers/{id:[0-9]+}/bookings", api.GetBookerBookings).Name("get_booker_bookings").Methods("GET", "OPTIONS")

	// Machines
	v1.HandleFunc("/machines", api.GetMachines).Name("get_machines").Methods("GET", "OPTIONS")
	v1.HandleFunc("/machines", api.AddMachine).Name("add_machine").Methods("POST", "OPTIONS")
	v1.HandleFunc("/machines/{id:[0-9]+}", api.GetMachine).Name("get_machine").Methods("GET", "OPTIONS")
	v1.HandleFunc("/machines/{id:[0-9]+}", api.UpdateMachine).Name("update_machine").Methods("PUT", "OPTIONS")
	v1.HandleFunc("/machines/{id:[0-9]+}", api.RemoveMachine).Name("remove_machine").Methods("DELETE", "OPTIONS")

	// Slots
	v1.HandleFunc("/slots", api.GetSlots).Name("get_slots").Methods("GET", "OPTIONS")
	v1.HandleFunc("/slots", api.AddSlot).Name("add_slot").Methods("POST", "OPTIONS")
	v1.HandleFunc("/slots/{id:[0-9]+}", api.GetSlot).Name("get_slot").Methods("GET", "OPTIONS")
	v1.HandleFunc("/slots/{id:[0-9]+}", api.UpdateSlot).Name("update_slot").Methods("PUT", "OPTIONS")
	v1.HandleFunc("/slots/{id:[0-9]+}", api.RemoveSlot).Name("remove_slot").Methods("DELETE", "OPTIONS")

	// Bookings
	v1.HandleFunc("/bookings", api.GetBookings).Name("get_bookings").Methods("GET", "OPTIONS")
	v1.HandleFunc("/bookings", api.AddBooking).Name("add_booking").Methods("POST", "OPTIONS")
	v1.HandleFunc("/bookings/{id:[0-9]+}", api.GetBooking).Name("get_booking").Methods("GET", "OPTIONS")
	v1.HandleFunc("/bookings/{id:[0-9]+}", api.UpdateBooking).Name("update_booking").Methods("PUT", "OPTIONS")
	v1.HandleFunc("/bookings/{id:[0-9]+}", api.RemoveBooking).Name("remove_booking").Methods("DELETE", "OPTIONS")
	v1.HandleFunc("/bookings/{id:[0-9]+}/notifications", api.RemoveBooking).Name("get_booking_notifications").Methods("GET", "OPTIONS")

	// Schedule
	v1.HandleFunc(`/schedule/{start:\d{4}-\d{2}-\d{2}}/{end:\d{4}-\d{2}-\d{2}}`, api.GetSchedule).Name("get_month_schedule").Methods("GET", "OPTIONS")

	// Notificationos

	fmt.Println("Listening on", *httpListen, "...")

	err := http.ListenAndServe(
		*httpListen,
		middleware.Adapt(
			r,
			middleware.Logger(),
			middleware.CORS(),
		),
	)

	if err != nil {
		panic(err)
	}
}
