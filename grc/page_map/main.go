package main

import (
    "flag"
    "fmt"
)

import "github.com/jetnoli/go-router/router"

func main(){
    flag.Parse()
    path := flag.Arg(0)
    fmt.Println("path ", path)
    router.CreateAssetsFile("asset_map.json", path)
}