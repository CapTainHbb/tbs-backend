package usecase

import "time"

type CreateProjectRequest struct {
	Name 					string
	Description 			string
	StartDate 				time.Time
	EndDate 				time.Time
	ProposedBudget 			float64
	Status 					string
	OwnerID 				int
}

type UpdateProjectRequest struct {
	ID 						int
	Name 					string
	Description 			string
	StartDate 				time.Time
	EndDate 				time.Time
	ProposedBudget 			float64
	Status 					string
	OwnerID 				int
}