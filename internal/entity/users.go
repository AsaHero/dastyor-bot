package entity

import "time"

type Users struct {
	ID        int64
	Username  *string
	FirstName *string
	LastName  *string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
