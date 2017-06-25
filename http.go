package main

import (
    "fmt"
    "encoding/json"
    "io"
    "net/http"
    "strconv"
)

// ================================================= //
// ================================================= //

func launchHttpServer() {

    return

    http.HandleFunc("/starting", postStarting)
    http.HandleFunc("/onServer", postOnServer)
    http.HandleFunc("/onInterne", postOnInterne)
    http.HandleFunc("/ending", postEnding)

    http.HandleFunc("/task", postCreateTask)

    err := http.ListenAndServe(":8080", nil)

    if err != nil {
        fmt.Println("starting listening for payload messages")
    } else {
        fmt.Errorf("an error occured while starting payload server %s", err.Error())
    }

}

// ================================================= //
// ================================================= //

// task array from http request
type TaskCollection struct {
    Name string `json:"name"`
    Tasks []Task `json:"data"`
}

// ================================================= //
// ================================================= //

func postStarting(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    fmt.Println("Request received")

    fmt.Println("___________________________________")
    // Read the body into a string for json decoding
    var content = &JsonB{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    fmt.Println( content )
    fmt.Println("___________________________________")

    w.WriteHeader(http.StatusOK)
}

// ================================================= //
// ================================================= //

func postOnServer(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Read the body into a string for json decoding
    var content = &TaskCollection{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    fmt.Println("Request received")

    // Go through each payload and queue items individually to be posted to S3
    for _, task := range content.Tasks {

        // Push the work onto the queue.
        TaskQueue <- task

        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
    }

    w.WriteHeader(http.StatusOK)
}

// ================================================= //
// ================================================= //

func postOnInterne(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Read the body into a string for json decoding
    var content = &TaskCollection{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    fmt.Println("Request received")

    // Go through each payload and queue items individually to be posted to S3
    for _, task := range content.Tasks {

        // Push the work onto the queue.
        TaskQueue <- task

        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
    }

    w.WriteHeader(http.StatusOK)
}

// ================================================= //
// ================================================= //

func postEnding(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Read the body into a string for json decoding
    var content = &TaskCollection{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    fmt.Println("Request received")

    // Go through each payload and queue items individually to be posted to S3
    for _, task := range content.Tasks {

        // Push the work onto the queue.
        TaskQueue <- task

        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
    }

    w.WriteHeader(http.StatusOK)
}

// ================================================= //
// ================================================= //

// handle the http post request on post localhost:8080/task
func postCreateTask(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Read the body into a string for json decoding
    var content = &TaskCollection{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    fmt.Println("Request received")

    // Go through each payload and queue items individually to be posted to S3
    for _, task := range content.Tasks {

        // Push the work onto the queue.
        TaskQueue <- task

        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
    }

    w.WriteHeader(http.StatusOK)
}
