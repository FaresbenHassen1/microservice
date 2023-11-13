package models

import (
	"github.com/google/uuid"
)

type Users struct {
	Id   uuid.UUID `JSON:"id_user"`
	Name string    `JSON:"name"`
}
