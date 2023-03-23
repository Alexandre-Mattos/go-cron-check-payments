package models

import "gorm.io/gorm"

type Conta struct {
	gorm.Model
	ID                      int
	EmpresaID               int
	DataVencimento          string
	GeracaoAutomaticaBoleto bool
	ReferenciaInicial       string
	ReferenciaFormatada     string
	OrigemType              string
	OrigemID                int
}

type ContasResponse struct {
	Data []Conta
}
