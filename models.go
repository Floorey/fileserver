package main

import "time"

type File struct {
	ID         uint `grom:"primaryKey"`
	Filename   string
	Filepath   string
	UploadDate time.Time
}
