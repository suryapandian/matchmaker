package uuid

import (
	uuid "github.com/satori/go.uuid"
)

func IsValidUUID(input string) bool {
	return !(uuid.FromStringOrNil(input) == uuid.Nil)
}

func NewUUID() string {
	return uuid.NewV4().String()
}
