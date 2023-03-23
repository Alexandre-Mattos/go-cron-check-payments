package models

import "gorm.io/gorm"

type Cobranca struct {
	gorm.Model
	ContaID uint
	Status  string
}
