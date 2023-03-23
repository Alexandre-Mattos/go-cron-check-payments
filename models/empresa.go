package models

import (
	"gorm.io/gorm"
)

type Empresa struct {
	gorm.Model
	ID                   int
	AsaasKey             string
	Configuracoes        ConfiguracoesEmpresa
	DiasBoletoAutomatico string
}
