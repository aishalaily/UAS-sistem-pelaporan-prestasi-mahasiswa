package main

import (
    "fmt"
    "uas-go/utils"
)

func tes() {
    h, _ := utils.HashPassword("admin123")
    fmt.Println(h)
}
