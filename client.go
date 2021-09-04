package main

import (
	"sync"
	"time"
)

type Status int

const (
	Request_Record_Limit = 400
	Login_Record_Limit   = 200

	UNBAN Status = 0
	BAN   Status = 1
)

func (s Status) String() string {
	return [...]string{"UNBAN", "BAN"}[s]
}

func (s Status) EnumIndex() int {
	return int(s)
}

type client struct {
	ip              string
	status          Status
	request_records []time.Time
	login_records   []time.Time
	sync.RWMutex
}

func NewClient(ip string, times ...time.Time) *client {
	c := &client{}
	c.ip = ip
	for i, t := range times {
		switch i {
		case 0:
			c.request_records = append(c.request_records, t)
		case 1:
			c.login_records = append(c.login_records, t)
		}
	}
	return c
}

func (c *client) GetStatus() Status {
	c.RLock()
	defer c.RUnlock()
	return c.status
}

func (c *client) SetStatus(status Status) {
	c.Lock()
	c.status = status
	c.Unlock()
}

func (c *client) CountRequestTimes(period time.Duration, now time.Time) int {
	var count int
	c.RLock()
	for _, t := range c.request_records {
		if now.Sub(t) <= period && now.Sub(t) >= 0 {
			count += 1
		}
	}
	c.RUnlock()
	return count
}

func (c *client) CountLoginTimes(period time.Duration, now time.Time) int {
	var count int
	c.RLock()
	for _, t := range c.login_records {
		if now.Sub(t) <= period && now.Sub(t) >= 0 {
			count += 1
		}
	}
	c.RUnlock()
	return count
}

func (c *client) AddRequestRecords(req time.Time) {
	c.Lock()
	c.request_records = append([]time.Time{req}, c.request_records...)
	if len(c.request_records) > Request_Record_Limit {
		c.request_records = c.request_records[0:Request_Record_Limit]
	}
	c.Unlock()
}

func (c *client) AddLoginRecords(login time.Time) {
	c.Lock()
	c.login_records = append([]time.Time{login}, c.login_records...)
	if len(c.login_records) > Login_Record_Limit {
		c.login_records = c.login_records[0:Login_Record_Limit]
	}
	c.Unlock()
}

func (c *client) Unban() string {
	c.SetStatus(UNBAN)
	return UNBAN.String() + "," + c.ip
}

func (c *client) Ban() string {
	c.SetStatus(BAN)
	return BAN.String() + "," + c.ip
}
