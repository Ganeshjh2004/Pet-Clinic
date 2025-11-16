package handlers

import (
	"encoding/json"
	"net/http"
	"petclinic/models"
	"petclinic/utils"
	"strconv"
)

func PetsHandler(w http.ResponseWriter, r *http.Request) {
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
			pets := models.GetPetsByOwnerID(claims.UserID)
			if pets == nil {
				utils.Warn("No pets found for owner")
				pets = []models.Pet{}
			}
			json.NewEncoder(w).Encode(pets)
		} else {
			pets := models.GetAllPets()
			if pets == nil {
				utils.Warn("No pets found in database")
				pets = []models.Pet{}
			}
			json.NewEncoder(w).Encode(pets)
		}

	case http.MethodPost:
		var pet models.Pet
		err := json.NewDecoder(r.Body).Decode(&pet)
		if err != nil {
			utils.Error("Failed to decode POST body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Owner role: always set owner_id to claims.UserID
		if claims.Role == "owner" {
			pet.OwnerID = claims.UserID
		} else if claims.Role != "staff" && claims.Role != "admin" {
			// If not owner, staff, or admin, forbid
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// For staff and admin, owner_id from body is used

		err = models.AddPet(pet)
		if err != nil {
			utils.Error("Error adding pet in DB: %v", err)
			http.Error(w, "Failed to add pet", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Warn("Invalid pet ID provided for update")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		existingPet, err := models.GetPetByID(id)
		if err != nil || existingPet == nil {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		// Only staff, admin, or owner of pet can update
		if claims.Role != "staff" && claims.Role != "admin" && existingPet.OwnerID != claims.UserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var pet models.Pet
		err = json.NewDecoder(r.Body).Decode(&pet)
		if err != nil {
			utils.Error("Failed to decode PUT body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = models.UpdatePet(id, pet)
		if err != nil {
			utils.Error("Error updating pet in DB: %v", err)
			http.Error(w, "Failed to update pet", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Warn("Invalid pet ID for deletion")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		existingPet, err := models.GetPetByID(id)
		if err != nil || existingPet == nil {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		// Only staff, admin, or owner of pet can delete
		if claims.Role != "staff" && claims.Role != "admin" && existingPet.OwnerID != claims.UserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		err = models.DeletePet(id)
		if err != nil {
			utils.Error("Error deleting pet in DB: %v", err)
			http.Error(w, "Failed to delete pet", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		utils.Warn("Unsupported method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
