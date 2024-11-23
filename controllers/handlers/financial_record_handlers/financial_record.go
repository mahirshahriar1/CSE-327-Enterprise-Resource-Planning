// Package financial_record_handlers provides HTTP handlers for managing financial records.
// It enables CRUD operations on financial records through HTTP endpoints.
package financial_record_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"erp/models"

	"github.com/gorilla/mux"
)

// FinancialRecordHandler handles HTTP requests for managing financial records.
// It uses a RecordStore to interact with the database for CRUD operations.
type FinancialRecordHandler struct {
	RecordStore models.FinancialRecordStore // RecordStore is the interface for managing financial records in the database.
}

// RegisterRoutes registers the HTTP routes for managing financial records.
//
// Parameters:
//   - router: A Gorilla Mux router where the routes will be registered.
//   - recordStore: An implementation of the FinancialRecordStore interface for database operations.
func RegisterRoutes(router *mux.Router, recordStore models.FinancialRecordStore) {
	handler := &FinancialRecordHandler{RecordStore: recordStore}

	// Register routes for financial record management
	router.HandleFunc("/records", handler.CreateRecord).Methods("POST")
	router.HandleFunc("/records/{id:[0-9]+}", handler.GetRecord).Methods("GET")
	router.HandleFunc("/records/{id:[0-9]+}", handler.UpdateRecord).Methods("PUT")
	router.HandleFunc("/records/{id:[0-9]+}", handler.DeleteRecord).Methods("DELETE")
}

// CreateRecord handles HTTP POST requests to create a new financial record.
//
// HTTP Method: POST
// URL Path: /records
//
// Request Body:
//   - JSON representation of a FinancialRecord object (excluding the ID).
//
// Response:
//   - Status Code: 201 (Created) if the record is successfully created.
//   - JSON representation of the created record on success.
//   - Status Code: 400 (Bad Request) if the input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the record creation fails.
func (h *FinancialRecordHandler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	var record models.FinancialRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := h.RecordStore.CreateFinancialRecord(&record); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create financial record: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetRecord handles HTTP GET requests to retrieve a financial record by its ID.
//
// HTTP Method: GET
// URL Path: /records/{id}
//
// Response:
//   - Status Code: 200 (OK) with the record data in JSON format if found.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 404 (Not Found) if the record with the specified ID does not exist.
func (h *FinancialRecordHandler) GetRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	record, err := h.RecordStore.GetFinancialRecordByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Record not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateRecord handles HTTP PUT requests to update an existing financial record.
//
// HTTP Method: PUT
// URL Path: /records/{id}
//
// Request Body:
//   - JSON representation of a FinancialRecord object (excluding the ID, which is taken from the URL).
//
// Response:
//   - Status Code: 200 (OK) with the updated record data in JSON format if successful.
//   - Status Code: 400 (Bad Request) if the ID or input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the update operation fails.
func (h *FinancialRecordHandler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	var record models.FinancialRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	record.ID = id
	if err := h.RecordStore.UpdateFinancialRecord(&record); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update record: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteRecord handles HTTP DELETE requests to delete a financial record by its ID.
//
// HTTP Method: DELETE
// URL Path: /records/{id}
//
// Response:
//   - Status Code: 204 (No Content) if the record is successfully deleted.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 500 (Internal Server Error) if the deletion operation fails.
func (h *FinancialRecordHandler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	if err := h.RecordStore.DeleteFinancialRecord(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete record: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
