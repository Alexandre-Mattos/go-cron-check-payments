package models

import "gorm.io/gorm"

type ConfiguracoesEmpresa struct {
	gorm.Model
	EmpresaID            int
	DiasBoletoAutomatico int
}
