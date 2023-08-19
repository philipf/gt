// A person is the entity for which we associate their Calendars, Working days, Holidays and preferred working hours.
package domain

import "github.com/google/uuid"

type Person struct {
	ID         uuid.UUID
	Name       string
	Email      string
	ExternalId string // This is the ID of the person in the external system (e.g. Google). More thinking is needed as it is possible that a person may have multiple external IDs.
	IsMe       bool   // This is a flag to indicate if this is the person who is currently logged in.
}
