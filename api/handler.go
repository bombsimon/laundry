package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bombsimon/laundry"
	"github.com/gorilla/mux"
)

type LaundryAPI struct {
	laundry *laundry.Laundry
}

func New(l *laundry.Laundry) *LaundryAPI {
	api := LaundryAPI{
		laundry: l,
	}

	return &api
}

func (api *LaundryAPI) GetBookers(w http.ResponseWriter, r *http.Request) {
	b, _ := api.laundry.GetBookers()

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) AddBooker(w http.ResponseWriter, r *http.Request) {
	var inRequest laundry.Booker
	if err := api.getJSONBody(&inRequest, r.Body, w, false); err != nil {
		return
	}

	b, err := api.laundry.AddBooker(&inRequest)
	if err != nil {
		api.RespondJSON(ErrorResponse{400, []string{err.Error()}}, w)
		return
	}

	jb, _ := json.Marshal(b)
	w.Write(jb)
}

func (api *LaundryAPI) GetBooker(w http.ResponseWriter, r *http.Request) {
	bookerId, _ := strconv.Atoi(mux.Vars(r)["id"])
	b, err := api.laundry.GetBooker(bookerId)
	if err != nil {
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
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	var inRequest laundry.Booker
	if err := api.getJSONBody(&inRequest, r.Body, w, false); err != nil {
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
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	if err = api.laundry.RemoveBooker(b); err != nil {
		w.Write(err.(*laundry.LaundryError).AsJSON())
		return
	}

	var empty = struct{}{}

	jb, _ := json.Marshal(&empty)
	w.Write(jb)
}

func (api *LaundryAPI) GetBookerBookings(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) GetMachines(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) AddMachine(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) GetMachine(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) UpdateMachine(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) RemoveMachine(w http.ResponseWriter, r *http.Request) {
}

func (api *LaundryAPI) GetSlots(w http.ResponseWriter, r *http.Request) {
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

func (api *LaundryAPI) getJSONBody(i interface{}, b io.ReadCloser, w http.ResponseWriter, restoreBody bool) error {
	decoder := json.NewDecoder(b)

	if err := decoder.Decode(i); err != nil {
		api.RespondJSON(ErrorResponse{400, []string{err.Error()}}, w)
		return err
	}

	return nil
}
