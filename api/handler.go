package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bombsimon/laundry"
	"github.com/gorilla/mux"
)

// LaundryAPI represents an API to the laundry service
type LaundryAPI struct {
	laundry *laundry.Laundry
}

// New will create a new LaundryAPI and add the passed Laundry service
// in the internal field laundry.
func New(l *laundry.Laundry) *LaundryAPI {
	api := LaundryAPI{
		laundry: l,
	}

	return &api
}

// GetBookers is the HTTP handler to get bookers
func (api *LaundryAPI) GetBookers(w http.ResponseWriter, r *http.Request) {
	b, _ := api.laundry.GetBookers()

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

// AddBooker is the HTTP handler to add a booker
func (api *LaundryAPI) AddBooker(w http.ResponseWriter, r *http.Request) {
	var inRequest laundry.Booker
	if err := api.getJSONBody(&inRequest, r.Body); err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	if inRequest.Identifier == "" {
		err := laundry.NewError("Missing identifier in request").WithStatus(http.StatusBadRequest)
		w.WriteHeader(err.Status)
		w.Write(err.AsJSON())
		return
	}

	b, err := api.laundry.AddBooker(&inRequest)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) GetBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := api.laundry.GetBooker(bookerId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := api.laundry.GetBooker(bookerId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	var inRequest laundry.Booker
	if err := api.getJSONBody(&inRequest, r.Body); err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	defer r.Body.Close()

	b.Phone = inRequest.Phone
	b.Email = inRequest.Email

	b, _ = api.laundry.UpdateBooker(b)

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := api.laundry.GetBooker(bookerId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	if err = api.laundry.RemoveBooker(b); err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetBookerBookings(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := api.laundry.GetBooker(bookerId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	bookings, err := api.laundry.GetBookerBookings(b)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(bookings)
	w.Write(jb)
}

func (api *LaundryAPI) GetMachines(w http.ResponseWriter, r *http.Request) {
	m, _ := api.laundry.GetMachines()

	jb, _ := json.Marshal(m)
	w.Write(jb)

}

func (api *LaundryAPI) AddMachine(w http.ResponseWriter, r *http.Request) {
	var inRequest laundry.Machine
	if err := api.getJSONBody(&inRequest, r.Body); err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	if inRequest.Info == "" {
		err := laundry.NewError("Missing info in request").WithStatus(http.StatusBadRequest)
		w.WriteHeader(err.Status)
		w.Write(err.AsJSON())
		return
	}

	m, err := api.laundry.AddMachine(&inRequest)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)

}

func (api *LaundryAPI) GetMachine(w http.ResponseWriter, r *http.Request) {
	machineId, _ := strconv.Atoi(mux.Vars(r)["id"])
	m, err := api.laundry.GetMachine(machineId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)
}

func (api *LaundryAPI) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	machineId, _ := strconv.Atoi(mux.Vars(r)["id"])
	m, err := api.laundry.GetMachine(machineId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	var inRequest laundry.Machine
	if err := api.getJSONBody(&inRequest, r.Body); err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	defer r.Body.Close()

	if inRequest.Info == "" {
		err := laundry.NewError("Missing field info").WithStatus(http.StatusBadRequest)
		w.WriteHeader(err.Status)
		w.Write(err.AsJSON())
		return
	}

	m.Info = inRequest.Info
	m.Working = inRequest.Working

	m, err = api.laundry.UpdateMachine(m)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(m)
	w.Write(jb)
}

func (api *LaundryAPI) RemoveMachine(w http.ResponseWriter, r *http.Request) {
	machineId, _ := strconv.Atoi(mux.Vars(r)["id"])
	m, err := api.laundry.GetMachine(machineId)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	if err = api.laundry.RemoveMachine(m); err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetSlots(w http.ResponseWriter, r *http.Request) {
	s, err := api.laundry.GetSlots()
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
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

	s, err := api.laundry.GetIntervalSchedule(start, end)
	if err != nil {
		w.WriteHeader(err.(*laundry.LaundryError).Status)
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	jb, _ := json.Marshal(s)
	w.Write(jb)
}

func (api *LaundryAPI) getJSONBody(i interface{}, b io.ReadCloser) error {
	decoder := json.NewDecoder(b)

	if err := decoder.Decode(i); err != nil {
		return laundry.ExtError(err).WithStatus(400)
	}

	return nil
}
