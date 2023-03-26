package models

import "gorm.io/gorm"

func (Conta) TableName() string {
	return "contas"
}

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
	CustomerId              string
	Valor                   float32
	Saldo                   float32
	Descricao               string
}

type ContasResponse struct {
	Data []Conta
}

func (Multa) TableName() string {
	return "multas"
}

type Multa struct {
	gorm.Model
	Porcentagem float32
}

func (Juros) TableName() string {
	return "juros"
}

type Juros struct {
	gorm.Model
	Porcentagem float32
}

func (Desconto) TableName() string {
	return "descontos"
}

type Desconto struct {
	DescontarAteVencimento bool
	PrazoAteVenciemnto     int32
	Valor                  float32
}
