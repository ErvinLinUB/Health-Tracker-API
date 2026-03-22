package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (app *application) listMeals(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, user_id, food_name, calories FROM meals ORDER BY id`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var meals []Meal
	for rows.Next() {
		var m Meal
		err := rows.Scan(&m.ID, &m.UserID, &m.FoodName, &m.Calories)
		if err != nil {
			app.serverError(w, err)
			return
		}
		meals = append(meals, m)
	}
	if err = rows.Err(); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"meals": meals}, nil)
}

func (app *application) getMeal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	query := `SELECT id, user_id, food_name, calories FROM meals WHERE id = $1`

	var m Meal
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.UserID, &m.FoodName, &m.Calories)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"meal": m}, nil)
}

func (app *application) createMeal(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   int    `json:"user_id"`
		FoodName string `json:"food_name"`
		Calories int    `json:"calories"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := newValidator()
	v.Check(input.UserID > 0, "user_id", "must be provided")
	v.Check(input.FoodName != "", "food_name", "must be provided")
	v.Check(input.Calories > 0, "calories", "must be greater than 0")

	if !v.Valid() {
		app.failedValidation(w, v.Errors)
		return
	}

	query := `INSERT INTO meals (user_id, food_name, calories) VALUES ($1, $2, $3) RETURNING id`

	var newID int
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.db.QueryRowContext(ctx, query, input.UserID, input.FoodName, input.Calories).Scan(&newID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	newMeal := Meal{ID: newID, UserID: input.UserID, FoodName: input.FoodName, Calories: input.Calories}
	extra := http.Header{"Location": []string{"/v1/meals/" + strconv.Itoa(newID)}}
	app.writeJSON(w, http.StatusCreated, envelope{"meal": newMeal}, extra)
}

func (app *application) updateMeal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	var input struct {
		UserID   int    `json:"user_id"`
		FoodName string `json:"food_name"`
		Calories int    `json:"calories"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := newValidator()
	v.Check(input.UserID > 0, "user_id", "must be provided")
	v.Check(input.FoodName != "", "food_name", "must be provided")
	v.Check(input.Calories > 0, "calories", "must be greater than 0")

	if !v.Valid() {
		app.failedValidation(w, v.Errors)
		return
	}

	query := `UPDATE meals SET user_id = $1, food_name = $2, calories = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	result, err := app.db.ExecContext(ctx, query, input.UserID, input.FoodName, input.Calories, id)
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

	updated := Meal{ID: id, UserID: input.UserID, FoodName: input.FoodName, Calories: input.Calories}
	app.writeJSON(w, http.StatusOK, envelope{"meal": updated}, nil)
}

func (app *application) deleteMeal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	query := `DELETE FROM meals WHERE id = $1`

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
