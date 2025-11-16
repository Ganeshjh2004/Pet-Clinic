package handlers

import (
	"encoding/json"
	"net/http"
	"petclinic/models"
	"petclinic/utils"
	"strconv"
)

func AppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	utils.Info("Received %s request at %s", r.Method, r.URL.Path)

	claimsVal := r.Context().Value("userClaims")
	claims, ok := claimsVal.(*utils.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if claims.Role == "owner" {
			apts := models.GetAppointmentsByOwnerID(claims.UserID)
			json.NewEncoder(w).Encode(apts)
		} else {
			apts := models.GetAllAppointments()
			if apts == nil {
				utils.Warn("No appointments found")
			}
			json.NewEncoder(w).Encode(apts)
		}

	case http.MethodPost:
		if claims.Role != "staff" && claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var appointment models.Appointment
		err := json.NewDecoder(r.Body).Decode(&appointment)
		if err != nil {
			utils.Error("Failed to decode POST body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Set OwnerID to logged in user if you want owners to create own appointments, adjust as needed
		if claims.Role == "owner" {
			appointment.OwnerID = claims.UserID
		}

		err = models.AddAppointment(appointment)
		if err != nil {
			utils.Error("Database error: %v", err)
			http.Error(w, "Failed to create appointment", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Warn("Invalid appointment ID provided for update")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Authorization check: staff/admin or owner of appointment
		existingAppointment := models.GetAppointmentByID(id)
		if existingAppointment == nil {
			http.Error(w, "Appointment not found", http.StatusNotFound)
			return
		}
		if claims.Role != "staff" && claims.Role != "admin" && existingAppointment.OwnerID != claims.UserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var appointment models.Appointment
		err = json.NewDecoder(r.Body).Decode(&appointment)
		if err != nil {
			utils.Error("Failed to decode PUT body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = models.UpdateAppointment(id, appointment)
		if err != nil {
			utils.Error("Failed to update appointment: %v", err)
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Warn("Invalid appointment ID for deletion")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		existingAppointment := models.GetAppointmentByID(id)
		if existingAppointment == nil {
			http.Error(w, "Appointment not found", http.StatusNotFound)
			return
		}
		if claims.Role != "staff" && claims.Role != "admin" && existingAppointment.OwnerID != claims.UserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		err = models.DeleteAppointment(id)
		if err != nil {
			utils.Error("Failed to delete appointment: %v", err)
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		utils.Warn("Unsupported method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
