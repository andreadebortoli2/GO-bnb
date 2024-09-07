package dbrepo

import (
	"errors"
	"time"

	"github.com/andreadebortoli2/GO-bnb/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation insert a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction insert a restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID return true if availability exists for room id, and false if no avaliability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// if the room id is 2, then fail; otherwise pass
	if roomID == 2 {
		return false, errors.New("some error")
	}
	// if start equal "2030-01-01" return there are availability
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2030-01-01")
	if start == startDate {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	// if start and end equlas returen error
	if start == end {
		return rooms, errors.New("some error")
	}
	// if start equal "2030-01-01" return there are rooms
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2030-01-01")
	if start == startDate {
		rooms = []models.Room{{ID: 1}, {ID: 2}}
		return rooms, nil
	}
	return rooms, nil
}

// GetRoomByID gets room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}
