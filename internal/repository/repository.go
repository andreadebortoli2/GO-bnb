package repository

import (
	"time"

	"github.com/andreadebortoli2/GO-bnb/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	// reservation funcs
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	AllRooms() ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	// authentication funcs
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	// admin reservation funcs
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRoom(id int, start time.Time) error
	DeleteBlockByID(id int) error
}
