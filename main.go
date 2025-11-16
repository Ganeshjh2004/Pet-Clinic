package main

import (
	"log"
	"net/http"
	"petclinic/db"
	"petclinic/handlers"
	"petclinic/middleware"
)

func main() {
	// Initialize database connection
	db.InitDB()

	// Register API endpoints
	http.HandleFunc("/login", handlers.LoginHandler)

	// Pets: Only staff/admin can access
	http.Handle("/pets",
		middleware.Logging(
			middleware.AuthMiddleware(
				http.HandlerFunc(handlers.PetsHandler),
			),
		),
	)

	// Owners: staff, admin, and owner can access
	http.Handle("/owners",
		middleware.Logging(
			middleware.AuthMiddleware(
				middleware.RoleBasedAccess("staff", "admin", "owner")(
					http.HandlerFunc(handlers.OwnersHandler),
				),
			),
		),
	)

	// Appointments: staff, admin, and owner can access
	http.Handle("/appointments",
		middleware.Logging(
			middleware.AuthMiddleware(
				middleware.RoleBasedAccess("staff", "admin", "owner")(
					http.HandlerFunc(handlers.AppointmentsHandler),
				),
			),
		),
	)

	http.HandleFunc("/upload", handlers.UploadFileHandler)
	http.HandleFunc("/download", handlers.DownloadFileHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
