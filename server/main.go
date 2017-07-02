package main

import (
    "io"
    "fmt"
    "time"
    "net/http"
    "encoding/json"
)

func main() {

    fmt.Println("starting server http")

    http.HandleFunc("/create/starting", postStarting)
    http.HandleFunc("/create/onServer", postOnServer)
    http.HandleFunc("/create/onInterne", postOnInterne)
    http.HandleFunc("/create/ending", postEnding)



//    http.HandleFunc("/create/starting", postCreateStarting)
//    http.HandleFunc("/create/onServer", postCreateOnServer)
//    http.HandleFunc("/create/onInterne", postCreateOnInterne)
//    http.HandleFunc("/create/ending", postCreateEnding)
//
//    http.HandleFunc("/update/starting", postUpdateStarting)
//    http.HandleFunc("/update/onServer", postUpdateOnServer)
//    http.HandleFunc("/update/onInterne", postUpdateOnInterne)
//    http.HandleFunc("/update/ending", postUpdateEnding)
//
//    http.HandleFunc("/delete/starting", postDeleteStarting)
//    http.HandleFunc("/delete/onServer", postDeleteOnServer)
//    http.HandleFunc("/delete/onInterne", postDeleteOnInterne)
//    http.HandleFunc("/delete/ending", postDeleteEnding)

    err := http.ListenAndServe(":8080", nil)

    if err != nil {
        fmt.Println("starting listening for payload messages")
    } else {
        fmt.Errorf("an error occured while starting payload server %s", err.Error())
    }

    time.Sleep(time.Hour)
}

// ================================================= //
// ================================================= //

var MaxLength int64 = 2048

// Data received from the executor
type StepData struct {
    Name        string
    Arguments   JsonB
    Buffer      JsonB
}

// Starting step
func postStarting(w http.ResponseWriter, r *http.Request) {

    fmt.Println("========================================")

    fmt.Println("----------- Request received -----------")

    w.Header().Set("Content-Type", "application/json")

    if r.Method != "POST" {

        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    // Read data from body request
    var body = &StepData{}

    // Decode body json data
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&body)

    if err != nil {

        fmt.Errorf("an error occured while deserializing message")

        w.WriteHeader(http.StatusBadRequest)

        fmt.Println("========================================")

        return
    }

    fmt.Println( body )

    fmt.Println("------------- Body decoded -------------")

    // ========================================
    // ========================================
    // ============= DO THE WORK ==============
    // ========================================
    // ========================================

    var buffer = body.Buffer

    buffer["steps"] = []string{"starting"}

    // ========================================
    // ========================================
    // ============ / DO THE WORK =============
    // ========================================
    // ========================================

    fmt.Println("------------- Action done --------------")

    js, err := json.Marshal( buffer )

    if err != nil {

        http.Error(w, err.Error(), http.StatusInternalServerError)

        fmt.Println("========================================")

        return
    }

    w.Header().Set("Content-Type", "application/json")

    w.Write(js)

    fmt.Println("------------- Send response ------------")

    fmt.Println("========================================")
}

func postOnServer(w http.ResponseWriter, r *http.Request) {

    fmt.Println("========================================")

    fmt.Println("----------- Request received -----------")

    w.Header().Set("Content-Type", "application/json")

    if r.Method != "POST" {

        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    // Read data from body request
    var body = &StepData{}

    // Decode body json data
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&body)

    if err != nil {

        fmt.Errorf("an error occured while deserializing message")

        w.WriteHeader(http.StatusBadRequest)

        fmt.Println("========================================")

        return
    }

    fmt.Println( body )

    fmt.Println("------------- Body decoded -------------")

    // ========================================
    // ========================================
    // ============= DO THE WORK ==============
    // ========================================
    // ========================================

    var buffer = body.Buffer

    buffer["steps"] = []string{ "starting", "onServer" }

    // ========================================
    // ========================================
    // ============ / DO THE WORK =============
    // ========================================
    // ========================================

    fmt.Println("------------- Action done --------------")

    js, err := json.Marshal( buffer )

    if err != nil {

        http.Error(w, err.Error(), http.StatusInternalServerError)

        fmt.Println("========================================")

        return
    }

    w.Header().Set("Content-Type", "application/json")

    w.Write(js)

    fmt.Println("------------- Send response ------------")

    fmt.Println("========================================")
}



func postOnInterne(w http.ResponseWriter, r *http.Request) {

    fmt.Println("========================================")

    fmt.Println("----------- Request received -----------")

    w.Header().Set("Content-Type", "application/json")

    if r.Method != "POST" {

        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    // Read data from body request
    var body = &StepData{}

    // Decode body json data
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&body)

    if err != nil {

        fmt.Errorf("an error occured while deserializing message")

        w.WriteHeader(http.StatusBadRequest)

        fmt.Println("========================================")

        return
    }

    fmt.Println( body )

    fmt.Println("------------- Body decoded -------------")

    // ========================================
    // ========================================
    // ============= DO THE WORK ==============
    // ========================================
    // ========================================

    var buffer = body.Buffer

    buffer["data"] = map[string]interface{}{"name":"noemi","informations":[]string{"toto", "success"}}

    // ========================================
    // ========================================
    // ============ / DO THE WORK =============
    // ========================================
    // ========================================

    fmt.Println("------------- Action done --------------")

    js, err := json.Marshal( buffer )

    if err != nil {

        http.Error(w, err.Error(), http.StatusInternalServerError)

        fmt.Println("========================================")

        return
    }

    w.Header().Set("Content-Type", "application/json")

    w.Write(js)

    fmt.Println("------------- Send response ------------")

    fmt.Println("========================================")
}




func postEnding(w http.ResponseWriter, r *http.Request) {

    fmt.Println("========================================")

    fmt.Println("----------- Request received -----------")

    w.Header().Set("Content-Type", "application/json")

    if r.Method != "POST" {

        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    // Read data from body request
    var body = &StepData{}

    // Decode body json data
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&body)

    if err != nil {

        fmt.Errorf("an error occured while deserializing message")

        w.WriteHeader(http.StatusBadRequest)

        fmt.Println("========================================")

        return
    }

    fmt.Println( body )

    fmt.Println("------------- Body decoded -------------")

    // ========================================
    // ========================================
    // ============= DO THE WORK ==============
    // ========================================
    // ========================================

    var buffer = body.Buffer

    // ========================================
    // ========================================
    // ============ / DO THE WORK =============
    // ========================================
    // ========================================

    fmt.Println("------------- Action done --------------")

    js, err := json.Marshal( buffer )

    if err != nil {

        http.Error(w, err.Error(), http.StatusInternalServerError)

        fmt.Println("========================================")

        return
    }

    w.Header().Set("Content-Type", "application/json")

    w.Write(js)

    fmt.Println("------------- Send response ------------")

    fmt.Println("========================================")
}




















































