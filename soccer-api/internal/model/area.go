package model

type Area struct {
	BaseModel

	ExternalId  string `gorm:"type:varchar(255);uniqueIndex"`
	Name        string `gorm:"type:varchar(255)"`
	CountryCode string `gorm:"type:varchar(10)"`
}
