package main

import "time"

type Terminal struct {
	ID             int
	MAC            string
	CPUUsage       float64
	MemUsage       float64
	LastModifyTime time.Time
}

type ChangeLog struct {
	TerminalID     int
	OldMAC         string
	NewMAC         string
	Timestamp      int64
	LastModifyTime time.Time
}
