package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, email FROM users ORDER BY id`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			app.serverError(w, err)
			return
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	query := `SELECT id, name, email FROM users WHERE id = $1`

	var u User
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"user": u}, nil)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := newValidator()
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(input.Email != "", "email", "must be provided")

	if !v.Valid() {
		app.failedValidation(w, v.Errors)
		return
	}

	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`

	var newID int
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.db.QueryRowContext(ctx, query, input.Name, input.Email).Scan(&newID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	newUser := User{ID: newID, Name: input.Name, Email: input.Email}
	extra := http.Header{"Location": []string{"/v1/users/" + strconv.Itoa(newID)}}
	app.writeJSON(w, http.StatusCreated, envelope{"user": newUser}, extra)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := newValidator()
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(input.Email != "", "email", "must be provided")

	if !v.Valid() {
		app.failedValidation(w, v.Errors)
		return
	}

	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	result, err := app.db.ExecContext(ctx, query, input.Name, input.Email, id)
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

	updated := User{ID: id, Name: input.Name, Email: input.Email}
	app.writeJSON(w, http.StatusOK, envelope{"user": updated}, nil)
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	query := `DELETE FROM users WHERE id = $1`

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
