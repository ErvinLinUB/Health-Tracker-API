package main

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Workout struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	Type            string `json:"type"`
	DurationMinutes int    `json:"duration_minutes"`
	CaloriesBurned  int    `json:"calories_burned"`
}

type Meal struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	FoodName string `json:"food_name"`
	Calories int    `json:"calories"`
}
