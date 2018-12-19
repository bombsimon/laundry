package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bombsimon/laundry/pkg"
	"github.com/bombsimon/laundry/pkg/api/middleware"
	"github.com/bombsimon/laundry/pkg/laundry"
	"github.com/gorilla/mux"
)

// LaundryAPI represents an API to the laundry service
type LaundryAPI struct {
	booking pkg.BookingHandler
	machine pkg.MachineHandler
	slot    pkg.SlotHandler
}

// New will create a new LaundryAPI and add the passed Laundry service in the
// internal field laundry.
func New(laundryService *laundry.Laundry) *LaundryAPI {
	api := LaundryAPI{
		booking: laundryService,
		machine: laundryService,
		slot:    laundryService,
	}

	return &api
}

// GetBookers is the HTTP handler to get bookers
func (api *LaundryAPI) GetBookers(w http.ResponseWriter, r *http.Request) {
	b, err := api.booking.GetBookers()
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

// AddBooker is the HTTP handler to add a booker
func (api *LaundryAPI) AddBooker(w http.ResponseWriter, r *http.Request) {
	var inRequest pkg.Booker

	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	b, err := api.booking.AddBooker(&inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) GetBooker(w http.ResponseWriter, r *http.Request) {
	bookerID, _ := strconv.Atoi(mux.Vars(r)["id"])

	b, err := api.booking.GetBooker(bookerID)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateBooker(w http.ResponseWriter, r *http.Request) {
	bookerID, _ := strconv.Atoi(mux.Vars(r)["id"])

	var inRequest pkg.Booker
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	b, err := api.booking.UpdateBooker(bookerID, &inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveBooker(w http.ResponseWriter, r *http.Request) {
	bookerID, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := api.booking.RemoveBookerByID(bookerID); err != nil {
		renderError(err, w)
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetBookerBookings(w http.ResponseWriter, r *http.Request) {
	bookerID, _ := strconv.Atoi(mux.Vars(r)["id"])

	bookings, err := api.booking.GetBookerBookingsByID(bookerID)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(bookings)
	w.Write(jb)
}

func (api *LaundryAPI) GetMachines(w http.ResponseWriter, r *http.Request) {
	m, _ := api.machine.GetMachines()

	jb, _ := json.Marshal(m)
	w.Write(jb)

}

func (api *LaundryAPI) AddMachine(w http.ResponseWriter, r *http.Request) {
	var inRequest pkg.Machine
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	m, err := api.machine.AddMachine(&inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)

}

func (api *LaundryAPI) GetMachine(w http.ResponseWriter, r *http.Request) {
	machineID, _ := strconv.Atoi(mux.Vars(r)["id"])

	m, err := api.machine.GetMachine(machineID)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	machineID, _ := strconv.Atoi(mux.Vars(r)["id"])

	var inRequest pkg.Machine
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	m, err := api.machine.UpdateMachine(machineID, &inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveMachine(w http.ResponseWriter, r *http.Request) {
	machineID, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := api.machine.RemoveMachineByID(machineID); err != nil {
		renderError(err, w)
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetSlots(w http.ResponseWriter, r *http.Request) {
	s, err := api.slot.GetSlots()
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func (api *LaundryAPI) AddSlot(w http.ResponseWriter, r *http.Request) {
	var inRequest pkg.Slot
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	s, err := api.slot.AddSlot(&inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func (api *LaundryAPI) GetSlot(w http.ResponseWriter, r *http.Request) {
	slotID, _ := strconv.Atoi(mux.Vars(r)["id"])

	s, err := api.slot.GetSlot(slotID)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateSlot(w http.ResponseWriter, r *http.Request) {
	slotID, _ := strconv.Atoi(mux.Vars(r)["id"])

	var inRequest pkg.Slot
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	s, err := api.slot.UpdateSlot(slotID, &inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveSlot(w http.ResponseWriter, r *http.Request) {
	slotID, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := api.slot.RemoveSlotByID(slotID); err != nil {
		renderError(err, w)
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetBookings(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) AddBooking(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) GetBooking(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) UpdateBooking(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) RemoveBooking(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) GetSchedule(w http.ResponseWriter, r *http.Request) {
	start, _ := mux.Vars(r)["start"]
	end, _ := mux.Vars(r)["end"]

	format := "2006-01-02"
	s, _ := time.Parse(format, start)
	e, _ := time.Parse(format, end)

	if s.After(e) {
		renderError(errors.New("Start time cannot be after end time"), w)
	}

	schedule, err := api.slot.GetSchedule(s, e)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(schedule)
	w.Write(jb)
}

func getJSONBody(i interface{}, b io.ReadCloser) error {
	defer b.Close()

	decoder := json.NewDecoder(b)

	if err := decoder.Decode(i); err != nil {
		return err
	}

	return nil
}

func renderError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"errors":["` + err.Error() + `"]}`))

	if lrw, ok := w.(*middleware.LoggingResponseWriter); ok {
		lrw.WriteError(err)
	}
}
