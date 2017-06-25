package main

import (
    "fmt"
    "encoding/json"
    "io"
    "net/http"
    "time"
)

var (
    MaxWorker       = 20  //os.Getenv("MAX_WORKERS")
    MaxDatabaseWorker       = 20  //os.Getenv("MAX_WORKERS")
    MaxQueue        = 5 //os.Getenv("MAX_QUEUE")
    MaxLength int64 = 20480
)


func main() {

    fmt.Println("starting server http")

    http.HandleFunc("/starting", postStarting)

    http.HandleFunc("/onServer", postOnServer)

//    http.HandleFunc("/onInterne", postOnInterne)
//    http.HandleFunc("/ending", postEnding)

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

type HttpOut struct {
    Name string
    Arguments JsonB
    Buffer JsonB
}


type Profile struct {
      Name    string
    Hobbies []string
}

func postStarting(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    fmt.Println("Request received")

    fmt.Println("___________________________________")
    // Read the body into a string for json decoding
    var content = &HttpOut{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    fmt.Println( content )
    fmt.Println("___________________________________")


    profile := Profile{"Alex", []string{"snowboarding", "programming"}}

    js, err := json.Marshal(profile)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)





//    w.WriteHeader(http.StatusOK)
}

func postOnServer(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    fmt.Println("Request received")

    fmt.Println("___________________________________")
    // Read the body into a string for json decoding
    var content = &HttpOut{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    fmt.Println( content )
    fmt.Println("___________________________________")


    profile := Profile{"Noemi", []string{"toto", "success"}}

    js, err := json.Marshal(profile)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)





//    w.WriteHeader(http.StatusOK)
}


//// ================================================= //
//// ================================================= //
//
//func postOnServer(w http.ResponseWriter, r *http.Request) {
//
//    if r.Method != "POST" {
//        w.WriteHeader(http.StatusMethodNotAllowed)
//        return
//    }
//
//    // Read the body into a string for json decoding
//    var content = &TaskCollection{}
//    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
//    if err != nil {
//        fmt.Errorf("an error occured while deserializing message")
//        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//        w.WriteHeader(http.StatusBadRequest)
//        return
//    }
//
//    fmt.Println("Request received")
//
//    // Go through each payload and queue items individually to be posted to S3
//    for _, task := range content.Tasks {
//
//        // Push the work onto the queue.
//        TaskQueue <- task
//
//        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
//    }
//
//    w.WriteHeader(http.StatusOK)
//}
//
//// ================================================= //
//// ================================================= //
//
//func postOnInterne(w http.ResponseWriter, r *http.Request) {
//
//    if r.Method != "POST" {
//        w.WriteHeader(http.StatusMethodNotAllowed)
//        return
//    }
//
//    // Read the body into a string for json decoding
//    var content = &TaskCollection{}
//    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
//    if err != nil {
//        fmt.Errorf("an error occured while deserializing message")
//        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//        w.WriteHeader(http.StatusBadRequest)
//        return
//    }
//
//    fmt.Println("Request received")
//
//    // Go through each payload and queue items individually to be posted to S3
//    for _, task := range content.Tasks {
//
//        // Push the work onto the queue.
//        TaskQueue <- task
//
//        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
//    }
//
//    w.WriteHeader(http.StatusOK)
//}
//
//// ================================================= //
//// ================================================= //
//
//func postEnding(w http.ResponseWriter, r *http.Request) {
//
//    if r.Method != "POST" {
//        w.WriteHeader(http.StatusMethodNotAllowed)
//        return
//    }
//
//    // Read the body into a string for json decoding
//    var content = &TaskCollection{}
//    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
//    if err != nil {
//        fmt.Errorf("an error occured while deserializing message")
//        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//        w.WriteHeader(http.StatusBadRequest)
//        return
//    }
//
//    fmt.Println("Request received")
//
//    // Go through each payload and queue items individually to be posted to S3
//    for _, task := range content.Tasks {
//
//        // Push the work onto the queue.
//        TaskQueue <- task
//
//        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
//    }
//
//    w.WriteHeader(http.StatusOK)
//}
//
