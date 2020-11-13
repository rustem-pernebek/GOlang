package models

//var ErrNoRecord = errors.New("models: no matching record found")

type Book struct {
	Id string
	Name string
	Author string
	Price int

}

type User struct {
	Id int
	Login string
	Passrword string
	Name string
	Phone string
}

type Question struct {
	Id      int
	Quest   string
	Answers []*Answer

}

type Answer struct {
	Id     int
	Answer string
}

type Cart struct {
	User  User
	books *[]Book
}
