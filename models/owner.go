package models

import (
	"database/sql"
	"petclinic/db"
	"petclinic/utils"
)

type Owner struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
}

func GetAllOwners() []Owner {
	rows, err := db.DB.Query("SELECT id, name, contact, email FROM owners")
	if err != nil {
		utils.Error("Failed to fetch owners: %v", err)
		return nil
	}
	defer rows.Close()
	var owners []Owner
	for rows.Next() {
		var o Owner
		err := rows.Scan(&o.ID, &o.Name, &o.Contact, &o.Email)
		if err != nil {
			utils.Warn("Failed to scan owner row: %v", err)
			continue
		}
		owners = append(owners, o)
	}
	if err = rows.Err(); err != nil {
		utils.Error("Rows error in GetAllOwners: %v", err)
	}
	return owners
}

func AddOwner(o Owner) error {
	_, err := db.DB.Exec("INSERT INTO owners (name, contact, email) VALUES ($1, $2, $3)", o.Name, o.Contact, o.Email)
	if err != nil {
		utils.Error("AddOwner DB error: %v", err)
	}
	return err
}

func UpdateOwner(id int, o Owner) error {
	_, err := db.DB.Exec("UPDATE owners SET name=$1, contact=$2, email=$3 WHERE id=$4", o.Name, o.Contact, o.Email, id)
	if err != nil {
		utils.Error("UpdateOwner DB error: %v", err)
	}
	return err
}

func DeleteOwner(id int) error {
	_, err := db.DB.Exec("DELETE FROM owners WHERE id=$1", id)
	if err != nil {
		utils.Error("DeleteOwner DB error: %v", err)
	}
	return err
}

func GetOwnerByID(id int) (*Owner, error) {
	var o Owner
	err := db.DB.QueryRow("SELECT id, name, contact, email FROM owners WHERE id=$1", id).
		Scan(&o.ID, &o.Name, &o.Contact, &o.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Warn("No owner found with id: %d", id)
			return nil, nil
		}
		utils.Error("GetOwnerByID DB error: %v", err)
		return nil, err
	}
	return &o, nil
}
