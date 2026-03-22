# Health Tracker API

**Test 1 - Systems Programming and Computer Organization**

**Checkpoint 1**

- Database

    - Created healthtracker PostgreSQL database
    - Created users, workouts, and meals tables via SQL migrations
Added sample data (3 rows per table) via a 4th migration

- API Endpoints

    - Full CRUD for users (list, get, create, update, delete)
    - Full CRUD for workouts (list, get, create, update, delete)
    - Full CRUD for meals (list, get, create, update, delete)
    - Health check endpoint

- Code Structure

    - Refactored into multiple files (main, routes, helpers, models, - validator, users, workouts, meals)
    - JSON input validation using the Validator
    - Proper error responses (400, 404, 422, 500)
    - Data saved to and retrieved from the database in JSON format

- Tested with curl

    - Verified all endpoints work correctly
    - Confirmed data is being written to and read from the database

Completion date: March 22nd, 2026

*- Ervin Lin*