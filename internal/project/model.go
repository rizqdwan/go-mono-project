package project

import "time"

type Project struct {
	ID            string
	ApplicationID int64
	Name          string
	Year          string
	ParentCode    string
	StartDate     time.Time
	DueDate       time.Time
	StatusID      string
	Category      string
	CreatedAt     time.Time
	UpdateAt      time.Time
}
