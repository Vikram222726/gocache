package main

import (
	"fmt"
	"net/http"
	"sync"
)

var header *DoublyLinkedList
var tail *DoublyLinkedList
const portNumber = ":6060"
const cacheSize int = 5
var cacheItemsCount int = 0
var mutexLock sync.Mutex
var cacheMap = make(map[string]*DoublyLinkedList)

// This function is the entry point for implementing our cache...
func main(){
	fmt.Println("Let's Build a simple caching service using go lang...")

	fmt.Println("Starting the server on port 6060...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Welcome to simple go lang cache")
	})

	http.HandleFunc("/putItem", InsertItemToCache)
	http.HandleFunc("/getItem", GetSingleCacheItem)
	http.HandleFunc("/getItems", GetAllCacheItems)
	http.HandleFunc("/updateItem", UpdateCacheItem)
	http.HandleFunc("/deleteItem", DeleteCacheItem)

	http.ListenAndServe(portNumber, nil)
}