package main

import (

)

type Step struct {
    Name string `json:"name"`
    Method string `json:"method"`
    Url string `json:"url"`
    EndStep bool `json:"end_step"`
}
