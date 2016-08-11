package main

import (
	"fmt"
	"net/http"
)

type criticalServices []struct {
	Node        string `json:"Node"`
	CheckID     string `json:"CheckID"`
	Name        string `json:"Name"`
	Status      string `json:"Status"`
	Notes       string `json:"Notes"`
	Output      string `json:"Output"`
	ServiceID   string `json:"ServiceID"`
	ServiceName string `json:"ServiceName"`
}

func main() {
	resp, err := http.Get("http://example.com/")
	if err != nil {

	}

	fmt.Print("hello")
}
