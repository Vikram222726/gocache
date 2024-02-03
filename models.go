package main

import "time"

type DoublyLinkedList struct{
	Key string `json:"key"`
	Value string `json:"value"`
	TTL time.Time `json:"item_ttl"`
	LeftPointer *DoublyLinkedList
	RightPointer *DoublyLinkedList
}

type ResultSet struct{
	Key string `json:"item_key"`
	Value string `json:"item_value"`
	TTL time.Time `json:"item_entry_time"`
}