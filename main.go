package main

import (
	"log"
	"net/http"
	"time"

	"course-registration/internal/httpapi"
	"course-registration/internal/service"
	"course-registration/internal/store"
)

func main() {
	st := store.NewStore()

	// Background audit worker (goroutine) + channel-based logic
	audit := service.NewAuditWorker(200)
	audit.Start()
	defer audit.Stop()

	enrollSvc := service.NewEnrollmentService(st, audit)

	router := httpapi.NewRouter(st, enrollSvc, audit)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Milestone 2 server running on http://localhost:8080")
	log.Println("Health: GET http://localhost:8080/health")
	log.Fatal(srv.ListenAndServe())
}
