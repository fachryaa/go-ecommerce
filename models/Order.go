package models

type Order struct {
	OrderId    int64  `gorm:"primaryKey" json:"orderId"`
	UserId     int64  `gorm:"int" json:"userId"`
	TotalPrice uint64 `gorm:"int" json:"totalPrice"`
	FullName   string `gorm:"varchar(300)" json:"fullName"`
	Address    string `gorm:"varchar(300)" json:"address"`
	Phone      string `gorm:"varchar(300)" json:"phone"`
}
