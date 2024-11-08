package models

import "time"

// Leave represents employee leave
type Leave struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	LeaveType string    `json:"leave_type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
}

// LeaveStore defines an interface for leave-related database operations
type LeaveStore interface {
	CreateLeave(leave *Leave) error
	GetLeaveByUserID(userID int) ([]*Leave, error)
	UpdateLeaveStatus(id int, status string) error
	DeleteLeave(id int) error
}
