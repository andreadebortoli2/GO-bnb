package dbrepo

import (
	"errors"
	"time"

	"github.com/andreadebortoli2/GO-bnb/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

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

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	return reservation, nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

func (m *testDBRepo) InsertBlockForRoom(id int, start time.Time) error {
	return nil
}

func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
