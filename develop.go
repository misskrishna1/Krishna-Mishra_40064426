package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"grom.io/gorm"
	_ "github.com/gorm.io/driver/mysql"
)

type Player struct {
	gorm.Model
	PlayerID int    `json:"playerId" gorm:"primary_key"`
	Name     string `json:"name"`
	Team     string `json:"team"`
}

type PlayerScores struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Match    string `json:"match"`
	Runs     int    `json:"runs"`
	Wickets  int    `json:"wickets"`
	PlayerID int    `json:"playerId" gorm:"foreign_key:"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:Krishna123*(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.Exec("CREATE DATABASE players_db")
	db.Exec("USE players_db")

	db.AutoMigrate(&Player{}, &PlayerScores{})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/player", createPlayer).Methods("POST")
	router.HandleFunc("/player/{playerId}/score", createPlayerScore).Methods("POST")
	router.HandleFunc("/players/{playerId}", getPlayer).Methods("GET")
	router.HandleFunc("/players", getPlayers).Methods("GET")
	router.HandleFunc("/players/scores", getPlayerScores).Methods("GET")
	
	initDB()

	log.Fatal(http.ListenAndServe(":8051", router))
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	json.NewDecoder(r.Body).Decode(&player)
	
	db.Create(&player)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func createPlayerScore(w http.ResponseWriter, r *http.Request) {
	var playerscore PlayerScores
	json.NewDecoder(r.Body).Decode(&playerscore)
	db.Create(&playerscore)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playerscore)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var players []Player
	db.Find(&players)
	json.NewEncoder(w).Encode(players)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputPlayerID := params["playerId"]

	var player Player
	db.First(&player, inputPlayerID)
	json.NewEncoder(w).Encode(player)
}

func getPlayerScores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var playerscores []PlayerScores
	db.Find(&playerscores)
	json.NewEncoder(w).Encode(playerscores)
}
