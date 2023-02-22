package models

type Cart struct {
	CartId     int64  `gorm:"primaryKey" json:"cartId"`
	ProductId  int64  `gorm:"int;not null" json:"productId"`
	UserId     int64  `gorm:"int" json:"userId"`
	Amount     uint64 `gorm:"int;not null" json:"amount"`
	TotalPrice uint64 `gorm:"int" json:"totalPrice"`
}
