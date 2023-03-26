package models

import "gorm.io/gorm"

func (ConfiguracoesEmpresa) TableName() string {
	return "configuracoes"
}

type ConfiguracoesEmpresa struct {
	gorm.Model
	EmpresaID            int
	DiasBoletoAutomatico int
}
