package repository

import "github.com/andreadebortoli2/GO-bnb/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(r models.RoomRestriction) error
}
