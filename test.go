package main

import (
    "fmt"
    "uas-go/utils"
)

func main() {
    h, _ := utils.HashPassword("admin123")
    fmt.Println(h)
}
