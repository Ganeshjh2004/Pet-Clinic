package models

import (
	"database/sql"
	"petclinic/db"
	"petclinic/utils"
	"time"
)

type Appointment struct {
	ID      int       `json:"id"`
	Date    time.Time `json:"date"`
	Time    string    `json:"time"`
	PetID   int       `json:"pet_id"`
	Reason  string    `json:"reason"`
	OwnerID int       `json:"owner_id"` // Add owner ID field for ownership tracking
}

func GetAllAppointments() []Appointment {
	rows, err := db.DB.Query("SELECT id, date, time, pet_id, reason FROM appointments")
	if err != nil {
		utils.Error("Failed to fetch appointments: %v", err)
		return nil
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		err := rows.Scan(&a.ID, &a.Date, &a.Time, &a.PetID, &a.Reason)
		if err != nil {
			utils.Warn("Failed to scan appointment row: %v", err)
			continue
		}
		appointments = append(appointments, a)
	}
	if err = rows.Err(); err != nil {
		utils.Error("Rows error in GetAllAppointments: %v", err)
	}
	return appointments
}

func AddAppointment(a Appointment) error {
	_, err := db.DB.Exec("INSERT INTO appointments (date, time, pet_id, reason, owner_id) VALUES ($1, $2, $3, $4, $5)", a.Date, a.Time, a.PetID, a.Reason, a.OwnerID)
	if err != nil {
		utils.Error("AddAppointment DB error: %v", err)
	}
	return err
}

func UpdateAppointment(id int, a Appointment) error {
	_, err := db.DB.Exec("UPDATE appointments SET date=$1, time=$2, pet_id=$3, reason=$4 WHERE id=$5", a.Date, a.Time, a.PetID, a.Reason, id)
	if err != nil {
		utils.Error("UpdateAppointment DB error: %v", err)
	}
	return err
}

func DeleteAppointment(id int) error {
	_, err := db.DB.Exec("DELETE FROM appointments WHERE id=$1", id)
	if err != nil {
		utils.Error("DeleteAppointment DB error: %v", err)
	}
	return err
}

func GetAppointmentByID(id int) *Appointment {
	var a Appointment
	err := db.DB.QueryRow("SELECT id, date, time, pet_id, reason FROM appointments WHERE id=$1", id).Scan(&a.ID, &a.Date, &a.Time, &a.PetID, &a.Reason)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Warn("No appointment found with id: %d", id)
			return nil
		}
		utils.Error("GetAppointmentByID DB error: %v", err)
		return nil
	}
	return &a
}

func GetAppointmentsByOwnerID(ownerID int) []Appointment {
	rows, err := db.DB.Query(
		`SELECT a.id, a.date, a.time, a.pet_id, a.reason
         FROM appointments a
         JOIN pets p ON a.pet_id = p.id
         WHERE p.owner_id = $1`, ownerID)
	if err != nil {
		utils.Error("Failed to fetch appointments for owner %d: %v", ownerID, err)
		return nil
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		err := rows.Scan(&a.ID, &a.Date, &a.Time, &a.PetID, &a.Reason)
		if err != nil {
			utils.Warn("Failed to scan appointment row: %v", err)
			continue
		}
		appointments = append(appointments, a)
	}
	if err = rows.Err(); err != nil {
		utils.Error("Rows error in GetAppointmentsByOwnerID: %v", err)
	}
	return appointments
}
