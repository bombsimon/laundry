package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bombsimon/laundry"
	"github.com/bombsimon/laundry/errors"
	"github.com/gorilla/mux"
)

// LaundryAPI represents an API to the laundry service
type LaundryAPI struct {
	version string
}

// New will create a new LaundryAPI and add the passed Laundry service
// in the internal field laundry.
func New() *LaundryAPI {
	api := LaundryAPI{"v1"}

	return &api
}

// GetBookers is the HTTP handler to get bookers
func (api *LaundryAPI) GetBookers(w http.ResponseWriter, r *http.Request) {
	b, err := laundry.GetBookers()
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

// AddBooker is the HTTP handler to add a booker
func (api *LaundryAPI) AddBooker(w http.ResponseWriter, r *http.Request) {
	var inRequest laundry.Booker
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	b, err := laundry.AddBooker(&inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) GetBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := laundry.GetBooker(bookerId)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])

	var inRequest laundry.Booker
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	b, err := laundry.UpdateBooker(bookerId, &inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := laundry.GetBooker(bookerId)
	if err != nil {
		renderError(err, w)
		return
	}

	if err = laundry.RemoveBooker(b); err != nil {
		renderError(err, w)
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetBookerBookings(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := laundry.GetBooker(bookerId)
	if err != nil {
		renderError(err, w)
		return
	}

	bookings, err := laundry.GetBookerBookings(b)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(bookings)
	w.Write(jb)
}

func (api *LaundryAPI) GetMachines(w http.ResponseWriter, r *http.Request) {
	m, _ := laundry.GetMachines()

	jb, _ := json.Marshal(m)
	w.Write(jb)

}

func (api *LaundryAPI) AddMachine(w http.ResponseWriter, r *http.Request) {
	var inRequest laundry.Machine
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	if inRequest.Info == "" {
		err := errors.New("Missing info in request").WithStatus(http.StatusBadRequest)
		renderError(err, w)
		return
	}

	m, err := laundry.AddMachine(&inRequest)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)

}

func (api *LaundryAPI) GetMachine(w http.ResponseWriter, r *http.Request) {
	machineId, _ := strconv.Atoi(mux.Vars(r)["id"])

	m, err := laundry.GetMachine(machineId)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	machineId, _ := strconv.Atoi(mux.Vars(r)["id"])
	m, err := laundry.GetMachine(machineId)
	if err != nil {
		renderError(err, w)
		return
	}

	var inRequest laundry.Machine
	if err := getJSONBody(&inRequest, r.Body); err != nil {
		renderError(err, w)
		return
	}

	defer r.Body.Close()

	if inRequest.Info == "" {
		err := errors.New("Missing field info").WithStatus(http.StatusBadRequest)
		renderError(err, w)
		return
	}

	m.Info = inRequest.Info
	m.Working = inRequest.Working

	m, err = laundry.UpdateMachine(m)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveMachine(w http.ResponseWriter, r *http.Request) {
	machineId, _ := strconv.Atoi(mux.Vars(r)["id"])
	m, err := laundry.GetMachine(machineId)
	if err != nil {
		renderError(err, w)
		return
	}

	if err = laundry.RemoveMachine(m); err != nil {
		renderError(err, w)
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetSlots(w http.ResponseWriter, r *http.Request) {
	s, err := laundry.GetSlots()
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func (api *LaundryAPI) AddSlot(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) GetSlot(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) UpdateSlot(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) RemoveSlot(w http.ResponseWriter, r *http.Request) {
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

func (api *LaundryAPI) GetMonthSchedule(w http.ResponseWriter, r *http.Request) {
	start, _ := mux.Vars(r)["start"]
	end, _ := mux.Vars(r)["end"]

	s, err := laundry.GetIntervalSchedule(start, end)
	if err != nil {
		renderError(err, w)
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func getJSONBody(i interface{}, b io.ReadCloser) *errors.LaundryError {
	defer b.Close()

	decoder := json.NewDecoder(b)

	if err := decoder.Decode(i); err != nil {
		return errors.New(err).Add("Could not marshal JSON from body").WithStatus(http.StatusBadRequest)
	}

	return nil
}

func renderError(err *errors.LaundryError, w http.ResponseWriter) {
	w.WriteHeader(err.Status)
	w.Write(err.AsJSON())

}
