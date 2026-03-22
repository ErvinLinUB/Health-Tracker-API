package main

import "net/http"

func serve(app *application) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/users", app.listUsers)
	mux.HandleFunc("GET /v1/users/{id}", app.getUser)
	mux.HandleFunc("POST /v1/users", app.createUser)
	mux.HandleFunc("PUT /v1/users/{id}", app.updateUser)
	mux.HandleFunc("DELETE /v1/users/{id}", app.deleteUser)

	mux.HandleFunc("GET /v1/workouts", app.listWorkouts)
	mux.HandleFunc("GET /v1/workouts/{id}", app.getWorkout)
	mux.HandleFunc("POST /v1/workouts", app.createWorkout)
	mux.HandleFunc("PUT /v1/workouts/{id}", app.updateWorkout)
	mux.HandleFunc("DELETE /v1/workouts/{id}", app.deleteWorkout)

	mux.HandleFunc("GET /v1/meals", app.listMeals)
	mux.HandleFunc("GET /v1/meals/{id}", app.getMeal)
	mux.HandleFunc("POST /v1/meals", app.createMeal)
	mux.HandleFunc("PUT /v1/meals/{id}", app.updateMeal)
	mux.HandleFunc("DELETE /v1/meals/{id}", app.deleteMeal)

	mux.HandleFunc("GET /v1/health", app.health)

	return http.ListenAndServe(":4000", mux)
}
