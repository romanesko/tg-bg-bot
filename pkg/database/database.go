package database

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/config"
	"database/sql"
	"fmt"
	"github.com/go-telegram/bot/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

type Db struct {
	DB *sql.DB
}

func NewDb() *Db {
	var cfg = config.GetConfig()
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		panic(err)

	}

	api := api.NewApi(cfg.Api.Host, cfg.Api.Key)
	createTables(db, api)
	return &Db{DB: db}
}

func createTables(db *sql.DB, api *api.Api) {
	log.Println("Creating tables")

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tg_chat_id INTEGER NOT NULL,
		tg_username TEXT NOT NULL,
		first_name TEXT,
		birthday date,
		birthtime time,
		time_unknown boolean not null DEFAULT false,
		country_id INTEGER,
		country_name INTEGER,
		city_id INTEGER,
		city_name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cities (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		display_name TEXT NOT NULL,
		country_name text NOT NULL,
		country_id integer NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	exists := false
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM cities)").Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		cities, err := api.GetCities()
		if err != nil {
			log.Fatal(err)
		}
		for _, city := range cities {
			_, err := db.Exec("INSERT INTO cities (name, display_name, country_name, country_id) VALUES (?, ?, ?, ?)", strings.ToLower(city.CityName), city.DisplayName, city.CountryName, city.CountryID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

type User struct {
	TgChatID    int64
	TgUsername  string
	FirstName   string
	BirthDate   *time.Time
	BirthTime   *time.Time
	CountryID   *int
	CountryName *string
	CityID      *int
	CityName    *string
	TimeUnknown bool
}

func (db *Db) GetUser(chat models.Chat) (User, error) {

	exists := false
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE tg_chat_id = ?)", chat.ID).Scan(&exists)
	if err != nil {
		return User{}, fmt.Errorf("GetUser: %v", err)
	}

	user := User{}

	if !exists {
		user.TgChatID = chat.ID
		user.TgUsername = chat.Username
		user.FirstName = chat.FirstName
		_, err := db.DB.Exec("INSERT INTO users (tg_chat_id, tg_username, first_name) VALUES (?, ?, ?)", chat.ID, chat.Username, chat.FirstName)
		if err != nil {
			return user, fmt.Errorf("CreateUser: %v", err)
		}
		return user, nil
	}

	err = db.DB.QueryRow("SELECT tg_chat_id, tg_username, first_name, birthday, birthtime, country_id, country_name, city_id, city_name,time_unknown FROM users WHERE tg_chat_id = ?", chat.ID).Scan(&user.TgChatID, &user.TgUsername, &user.FirstName, &user.BirthDate, &user.BirthTime, &user.CountryID, &user.CountryName, &user.CityID, &user.CityName, &user.TimeUnknown)
	if err != nil {
		return user, fmt.Errorf("GetUser: %v", err)
	}

	return user, nil
}

func (db *Db) UpdateUser(user *User) error {
	_, err := db.DB.Exec("UPDATE users SET first_name = ?, birthday = ?, birthtime = ?,country_id = ?, country_name = ?, city_id = ?, city_name = ?, time_unknown = ? WHERE tg_chat_id = ?",
		user.FirstName, user.BirthDate, user.BirthTime, user.CountryID, user.CountryName, user.CityID, user.CityName, user.TimeUnknown, user.TgChatID)
	if err != nil {
		return fmt.Errorf("UpdateUser: %v", err)
	}
	return nil

}

type City struct {
	ID          int
	Name        string
	DisplayName string
	CountryName string
	CountryID   int
}

func (db *Db) GetCity(cityId int) (City, error) {
	city := City{}
	err := db.DB.QueryRow("SELECT * FROM cities WHERE id = ?", cityId).Scan(&city.ID, &city.Name, &city.DisplayName, &city.CountryName, &city.CountryID)
	if err != nil {
		return city, fmt.Errorf("GetCity: %v", err)
	}
	return city, nil

}

func (db *Db) FindCity(name string) ([]City, error) {
	cities := []City{}

	rows, err := db.DB.Query("SELECT * FROM cities where LOWER(name) like ?", fmt.Sprintf("%s%s", name, "%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var city City
		err = rows.Scan(&city.ID, &city.Name, &city.DisplayName, &city.CountryName, &city.CountryID)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}
