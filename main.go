package main

import (
    "github.com/memmaker/go/fxtools"
    "github.com/memmaker/go/recfile"
)

type MyRecord struct {
    Name string
    Age int
    Height float64
    IsStudent bool
    Address Address
}

type Address struct {
    Street string
    City string
    Zip int
}


func main() {
    recEncoder := recfile.NewEncoder(fxtools.MustCreate("test.rec"))

    err := recEncoder.Encode(MyRecord{
        Name: "John Doe",
        Age: 25,
        Height: 5.9,
        IsStudent: true,
        Address: Address{
            Street: "123 Main St",
            City: "Anytown",
            Zip: 12345,
        },
    })
    if err != nil {
        panic(err)
    }
}