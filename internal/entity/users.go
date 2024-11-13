package entity

import "time"

type Users struct {
	GUID       string
	ExternalID int64
	Username   *string
	FirstName  *string
	LastName   *string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
