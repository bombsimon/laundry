package main

import (
	"net/http"
	"os"

	"github.com/bombsimon/laundry"
	"github.com/bombsimon/laundry/api"
	"github.com/bombsimon/laundry/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		configFile = kingpin.Flag("config-file", "Path to configuration file").Envar("LAUNDRY_CONFIG_FILE").String()
	)

	kingpin.Parse()

	service := laundry.New(*configFile)
	api := api.New(service)

	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	// Bookers
	v1.HandleFunc("/bookers", api.GetBookers).Name("get_bookers").Methods("GET")
	v1.HandleFunc("/bookers", api.AddBooker).Name("add_booker").Methods("POST")
	v1.HandleFunc("/bookers/{id:[0-9]+}", api.GetBooker).Name("get_booker").Methods("GET")
	v1.HandleFunc("/bookers/{id:[0-9]+}", api.UpdateBooker).Name("update_booker").Methods("PUT")
	v1.HandleFunc("/bookers/{id:[0-9]+}", api.RemoveBooker).Name("remove_booker").Methods("DELETE")
	v1.HandleFunc("/bookers/{id:[0-9]+}/bookings", api.GetBookerBookings).Name("get_booker_bookings").Methods("GET")

	// Machines
	v1.HandleFunc("/machines", api.GetMachines).Name("get_machines").Methods("GET")
	v1.HandleFunc("/machines", api.AddMachine).Name("add_machine").Methods("POST")
	v1.HandleFunc("/machines/{id:[0-9]+}", api.GetMachine).Name("get_machine").Methods("GET")
	v1.HandleFunc("/machines/{id:[0-9]+}", api.UpdateMachine).Name("update_machine").Methods("PUT")
	v1.HandleFunc("/machines/{id:[0-9]+}", api.RemoveMachine).Name("remove_machine").Methods("DELETE")

	// Slots
	v1.HandleFunc("/slots", api.GetSlots).Name("get_slots").Methods("GET")
	v1.HandleFunc("/slots", api.AddSlot).Name("add_slot").Methods("POST")
	v1.HandleFunc("/slots/{id:[0-9]+}", api.GetSlot).Name("get_slot").Methods("GET")
	v1.HandleFunc("/slots/{id:[0-9]+}", api.UpdateSlot).Name("update_slot").Methods("PUT")
	v1.HandleFunc("/slots/{id:[0-9]+}", api.RemoveSlot).Name("remove_slot").Methods("DELETE")

	// Bookings
	v1.HandleFunc("/bookings", api.GetBookings).Name("get_bookings").Methods("GET")
	v1.HandleFunc("/bookings", api.AddBooking).Name("add_booking").Methods("POST")
	v1.HandleFunc("/bookings/{id:[0-9]+}", api.GetBooking).Name("get_booking").Methods("GET")
	v1.HandleFunc("/bookings/{id:[0-9]+}", api.UpdateBooking).Name("update_booking").Methods("PUT")
	v1.HandleFunc("/bookings/{id:[0-9]+}", api.RemoveBooking).Name("remove_booking").Methods("DELETE")
	v1.HandleFunc("/bookings/{id:[0-9]+}/notifications", api.RemoveBooking).Name("get_booking_notifications").Methods("GET")

	// Schedule
	v1.HandleFunc(`/schedule/{start:\d{4}-\d{2}-\d{2}}/{end:\d{4}-\d{2}-\d{2}}`, api.GetMonthSchedule).Name("get_month_schedule").Methods("GET")

	// Notificationos

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	service.Logger.Infof("Serving up at %s...", service.Config.HTTP.Listen)

	http.ListenAndServe(
		service.Config.HTTP.Listen,
		middleware.Adapt(
			loggedRouter,
			middleware.Notify(),
		),
	)
}
