package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Entry is a single media title.
type Entry struct {
	Title        string
	State        string
	Year         string
	TagID        string
	Rating       int
	Info1        string
	EpisodeDone  int
	EpisodeTotal int
}

// Database is the entirety of all media titles.
type Database struct {
	Movies []Entry
	Series []Entry
	Anime  []Entry
	Games  []Entry
	Books  []Entry
}

// NewDatabase creates a new Database and returns it.
func newDatabase() *Database {
	d := &Database{}

	d.load()
	d.save()

	return d
}

// Load fetches the user's database from a file.
func (d *Database) load() {
	path := os.Getenv("HOME") + "/.config/apollo/database.json"
	cont, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(cont, d)
	if err != nil {
		log.Fatal(err)
	}
}

// Save takes the database and saves it to a file.
func (d *Database) save() {
	cont, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}

	path := os.Getenv("HOME") + "/.config/apollo/database.json"
	err = ioutil.WriteFile(path, cont, 0644)
	if err != nil {
		log.Print(err)
	}
}
