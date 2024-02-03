package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func removeTailNode(){
	newTailNode := tail.LeftPointer
	tail.LeftPointer = nil
	newTailNode.RightPointer = nil
	tail = newTailNode
	cacheItemsCount--
}

// This function creates a doubly linked list
func InsertItemToCache(w http.ResponseWriter, r *http.Request){
	fmt.Println("Started inserting new item to cache")
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Data not sent for creating an item")
	}

	var newNode DoublyLinkedList
	err := json.NewDecoder(r.Body).Decode(&newNode)
	if err != nil {
		log.Fatal(err)
	}

	newNode.TTL = time.Now()

	itemKey := newNode.Key
	itemVal := newNode.Value
	itemTTL := newNode.TTL

	itemNode, itemPresent := cacheMap[itemKey]
	if itemPresent{
		itemNode.Value = itemVal
		itemNode.TTL = itemTTL
		json.NewEncoder(w).Encode(ResultSet{Key: itemKey, Value: itemVal, TTL: itemTTL})
		return
	}

	mutexLock.Lock()

	if header == nil && tail == nil{
		// Means the cache is empty currently..
		header, tail = &newNode, &newNode
		json.NewEncoder(w).Encode(ResultSet{Key: itemKey, Value: itemVal, TTL: itemTTL})
		cacheMap[newNode.Key] = &newNode
		cacheItemsCount++
		mutexLock.Unlock()
		return
	}else if cacheItemsCount == cacheSize {
		// Cache is full we'll have to remove the last item from tail and also from map
		delete(cacheMap, tail.Key)
		removeTailNode()
	}

	// Add this new node inside the Doubly Linked list
	newNode.RightPointer = header
	header.LeftPointer = &newNode
	header = &newNode
	cacheMap[newNode.Key] = &newNode
	cacheItemsCount++

	err = json.NewEncoder(w).Encode(ResultSet{Key: itemKey, Value: itemVal, TTL: itemTTL})
	if err != nil {
		log.Fatal(err)
	}
	mutexLock.Unlock()
}

// This function will return all items from the cache...
func GetAllCacheItems(w http.ResponseWriter, r *http.Request){
	fmt.Println("Fetching all items from cache")
	w.Header().Set("Content-Type", "application/json")

	headPtr := header
	resultSet := []ResultSet{}
	for headPtr != nil {
		nodeData := ResultSet{
			Key: headPtr.Key,
			Value: headPtr.Value,
			TTL: headPtr.TTL,
		}
		resultSet = append(resultSet, nodeData)
		headPtr = headPtr.RightPointer
	}

	err := json.NewEncoder(w).Encode(resultSet)
	if err != nil {
		log.Fatal(err)
	}
}

// This function will fetch single item from map and return the result
func GetSingleCacheItem(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	keyValue := r.FormValue("key")

	fmt.Println("Fetching and returning value for key: ", keyValue)
	itemNode, itemPresent := cacheMap[keyValue]

	if !itemPresent{
		json.NewEncoder(w).Encode(fmt.Sprintf("Key: %s not present in cache", keyValue))
		return
	}

	if itemNode != header{
		leftNode := itemNode.LeftPointer
		rightNode := itemNode.RightPointer
		if leftNode != nil{
			leftNode.RightPointer = rightNode
		}
		if rightNode != nil {
			rightNode.LeftPointer = leftNode
		}
		header.LeftPointer = itemNode
		itemNode.LeftPointer = nil
		itemNode.RightPointer = header
		header = itemNode
	}

	json.NewEncoder(w).Encode(ResultSet{Key: itemNode.Key, Value: itemNode.Value, TTL: itemNode.TTL})
}

// This function will update the key's value in our cache
func UpdateCacheItem(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var updatedData DoublyLinkedList
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updating Cache Item for key:", updatedData.Key)
	
	itemNode, itemPresent := cacheMap[updatedData.Key]
	if !itemPresent{
		json.NewEncoder(w).Encode(fmt.Sprintf("Key: %s not present in cache", updatedData.Key))
		return
	}

	itemNode.Value = updatedData.Value
	itemNode.TTL = time.Now()

	mutexLock.Lock()

	// Add logic to update this node in DLL
	if itemNode != header{
		leftNode := itemNode.LeftPointer
		rightNode := itemNode.RightPointer
		if leftNode != nil{
			leftNode.RightPointer = rightNode
		}
		if rightNode != nil {
			rightNode.LeftPointer = leftNode
		}
		header.LeftPointer = itemNode
		itemNode.LeftPointer = nil
		itemNode.RightPointer = header
		header = itemNode
	}

	json.NewEncoder(w).Encode(ResultSet{Key: itemNode.Key, Value: itemNode.Value})
	mutexLock.Unlock()
}

func DeleteCacheItem(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	keyValue := r.FormValue("key")

	fmt.Printf("Deleting key: %s from cache", keyValue)
	itemNode, itemPresent := cacheMap[keyValue]
	fmt.Println(keyValue, itemPresent)

	if !itemPresent{
		json.NewEncoder(w).Encode(fmt.Sprintf("Key: %s not present in cache", keyValue))
		return
	}

	mutexLock.Lock()

	// Removed key from our cache map
	delete(cacheMap, keyValue)

	// Add logic to delete key from DLL
	leftNode := itemNode.LeftPointer
	rightNode := itemNode.RightPointer
	if header == itemNode{
		header = rightNode
	}
	if tail == itemNode{
		tail = leftNode
	}
	if leftNode != nil{
		leftNode.RightPointer = rightNode
	}
	if rightNode != nil{
		rightNode.LeftPointer = leftNode
	}
	itemNode.LeftPointer, itemNode.RightPointer = nil, nil

	cacheItemsCount--

	json.NewEncoder(w).Encode(fmt.Sprintf("Successfully deleted %s Key from cache", keyValue))

	mutexLock.Unlock()
}