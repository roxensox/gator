package config

type Config struct {
	DB_URL      *string `json:"db_url"`
	CurrentUser *string `json:"current_user"`
}
