package main

import (
    "io"
    "fmt"
    "time"
    "net/http"
    "encoding/json"
)

type StepInDataStruct struct {
    Name        string
    Arguments   JsonB
    Buffer      JsonB
}

type StepOutDataStruct struct {
    Buffer      *JsonB
    Interval    *int
    Step        *string
    Comment     *string
}

func postHandler(w http.ResponseWriter, r *http.Request) {

    // Setup the response header
    w.Header().Set("Content-Type", "application/json")

    fmt.Println("========================================")

    fmt.Println("----------- Request received -----------")

    if r.Method != "POST" {

        w.WriteHeader(http.StatusMethodNotAllowed)

        return
    }

    fmt.Println("------------- Request post -------------")

    // Read data from body request
    var body = &StepInDataStruct{}

    // Decode body json data
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&body)

    if err != nil {

        fmt.Errorf("an error occured while deserializing message")

        w.WriteHeader(http.StatusBadRequest)

        fmt.Println("========================================")

        return
    }

    fmt.Println("------------- Body decoded -------------")

    // Show the request body data
    fmt.Println( body )

    fmt.Println("------------- Action begin -------------")








//    //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//    //  200  = next                                                         //
//    //______________________________________________________________________//
//
//    var buffer = body.Buffer
//
//    buffer["steps"] = []string{ r.URL.String() }
//
//    var stepOutData = &StepOutDataStruct{
//        Buffer: &buffer }
//


//    //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//    //  301  = next to '...' step or/and next in '...' interval of seconds  //
//    //______________________________________________________________________//
//
//    w.WriteHeader(http.StatusMovedPermanently)
//
//    var interval = 24 * 60 * 60 // 1 day
//    var step = "ending"
//    var buffer = body.Buffer
//
//    buffer["steps"] = []string{ r.URL.String() }
//
//    var stepOutData = &StepOutDataStruct{
//        Buffer: &buffer,
//        Interval: &interval,
//        Step: &step }


    //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
    //  420  = cancelled                                                    //
    //______________________________________________________________________//

    w.WriteHeader(420)

    var comment = "commentaire from api"
    var buffer = body.Buffer

    buffer["steps"] = []string{ r.URL.String() }

    var stepOutData = &StepOutDataStruct{
        Buffer: &buffer,
        Comment: &comment }


















    fmt.Println("------------- Action done --------------")

    // Re-encode the response
    js, err := json.Marshal( stepOutData )

    if err != nil {

        http.Error(w, err.Error(), http.StatusInternalServerError)

        fmt.Println("========================================")

        return
    }

    fmt.Println("------------- Send response ------------")

    // Write the json response
    w.Write(js)

    fmt.Println("========================================")
}

// ================================================= //
// ================================================= //




































































func main() {

    fmt.Println("starting server http")

    http.HandleFunc("/starting", postStarting)
    http.HandleFunc("/onServer", postOnServer)
    http.HandleFunc("/onInterne", postOnInterne)
    http.HandleFunc("/ending", postEnding)

    http.HandleFunc("/create/starting", postHandler)
    http.HandleFunc("/create/onServer", postHandler)
    http.HandleFunc("/create/onInterne", postHandler)
    http.HandleFunc("/create/ending", postHandler)

    http.HandleFunc("/update/starting", postHandler)
    http.HandleFunc("/update/onServer", postHandler)
    http.HandleFunc("/update/onInterne", postHandler)
    http.HandleFunc("/update/ending", postHandler)

    http.HandleFunc("/delete/starting", postHandler)
    http.HandleFunc("/delete/onServer", postHandler)
    http.HandleFunc("/delete/onInterne", postHandler)
    http.HandleFunc("/delete/ending", postHandler)

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
