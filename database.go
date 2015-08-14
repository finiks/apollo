package main

import (
    "log"
    "os"
    "io/ioutil"
    "encoding/json"
)

type Movie struct {
    Title string
    State string
    Year string
    ImdbID string
    Rating int
}

type Database struct {
    Movies []Movie
}

func createDatabase() *Database {
    d := &Database{}

    d.load()
    d.save()

    return d
}

func (d *Database) load() {
    path := os.Getenv("HOME") + "/.config/apollo/database.json"
    cont, err := ioutil.ReadFile(path)
    if err != nil {
        log.Print(err)
        return
    }

    err = json.Unmarshal(cont, d)
    if err != nil {
        log.Fatal(err)
    }
}

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
