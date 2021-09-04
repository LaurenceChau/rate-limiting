package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	timeLayout                           = "2/Jan/2006:15:04:05 -0700"
	max_allowed_requests_in_1min         = 40
	max_allowed_requests_in_10mins       = 100
	max_allowed_login_requests_in_10mins = 20
)

var (
	clientHistory map[string]*client
)

func init() {
	clientHistory = make(map[string]*client)
	rooms = GetAllRooms()
}

type message struct {
	client     *client
	accessTime time.Time
}

func main() {

	var wg sync.WaitGroup
	wg.Add(3)

	message_channel := make(chan message)
	quick_janitor := make(chan bool)
	quick_worker1 := make(chan bool)
	quick_worker2 := make(chan bool)
	time_channel := make(chan time.Time, 5)
	defer close(message_channel)
	defer close(quick_janitor)
	defer close(time_channel)
	defer close(quick_worker1)
	defer close(quick_worker2)

	go ddosMonitor_t(message_channel, quick_worker1, &wg)
	go ddosMonitor_t(message_channel, quick_worker2, &wg)
	go janitor_t(time_channel, quick_janitor, &wg)

	csvLog, err := os.Open("TestQ1.log")
	defer csvLog.Close()
	if err != nil {
		fmt.Println(err)
	}

	r := csv.NewReader(csvLog)
	r.Comma = ' '
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range records {

		ip := line[0]
		timeStr := strings.Trim(line[3], "[") + " " + strings.Trim(line[4], "]")
		path := strings.Split(line[5], " ")[1]

		accessTime, _ := time.Parse(timeLayout, timeStr)
		time_channel <- accessTime

		// check wheather clint is in the client history
		if client, ok := clientHistory[ip]; ok {
			// only add access records to the client only if the client is in unban state
			if client.GetStatus() == UNBAN {
				if path == "/login" {
					client.AddRequestRecords(accessTime)
					client.AddLoginRecords(accessTime)
				} else {
					client.AddRequestRecords(accessTime)
				}
				// only send to the client to ddosmonitory if the client is in unban state
				message_channel <- message{client: client, accessTime: accessTime}
			}
		} else {
			if path == "/login" {
				client_new := NewClient(ip, accessTime, accessTime)
				clientHistory[ip] = client_new
			} else {
				client_new := NewClient(ip, accessTime)
				clientHistory[ip] = client_new
			}
		}
	}
	quick_janitor <- true
	quick_worker1 <- true
	quick_worker2 <- true
	wg.Wait()
}

// determine wheather a client can be released from jail
func janitor_t(timeCh <-chan time.Time, quitCh <-chan bool, wg *sync.WaitGroup) {
	for {
		select {
		case now := <-timeCh:
			for _, room := range rooms {
				room.ReleasePrisoners(now)
			}
		case <-quitCh:
			wg.Done()
			break
		}
	}
}

// determine wheather a client need to in jail
func ddosMonitor_t(mesgCh <-chan message, quitCh <-chan bool, wg *sync.WaitGroup) {

	min1, _ := time.ParseDuration("1m")
	min10, _ := time.ParseDuration("10m")
	hour1, _ := time.ParseDuration("1h")
	hour2, _ := time.ParseDuration("2h")

	for {
		select {
		case mesg := <-mesgCh:
			switch {
			// ban for 10 minute if request > 40 in the past 1 minute
			case mesg.client.CountRequestTimes(min1, mesg.accessTime) > max_allowed_requests_in_1min:
				if mesg.client.GetStatus() != BAN {
					p := prisoner{client: mesg.client, entry_time: mesg.accessTime}
					rooms[min10].AddPrisoner(p)
				}
			// ban for 1 hour if request > 100 in the past 10 minute
			case mesg.client.CountRequestTimes(min10, mesg.accessTime) > max_allowed_requests_in_10mins:
				if mesg.client.GetStatus() != BAN {
					p := prisoner{client: mesg.client, entry_time: mesg.accessTime}
					rooms[hour1].AddPrisoner(p)
				}
			// ban for 2 hours if login request > 20 in the past 10 minute
			case mesg.client.CountLoginTimes(min10, mesg.accessTime) > max_allowed_login_requests_in_10mins:
				if mesg.client.GetStatus() != BAN {
					p := prisoner{client: mesg.client, entry_time: mesg.accessTime}
					rooms[hour2].AddPrisoner(p)
				}
			}
		case <-quitCh:
			wg.Done()
			break
		}
	}
}
