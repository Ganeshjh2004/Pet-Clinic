package handlers

import (
	"encoding/json"
	"net/http"
	"petclinic/models"
	"petclinic/utils"
	"strconv"
)

func OwnersHandler(w http.ResponseWriter, r *http.Request) {
	utils.Info("Received %s request at %s", r.Method, r.URL.Path)

	claimsVal := r.Context().Value("userClaims")
	claims, ok := claimsVal.(*utils.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Staff/admin see all; owner sees only self
		if claims.Role == "owner" {
			owner, err := models.GetOwnerByID(claims.UserID)
			if err != nil || owner == nil {
				http.Error(w, "Owner not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode([]*models.Owner{owner})
		} else {
			owners := models.GetAllOwners()
			if owners == nil {
				utils.Warn("No owners found in database")
			}
			json.NewEncoder(w).Encode(owners)
		}

	case http.MethodPost:
		if claims.Role != "staff" && claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var owner models.Owner
		err := json.NewDecoder(r.Body).Decode(&owner)
		if err != nil {
			utils.Error("Failed to decode POST body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		err = models.AddOwner(owner)
		if err != nil {
			utils.Error("Error adding owner in DB: %v", err)
			http.Error(w, "Failed to add owner", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Warn("Invalid owner ID provided for update")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		// Only staff/admin or owner themselves can update
		if claims.Role != "staff" && claims.Role != "admin" && claims.UserID != id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var owner models.Owner
		err = json.NewDecoder(r.Body).Decode(&owner)
		if err != nil {
			utils.Error("Failed to decode PUT body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		err = models.UpdateOwner(id, owner)
		if err != nil {
			utils.Error("Error updating owner in DB: %v", err)
			http.Error(w, "Failed to update owner", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Warn("Invalid owner ID for deletion")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		// Only staff/admin can delete
		if claims.Role != "staff" && claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		err = models.DeleteOwner(id)
		if err != nil {
			utils.Error("Error deleting owner in DB: %v", err)
			http.Error(w, "Failed to delete owner", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		utils.Warn("Unsupported method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
