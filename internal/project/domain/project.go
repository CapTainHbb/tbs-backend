package domain

import (
	"time"
)


type Project struct {
	ID 					int
	Name				string
	Description			string
	StartDate			time.Time
	EndDate				time.Time
	OwnerID				int
	ProposedBudget		float64
	Status				string
}