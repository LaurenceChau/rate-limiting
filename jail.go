package main

import (
	"fmt"
	"sync"
	"time"
)

type prisoner struct {
	client     *client
	entry_time time.Time
}

type room struct {
	prisioners  map[string]prisoner
	lock_period time.Duration
	sync.RWMutex
}

var (
	room_10mins, room_1hour, room_2hours room
	rooms                                map[time.Duration]*room
)

func init() {
	min10, _ := time.ParseDuration("10m")
	hour1, _ := time.ParseDuration("1h")
	hour2, _ := time.ParseDuration("2h")
	room_10mins = room{prisioners: map[string]prisoner{}, lock_period: min10}
	room_1hour = room{prisioners: map[string]prisoner{}, lock_period: hour1}
	room_2hours = room{prisioners: map[string]prisoner{}, lock_period: hour2}
	rooms = make(map[time.Duration]*room)
}

func GetAllRooms() map[time.Duration]*room {
	rooms[room_10mins.lock_period] = &room_10mins
	rooms[room_1hour.lock_period] = &room_1hour
	rooms[room_2hours.lock_period] = &room_2hours
	return rooms
}

func (r *room) AddPrisoner(p prisoner) {
	r.Lock()
	r.prisioners[p.client.ip] = p
	r.Unlock()
	fmt.Printf("%d,%s\n", p.entry_time.Unix(), p.client.Ban())
}

func (r *room) ReleasePrisoners(now time.Time) int {
	var count int
	r.Lock()
	for name, p := range r.prisioners {
		if now.Sub(p.entry_time) >= r.lock_period {
			count += 1
			fmt.Printf("%d,%s\n", now.Unix(), p.client.Unban())
			delete(r.prisioners, name)
		}
	}
	r.Unlock()
	return count
}

func (r *room) GetNoOfPrisoner() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.prisioners)
}
