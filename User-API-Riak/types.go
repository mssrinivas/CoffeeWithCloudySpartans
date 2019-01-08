package main

type user struct {
	UserId     	string 	`json:"Userid"`
	Email 	   	string 	`json:"email"`
	UserType	string 	`json:"UserType"`
	Password   	string	`json:"Password"`
}

type Keys struct {
	Keys []string
}