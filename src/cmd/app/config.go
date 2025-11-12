package main

import "time"

type Config struct {
	HTTPServer HTTPServer
	Database   Database
}

type HTTPServer struct {
	Addr            string
	ShutdownTimeout time.Duration
}

type Database struct {
}
