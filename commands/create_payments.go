package commands

import (
	"go-cron-check-payments/database"
	"go-cron-check-payments/logger"
	"go-cron-check-payments/models"

	"log"
	"time"
)

func CreatePayments() error {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	var empresas []models.Empresa
	empresaInicio := 1
	dataGeracao := time.Now()

	db.Where("empresas.id >= ?", empresaInicio).
		Where("status = ?", "A").
		Where("asaas_key IS NOT NULL").
		Joins("JOIN configuracoes ON empresas.id = configuracoes.empresa_id").
		Select("empresas.id, empresas.asaas_key, configuracoes.dias_boleto_automatico").
		Group("empresas.id").
		Find(&empresas)

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Close()

	log.Info("info message")

	for _, empresa := range empresas {
		diasCobrancaAutomatica := dataGeracao.AddDate(0, 0, time.Parse(empresa.DiasBoletoAutomatico)).Format("2006-01-02")

		var contas []models.Conta
		db.Where("empresa_id = ?", empresa.ID).
			Where("data_vencimento = ?", diasCobrancaAutomatica).
			Where("status = ?", "A").
			Where("geracao_automatica_boleto = ?", true).
			Where("NOT EXISTS (SELECT 1 FROM cobrancas WHERE cobrancas.conta_id = contas.id AND cobrancas.status != 'canceled')").
			Joins("JOIN locacoes ON locacoes.id = contas.locacao_id").
			Joins("JOIN inquilinos ON inquilinos.locacao_id = locacoes.id").
			Where("inquilinos.forma_cobranca = ?", "boleto_bancario").
			Find(&contas)

		for _, conta := range contas {
			var cobrancas []models.Cobranca
			db.Where("conta_id = ?", conta.ID).
				Where("status != ?", "canceled").
				Find(&cobrancas)

			// Process cobrancas
			// ...
		}
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()

	return nil
}
