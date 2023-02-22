package models

type Product struct {
	ProductId       int64  `gorm:"primaryKey" json:"productId"`
	ProductName     string `gorm:"varchar(300);not null" json:"productName"`
	ProductCategory string `gorm:"varchar(300);not null" json:"productCategory"`
	ProductPrice    uint64 `gorm:"int;not null" json:"productPrice"`
	ProductStock    uint64 `gorm:"int;not null" json:"productStock"`
}
