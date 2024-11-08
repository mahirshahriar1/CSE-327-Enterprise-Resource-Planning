package models

import "time"

// Attendance represents employee attendance
type Attendance struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
	TotalHours float64 `json:"total_hours"`
}

// AttendanceStore defines an interface for attendance-related database operations
type AttendanceStore interface {
	CreateAttendance(attendance *Attendance) error
	GetAttendanceByUserID(userID int) ([]*Attendance, error)
	UpdateAttendance(attendance *Attendance) error
	DeleteAttendance(id int) error
}
