package entity

import "time"

type FilterParameters struct {
	Name          *string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
	UpdatedAtFrom *time.Time
	UpdatedAtTo   *time.Time
}
