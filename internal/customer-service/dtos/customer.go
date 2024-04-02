package dtos

import "time"

type CustomerDto struct {
	Id          uint      `json:"id"`
	Firstname   string    `json:"firstname"`
	Middlename  *string   `json:"middlename"`
	Lastname    string    `json:"lastname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}
