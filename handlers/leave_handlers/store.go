package leave_handlers

import (
	"database/sql"
	"erp/models"
)

// DBLeaveStore implements the LeaveStore interface for a SQL database
type DBLeaveStore struct {
	DB *sql.DB
}

// CreateLeave creates a new leave request in the database
func (s *DBLeaveStore) CreateLeave(leave *models.Leave) error {
	query := `INSERT INTO leave (user_id, leave_type, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.DB.QueryRow(query, leave.UserID, leave.LeaveType, leave.StartDate, leave.EndDate, leave.Status).Scan(&leave.ID)
	return err
}

// GetLeaveByUserID fetches all leave requests for a given user
func (s *DBLeaveStore) GetLeaveByUserID(userID int) ([]*models.Leave, error) {
	query := `SELECT id, user_id, leave_type, start_date, end_date, status FROM leave WHERE user_id = $1`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []*models.Leave
	for rows.Next() {
		leave := &models.Leave{}
		if err := rows.Scan(&leave.ID, &leave.UserID, &leave.LeaveType, &leave.StartDate, &leave.EndDate, &leave.Status); err != nil {
			return nil, err
		}
		leaves = append(leaves, leave)
	}
	return leaves, nil
}

// UpdateLeaveStatus updates the status of a leave request
func (s *DBLeaveStore) UpdateLeaveStatus(id int, status string) error {
	query := `UPDATE leave SET status = $1 WHERE id = $2`
	_, err := s.DB.Exec(query, status, id)
	return err
}
