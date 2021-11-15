package main

import (
    "fmt"
    "form3-organisation-client/accounts"
)
        
func main() {
    fmt.Println("Form3 Organisation Client")
    
    // account_data, err := accounts.Fetch("d5fc00de-1e0d-4bce-a486-63b146fe5093")
    // if err != nil {
    //     panic(err)
    // }

    // fmt.Println(account_data.Attributes.Name)

    err := accounts.Delete("7ea9c39e-3323-418e-9c9e-80d2767b587f", 0)
    if err != nil {
        panic(err)
    }
}

