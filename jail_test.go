package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetAllRooms(t *testing.T) {
	rooms := GetAllRooms()
	fmt.Println("TestGetAllRooms:", rooms)
}

func TestAddPrisoners(t *testing.T) {
	rooms := GetAllRooms()
	min10, _ := time.ParseDuration("10m")

	// new a client
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client_x := NewClient("222.219.188.0", accessTime)

	// client become prisoner
	p := prisoner{client: client_x, entry_time: accessTime}

	// add prisoner to jail
	rooms[min10].AddPrisoner(p)

	// verify: there should have 1 prisoner
	if len(rooms[min10].prisioners) != 1 {
		t.Errorf("No of prisoners = %d; want 1", len(rooms[min10].prisioners))
	}
}

func TestReleasePrisoners(t *testing.T) {
	rooms := GetAllRooms()
	min10, _ := time.ParseDuration("10m")

	// new a client
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client_x := NewClient("222.219.188.0", accessTime)

	// client become prisoner
	p := prisoner{client: client_x, entry_time: accessTime}

	// add prisoner to jail
	rooms[min10].AddPrisoner(p)

	// release prisoner
	currentTime, _ := time.Parse(timeLayout, "01/Jan/2019:01:55:02 +0800")
	rooms[min10].ReleasePrisoners(currentTime)
}
