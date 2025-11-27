package main

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title          string `gorm:"type:varchar(255);not null" json:"title"`
	AvailableSeats int    `gorm:"column:available_seats;not null;default:0" json:"available_seats"`
}
