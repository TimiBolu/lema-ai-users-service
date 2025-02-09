package models

type Address struct {
	ID      string `gorm:"primaryKey" json:"id"`
	UserID  string `gorm:"unique;not null;index" json:"userId"`
	Street  string `gorm:"not null" json:"street"`
	City    string `gorm:"not null" json:"city"`
	State   string `gorm:"not null" json:"state"`
	ZipCode string `gorm:"not null" json:"zipCode"`
}
