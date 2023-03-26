package models

import "gorm.io/gorm"

func (Cobranca) TableName() string {
	return "boletos"
}

type Cobranca struct {
	gorm.Model
	ID                int
	ContaID           int
	Status            bool
	TransactionID     string
	UrlGerarCobranca  string
	EmpresaID         int
	GeracaoAutomatica bool
	Situacao          string
}

func (CobrancaTaxa) TableName() string {
	return "cobranca_taxas"
}

type CobrancaTaxa struct {
	gorm.Model
	EmpresaID       int
	BoletoId        int
	ValorReal       float32
	ValorPercentual float32
	Tipo            string
}
