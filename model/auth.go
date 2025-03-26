package model

type Login struct {
	Email    string `gorm:"column:email;type:varchar(255);NOT NULL" json:"email"`
	Password string `json:"password"`
}
