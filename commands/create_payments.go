package commands

import (
	"fmt"
	"go-cron-check-payments/asaas"
	"go-cron-check-payments/database"
	"go-cron-check-payments/logger"
	"go-cron-check-payments/models"
	"strconv"
	"strings"

	"log"
	"time"
)

func CreatePayments() {

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	var contasID []string

	var empresas []models.Empresa
	empresaInicio := 1
	dataGeracao := time.Now()

	db.Model(&models.Empresa{}).
		Where("empresas.id >= ?", empresaInicio).
		Where("status = ?", "A").
		Where("asaas_key IS NOT NULL").
		Joins("JOIN configuracoes ON empresas.id = configuracoes.empresa_id").
		Select("empresas.id, empresas.asaas_key, configuracoes.dias_boleto_automatico").
		Group("empresas.id").
		Find(&empresas)

	if err != nil {
		panic(err)
	}

	for _, empresa := range empresas {

		var diasBoletoAutomatico int

		if empresa.DiasBoletoAutomatico != 0 {
			diasBoletoAutomatico = empresa.DiasBoletoAutomatico
		} else {
			diasBoletoAutomatico = 10
		}

		diasCobrancaAutomatica := dataGeracao.AddDate(0, 0, diasBoletoAutomatico).Format("2006-01-02")

		fmt.Println(diasCobrancaAutomatica)

		var contas []models.Conta
		db.Model(&models.Conta{}).
			Where("contas.empresa_id = ?", empresa.ID).
			Where("data_vencimento = ?", diasCobrancaAutomatica).
			Where("situacao = ?", "em_aberto").
			Where("contas.tipo = ?", "recebimento").
			Where("data_cancelamento IS NULL").
			Where("contas.origem_id IS NOT NULL").
			Where("geracao_automatica_boleto = ?", true).
			Where("NOT EXISTS (SELECT 1 FROM boletos WHERE boletos.conta_id = contas.id AND boletos.status != 0)").
			Joins("LEFT JOIN locacoes ON locacoes.id = contas.origem_id").
			Joins("JOIN locacao_inquilinos ON locacao_inquilinos.locacao_id = locacoes.id").
			Joins("LEFT JOIN clientes ON clientes.id = contas.cliente_id").
			Where("locacao_inquilinos.forma_cobranca = ?", "boleto_bancario").
			Select("contas.valor, contas.descricao, contas.empresa_id, contas.origem_id, clientes.customer_id as customer_id, contas.id, contas.data_vencimento").
			Find(&contas)

		for _, conta := range contas {
			fmt.Println(conta.ID)

			var pointer *models.Cobranca = new(models.Cobranca)
			pointer, success := asaas.CreatePayment(conta)
			if pointer != nil {
				cobranca := *pointer
				if success {
					contasID = append(contasID, strconv.Itoa(cobranca.ContaID))
				}
			}
		}

		if len(contasID) >= 1 {
			logger.Send("Contas geradas em "+time.Now().String()+": "+strings.Join(contasID, ","), "success")
		} else {
			logger.Send("Nenhuma cobrança gerada em: "+time.Now().Format("2006-01-02 15:04:05"), "warning")
		}
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
