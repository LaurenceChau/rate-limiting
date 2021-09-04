package main

import (
	"testing"
	"time"
)

func TestGetStatus(t *testing.T) {
	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime)

	// get status
	status := client1.GetStatus()

	// verify
	if status != UNBAN {
		t.Errorf("status = %d; want UNBAN", status)
	}
}

func TestSetStatus(t *testing.T) {

	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime)

	// set user status
	client1.SetStatus(BAN)

	// verify
	if client1.GetStatus() != BAN {
		t.Errorf("User Status = %s; want BAN", client1.GetStatus())
	}

}

func TestCountRequestTimes(t *testing.T) {
	// create client
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime)

	// count request times
	currentTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	period, _ := time.ParseDuration("1m")
	reqTimes := client1.CountRequestTimes(period, currentTime)

	// verify
	if reqTimes != 1 {
		t.Errorf("reqTimes = %d; want 1", reqTimes)
	}
}

func TestCountLoginTimes(t *testing.T) {
	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	loginTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime, loginTime)

	// count login times
	currentTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	period, _ := time.ParseDuration("1m")
	LoginTimes := client1.CountLoginTimes(period, currentTime)
	if LoginTimes != 1 {
		t.Errorf("LoginTimes = %d; want 0", LoginTimes)
	}
}

func TestAddRequestRecords(t *testing.T) {
	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime)

	// user request => add request times
	requestTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1.AddRequestRecords(requestTime)

	// verify
	currentTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	period, _ := time.ParseDuration("1m")
	count := client1.CountRequestTimes(period, currentTime)
	if count != 2 {
		t.Errorf("Request Times = %d; want 2", count)
	}
}

func TestAddLoginRecords(t *testing.T) {

	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	loginTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime, loginTime)

	// add login record
	loginRequest, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1.AddLoginRecords(loginRequest)

	// verify
	currentTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	period, _ := time.ParseDuration("1m")
	count := client1.CountLoginTimes(period, currentTime)
	if count != 2 {
		t.Errorf("Login Times = %d; want 1", count)
	}
}

func TestUnbanUser(t *testing.T) {

	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime)

	// // Release User
	result := client1.Unban()

	// verify1
	if result != "UNBAN,222.219.188.0" {
		t.Errorf("result = %s; want UNBAN,222.219.188.0", result)
	}
	// verify2
	if client1.GetStatus() != UNBAN {
		t.Errorf("status = %s; want UNBAN", client1.GetStatus())
	}
}

func TestBanUser(t *testing.T) {
	// create user
	timeLayout := "2/Jan/2006:15:04:05 -0700"
	accessTime, _ := time.Parse(timeLayout, "31/Dec/2018:23:55:02 +0800")
	client1 := NewClient("222.219.188.0", accessTime)

	// // Release User
	result := client1.Ban()

	// verify1
	if result != "BAN,222.219.188.0" {
		t.Errorf("result = %s; want BAN,222.219.188.0", result)
	}
	// verify2
	if client1.GetStatus() != BAN {
		t.Errorf("status = %s; want BAN", client1.GetStatus())
	}
}
