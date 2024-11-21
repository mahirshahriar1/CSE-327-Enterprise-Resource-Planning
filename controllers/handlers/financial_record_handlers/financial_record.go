// Package financial_record_handlers provides HTTP handlers for managing financial records.
package financial_record_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"erp/models"

	"github.com/gorilla/mux"
)

// FinancialRecordHandler struct handles HTTP requests for managing financial records.
type FinancialRecordHandler struct {
	RecordStore models.FinancialRecordStore
}

// RegisterRoutes registers the routes for financial record handlers.
func RegisterRoutes(router *mux.Router, recordStore models.FinancialRecordStore) {
	handler := &FinancialRecordHandler{RecordStore: recordStore}

	// Register routes for financial records
	router.HandleFunc("/records", handler.CreateRecord).Methods("POST")
	router.HandleFunc("/records/{id:[0-9]+}", handler.GetRecord).Methods("GET")
	router.HandleFunc("/records/{id:[0-9]+}", handler.UpdateRecord).Methods("PUT")
	router.HandleFunc("/records/{id:[0-9]+}", handler.DeleteRecord).Methods("DELETE")
}

// CreateRecord creates a new financial record.
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

// GetRecord retrieves a financial record by ID.
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

// UpdateRecord updates an existing financial record.
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

// DeleteRecord deletes a financial record by ID.
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
