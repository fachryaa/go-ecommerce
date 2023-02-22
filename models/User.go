package models

type User struct {
	UserId   int64  `gorm:"primaryKey" json:"userId"`
	UserName string `gorm:"varchar(300);unique;not null" json:"username"`
	Password string `gorm:"varchar(300);not null" json:"password"`
	FullName string `gorm:"varchar(300);not null" json:"fullName"`
	Address  string `gorm:"varchar(300);not null" json:"address"`
	Phone    string `gorm:"varchar(300);not null" json:"phone"`
}
