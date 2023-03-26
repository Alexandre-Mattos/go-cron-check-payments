package models

import (
	"gorm.io/gorm"
)

func (Empresa) TableName() string {
	return "empresas"
}

type Empresa struct {
	gorm.Model
	ID                   int
	AsaasKey             string
	Configuracoes        ConfiguracoesEmpresa
	DiasBoletoAutomatico int
}

func (EmpresaTaxa) TableName() string {
	return "empresa_taxas"
}

type EmpresaTaxa struct {
	gorm.Model
	WalletId        string
	ValorReal       float32
	ValorPercentual float32
	DeletedAt       gorm.DeletedAt `gorm:"-:all"`
}
