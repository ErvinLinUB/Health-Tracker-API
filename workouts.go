package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (app *application) listWorkouts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, user_id, type, duration_minutes, calories_burned FROM workouts ORDER BY id`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var workouts []Workout
	for rows.Next() {
		var wo Workout
		err := rows.Scan(&wo.ID, &wo.UserID, &wo.Type, &wo.DurationMinutes, &wo.CaloriesBurned)
		if err != nil {
			app.serverError(w, err)
			return
		}
		workouts = append(workouts, wo)
	}
	if err = rows.Err(); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"workouts": workouts}, nil)
}

func (app *application) getWorkout(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	query := `SELECT id, user_id, type, duration_minutes, calories_burned FROM workouts WHERE id = $1`

	var wo Workout
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.db.QueryRowContext(ctx, query, id).Scan(&wo.ID, &wo.UserID, &wo.Type, &wo.DurationMinutes, &wo.CaloriesBurned)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"workout": wo}, nil)
}

func (app *application) createWorkout(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID          int    `json:"user_id"`
		Type            string `json:"type"`
		DurationMinutes int    `json:"duration_minutes"`
		CaloriesBurned  int    `json:"calories_burned"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := newValidator()
	v.Check(input.UserID > 0, "user_id", "must be provided")
	v.Check(input.Type != "", "type", "must be provided")
	v.Check(input.DurationMinutes > 0, "duration_minutes", "must be greater than 0")
	v.Check(input.CaloriesBurned > 0, "calories_burned", "must be greater than 0")

	if !v.Valid() {
		app.failedValidation(w, v.Errors)
		return
	}

	query := `INSERT INTO workouts (user_id, type, duration_minutes, calories_burned) VALUES ($1, $2, $3, $4) RETURNING id`

	var newID int
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.db.QueryRowContext(ctx, query, input.UserID, input.Type, input.DurationMinutes, input.CaloriesBurned).Scan(&newID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	newWorkout := Workout{ID: newID, UserID: input.UserID, Type: input.Type, DurationMinutes: input.DurationMinutes, CaloriesBurned: input.CaloriesBurned}
	extra := http.Header{"Location": []string{"/v1/workouts/" + strconv.Itoa(newID)}}
	app.writeJSON(w, http.StatusCreated, envelope{"workout": newWorkout}, extra)
}

func (app *application) updateWorkout(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	var input struct {
		UserID          int    `json:"user_id"`
		Type            string `json:"type"`
		DurationMinutes int    `json:"duration_minutes"`
		CaloriesBurned  int    `json:"calories_burned"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := newValidator()
	v.Check(input.UserID > 0, "user_id", "must be provided")
	v.Check(input.Type != "", "type", "must be provided")
	v.Check(input.DurationMinutes > 0, "duration_minutes", "must be greater than 0")
	v.Check(input.CaloriesBurned > 0, "calories_burned", "must be greater than 0")

	if !v.Valid() {
		app.failedValidation(w, v.Errors)
		return
	}

	query := `UPDATE workouts SET user_id = $1, type = $2, duration_minutes = $3, calories_burned = $4 WHERE id = $5`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	result, err := app.db.ExecContext(ctx, query, input.UserID, input.Type, input.DurationMinutes, input.CaloriesBurned, id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		app.serverError(w, err)
		return
	}
	if rowsAffected == 0 {
		app.notFound(w)
		return
	}

	updated := Workout{ID: id, UserID: input.UserID, Type: input.Type, DurationMinutes: input.DurationMinutes, CaloriesBurned: input.CaloriesBurned}
	app.writeJSON(w, http.StatusOK, envelope{"workout": updated}, nil)
}

func (app *application) deleteWorkout(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	query := `DELETE FROM workouts WHERE id = $1`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	result, err := app.db.ExecContext(ctx, query, id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		app.serverError(w, err)
		return
	}
	if rowsAffected == 0 {
		app.notFound(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
