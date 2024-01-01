package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/aleale2121/interactive-presentation/util"
	_ "github.com/lib/pq"
)

const (
	presenatationData2 string = `[
		{
		"question": "What is your favorite pet?",
		"options": [
				{"key": "A", "value": "Dog"},
				{"key": "B", "value": "Cat"},
				{"key": "C", "value": "Crocodile"}
			]},
		{
		"question": "Which of the countries would you like to visit the most?",
		"options": [
			{"key": "A", "value": "Argentina"},
			{"key": "B", "value": "Austria"},
			{"key": "C", "value": "Australia"}
		]}
	]`

	presenatationData8 string = `[
		{
		  "question": "What is your favorite pet?",
		  "options": [
			{"key": "A", "value": "Dog"},
			{"key": "B", "value": "Cat"},
			{"key": "C", "value": "Crocodile"}
		  ]
		},
		{
		  "question": "Which of the countries would you like to visit the most?",
		  "options": [
			{"key": "A", "value": "Argentina"},
			{"key": "B", "value": "Austria"},
			{"key": "C", "value": "Australia"}
		  ]
		},
		{
		  "question": "What is your favorite color?",
		  "options": [
			{"key": "A", "value": "Red"},
			{"key": "B", "value": "Blue"},
			{"key": "C", "value": "Green"},
			{"key": "D", "value": "Yellow"},
			{"key": "E", "value": "Purple"}
		  ]
		},
		{
		  "question": "How do you prefer to travel?",
		  "options": [
			{"key": "A", "value": "Car"},
			{"key": "B", "value": "Train"},
			{"key": "C", "value": "Plane"},
			{"key": "D", "value": "Bicycle"},
			{"key": "E", "value": "Walking"}
		  ]
		},
		{
		  "question": "What is your favorite type of cuisine?",
		  "options": [
			{"key": "A", "value": "Italian"},
			{"key": "B", "value": "Japanese"},
			{"key": "C", "value": "Mexican"},
			{"key": "D", "value": "Indian"},
			{"key": "E", "value": "Mediterranean"}
		  ]
		},
		{
		  "question": "Which programming language do you prefer?",
		  "options": [
			{"key": "A", "value": "JavaScript"},
			{"key": "B", "value": "Python"},
			{"key": "C", "value": "Java"},
			{"key": "D", "value": "Go"},
			{"key": "E", "value": "Ruby"}
		  ]
		},
		{
		  "question": "What is your favorite movie genre?",
		  "options": [
			{"key": "A", "value": "Action"},
			{"key": "B", "value": "Comedy"},
			{"key": "C", "value": "Drama"},
			{"key": "D", "value": "Science Fiction"},
			{"key": "E", "value": "Thriller"}
		  ]
		},
		{
		  "question": "Do you prefer reading fiction or non-fiction?",
		  "options": [
			{"key": "A", "value": "Fiction"},
			{"key": "B", "value": "Non-Fiction"}
		  ]
		}
	  ]	  
	  `
)

var testQueries *Queries
var testDB *sql.DB

var presentationData map[int]string

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	if err != nil {
		log.Fatal("cannot marshal:", err)
		return
	}
	presentationData = map[int]string{
		2: presenatationData2,
		8: presenatationData8,
	}
	os.Exit(m.Run())
}
