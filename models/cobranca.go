package models

func (Cobranca) TableName() string {
	return "boletos"
}

type Cobranca struct {
	ID                int    `gorm:"column:id"`
	EmpresaID         int    `gorm:"column:empresa_id"`
	ContaID           int    `gorm:"column:conta_id"`
	Status            bool   `gorm:"column:status"`
	TransactionID     string `gorm:"column:transaction_id"`
	UrlGerarCobranca  string `gorm:"column:url_gerar_cobranca"`
	GeracaoAutomatica bool   `gorm:"column:geracao_automatica"`
	Situacao          string `gorm:"column:situacao"`
}

func (CobrancaTaxa) TableName() string {
	return "cobranca_taxas"
}

type CobrancaTaxa struct {
	EmpresaID       int     `gorm:"column:empresa_id"`
	BoletoId        int     `gorm:"column:boleto_id"`
	ValorReal       float32 `gorm:"column:valor_real"`
	ValorPercentual float32 `gorm:"column:valor_percentual"`
	Tipo            string  `gorm:"column:tipo"`
}
