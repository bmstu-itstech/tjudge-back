package app

type User struct {
	Id       	int
	Username 	string
	Fullname 	string
} 

type Team struct {
	Code     	string
	Name     	string
	Leader 		*User
	Contest 	int
	MaxSize 	int
	Members 	[]*User
}
