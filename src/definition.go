package main

type UserInfo struct {
	UserName  string
	StudentID string
}

type LoginInfo struct {
	StudentID         string
	LoggingTime       string
	SubmittingTime    string
	SubmittedFileName string
}

type AdminInfo struct {
	UserName string
	Password string
}
