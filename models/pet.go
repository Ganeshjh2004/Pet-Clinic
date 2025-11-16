package models

import (
	"database/sql"
	"petclinic/db"
	"petclinic/utils"
)

type Pet struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Breed   string `json:"breed"`
	OwnerID int    `json:"owner_id"`
	History string `json:"history"`
}

func GetAllPets() []Pet {
	rows, err := db.DB.Query("SELECT id, name, species, breed, owner_id, history FROM pets")
	if err != nil {
		utils.Error("Failed to fetch pets: %v", err)
		return nil
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var p Pet
		err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.OwnerID, &p.History)
		if err != nil {
			utils.Warn("Failed to scan pet row: %v", err)
			continue
		}
		pets = append(pets, p)
	}
	if err = rows.Err(); err != nil {
		utils.Error("Rows error in GetAllPets: %v", err)
	}
	return pets
}

func GetPetsByOwnerID(ownerID int) []Pet {
	rows, err := db.DB.Query("SELECT id, name, species, breed, owner_id, history FROM pets WHERE owner_id=$1", ownerID)
	if err != nil {
		utils.Error("Failed to fetch pets for owner %d: %v", ownerID, err)
		return nil
	}
	defer rows.Close()

	var pets []Pet
	for rows.Next() {
		var p Pet
		err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.OwnerID, &p.History)
		if err != nil {
			utils.Warn("Failed to scan pet row: %v", err)
			continue
		}
		pets = append(pets, p)
	}
	if err = rows.Err(); err != nil {
		utils.Error("Rows error in GetPetsByOwnerID: %v", err)
	}
	return pets
}

func AddPet(p Pet) error {
	_, err := db.DB.Exec("INSERT INTO pets (name, species, breed, owner_id, history) VALUES ($1, $2, $3, $4, $5)", p.Name, p.Species, p.Breed, p.OwnerID, p.History)
	if err != nil {
		utils.Error("AddPet DB error: %v", err)
	}
	return err
}

func UpdatePet(id int, p Pet) error {
	_, err := db.DB.Exec("UPDATE pets SET name=$1, species=$2, breed=$3, owner_id=$4, history=$5 WHERE id=$6", p.Name, p.Species, p.Breed, p.OwnerID, p.History, id)
	if err != nil {
		utils.Error("UpdatePet DB error: %v", err)
	}
	return err
}

func DeletePet(id int) error {
	_, err := db.DB.Exec("DELETE FROM pets WHERE id=$1", id)
	if err != nil {
		utils.Error("DeletePet DB error: %v", err)
	}
	return err
}

func GetPetByID(id int) (*Pet, error) {
	var p Pet
	err := db.DB.QueryRow("SELECT id, name, species, breed, owner_id, history FROM pets WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.OwnerID, &p.History)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Warn("No pet found with id: %d", id)
			return nil, nil
		}
		utils.Error("GetPetByID DB error: %v", err)
		return nil, err
	}
	return &p, nil
}
