package asaas

import (
	"errors"
	"go-cron-check-payments/database"
	"go-cron-check-payments/logger"
	"go-cron-check-payments/models"
	"os"
	"strconv"
	"strings"

	sdk "github.com/Alexandre-Mattos/go-asaas-sdk-main"
)

func CreatePayment(conta models.Conta, empresa models.Empresa) (*models.Cobranca, bool) {

	db, err := database.Connect()
	if err != nil {
		logger.Send(err.Error(), "error")
	}

	client := sdk.NewAsaasClient(empresa.AsaasKey)

	paymentBoleto := sdk.PaymentBoleto{
		Customer:    conta.CustomerId,
		DueDate:     conta.DataVencimento,
		Value:       conta.Valor,
		Description: conta.Descricao,
	}

	var empresaTaxa models.EmpresaTaxa
	db.Model(&models.EmpresaTaxa{}).
		Unscoped().
		Where("empresa_id = ?", conta.EmpresaID).
		Find(&empresaTaxa)

	var paymentConfig *sdk.PaymentBoleto = &paymentBoleto
	paymentConfig, err = MakeSplit(paymentBoleto, empresaTaxa)
	if err != nil {
		logger.Send(err.Error(), "error")
	}
	paymentBoleto = *paymentConfig

	var contaMulta models.Multa
	notFound := db.Model(&models.Multa{}).
		Unscoped().
		Where("multavel_type IS NOT NULL").
		Where("multavel_id = ?", conta.OrigemID).
		Select("porcentagem").
		Find(&contaMulta).
		Error

	if notFound != nil {
		paymentConfig, err = MakeFine(paymentBoleto, contaMulta)
		if err != nil {
			logger.Send(err.Error(), "error")
		}
		paymentBoleto = *paymentConfig
	}

	var contaJuro models.Juros
	notFound = db.Model(&models.Juros{}).
		Where("sujeito_juros_type IS NOT NULL").
		Where("sujeito_juros_id = ?", conta.OrigemID).
		Select("porcentagem").
		Find(&contaJuro).
		Error

	if notFound != nil {
		paymentConfig, err = MakeInterest(paymentBoleto, contaJuro)
		if err != nil {
			logger.Send(err.Error(), "error")
		}
		paymentBoleto = *paymentConfig
	}

	var contaDesconto models.Desconto
	notFound = db.Model(&models.Desconto{}).
		Where("descontavel_type IS NOT NULL").
		Where("descontavel_id = ?", conta.OrigemID).
		Select("valor, descontar_ate_vencimento,prazo_ate_vencimento").
		Find(&contaDesconto).
		Error

	if notFound != nil {
		paymentConfig, err = MakeDiscount(paymentBoleto, contaDesconto)
		if err != nil {
			logger.Send(err.Error(), "error")
		}
		paymentBoleto = *paymentConfig
	}

	var cobrancaExistente models.Cobranca
	errFind := db.Model(&models.Cobranca{}).
		Unscoped().
		Where("situacao = ?", "em_aberto").
		Select("transaction_id").
		Last(&cobrancaExistente)

	if errFind.Error != nil {

		boletoResponse, errAPI, err := client.UpdatePaymentBoleto("", cobrancaExistente.TransactionID, paymentBoleto)
		if err != nil {
			logger.Send(err.Error(), "error")
			return nil, false
		}
		if errAPI != nil {
			var errorsAsaas []string
			for _, errorAsaas := range errAPI.Errors {
				errorsAsaas = append(errorsAsaas, errorAsaas.Description)
			}
			logger.Send("["+strconv.Itoa(conta.EmpresaID)+"] Erro na requisição ASAAS da conta ["+strconv.Itoa(conta.ID)+"]: "+strings.Join(errorsAsaas, ","), "error")
		}
		if boletoResponse == nil {
			logger.Send("Response do cadastro do boleto da conta "+strconv.Itoa(conta.ID), "error")
			return nil, false
		}

		cobrancaExistente.TransactionID = boletoResponse.ID
		cobrancaExistente.UrlGerarCobranca = boletoResponse.BankSlipURL
		db.Save(&cobrancaExistente)

		return &cobrancaExistente, false
	}

	boletoResponse, errAPI, err := client.PaymentBoleto("", paymentBoleto)
	if err != nil {
		logger.Send(err.Error(), "error")
		return nil, false
	}
	if errAPI != nil {
		var errorsAsaas []string
		for _, errorAsaas := range errAPI.Errors {
			errorsAsaas = append(errorsAsaas, errorAsaas.Description)
		}
		logger.Send("["+strconv.Itoa(conta.EmpresaID)+"] Erro na requisição ASAAS da conta ["+strconv.Itoa(conta.ID)+"]: "+strings.Join(errorsAsaas, ","), "error")
		return nil, false
	}
	if boletoResponse == nil {
		logger.Send("Response do cadastro do boleto da conta "+strconv.Itoa(conta.ID), "error")
		return nil, false
	}

	cobranca := models.Cobranca{
		ContaID:           conta.ID,
		EmpresaID:         conta.EmpresaID,
		TransactionID:     boletoResponse.ID,
		UrlGerarCobranca:  boletoResponse.BankSlipURL,
		GeracaoAutomatica: true,
		Status:            true,
		Situacao:          "em_aberto",
	}

	db.Create(&cobranca)

	for _, split := range boletoResponse.Split {
		cobrancaTaxa := models.CobrancaTaxa{
			BoletoId:  cobranca.ID,
			EmpresaID: conta.EmpresaID,
			Tipo:      "imobia",
		}
		if split.FixedValue > 0 && split.PercentualValue > 0 {
			cobrancaTaxa.ValorReal = split.FixedValue
			cobrancaTaxa.ValorPercentual = split.PercentualValue

		} else if split.FixedValue > 0 && split.PercentualValue <= 0 {
			cobrancaTaxa.ValorReal = split.FixedValue
			cobrancaTaxa.ValorPercentual = 0

		} else {
			cobrancaTaxa.ValorReal = 0
			cobrancaTaxa.ValorPercentual = 0

		}

		db.Unscoped().Create(&cobrancaTaxa)
	}

	taxaAsaas := boletoResponse.Value - boletoResponse.NetValue

	if taxaAsaas != 0 {
		cobrancaTaxa := models.CobrancaTaxa{
			BoletoId:        cobranca.ID,
			EmpresaID:       conta.EmpresaID,
			Tipo:            "asaas",
			ValorReal:       taxaAsaas,
			ValorPercentual: 0,
		}

		db.Create(&cobrancaTaxa)
	}
	return &cobranca, true
}

func MakeSplit(cobranca sdk.PaymentBoleto, taxa models.EmpresaTaxa) (*sdk.PaymentBoleto, error) {

	if taxa.ValorReal <= 0 {
		return &cobranca, nil
	}

	taxaSplit, err := strconv.ParseFloat(os.Getenv("VALOR_TAXA"), 2)
	if err != nil {
		logger.Send(err.Error(), "error")
		return nil, err
	}

	valorFixo := taxa.ValorReal - float32(taxaSplit)

	if valorFixo > 0 {
		cobranca.Split = append(cobranca.Split, sdk.Split{
			WalletId:        os.Getenv("ASAAS_WALLET_ID"),
			FixedValue:      valorFixo,
			PercentualValue: 0,
		})
	} else if taxa.ValorPercentual > 0 {
		cobranca.Split = append(cobranca.Split, sdk.Split{
			WalletId:        os.Getenv("ASAAS_WALLET_ID"),
			FixedValue:      0,
			PercentualValue: taxa.ValorPercentual,
		})
	}

	return &cobranca, nil
}

func MakeFine(cobranca sdk.PaymentBoleto, multa models.Multa) (*sdk.PaymentBoleto, error) {

	if multa.Porcentagem <= 0 {
		return nil, errors.New("Porcentagem da multa menor igual a 0")
	}

	cobranca.Fine = sdk.PaymentFine{
		Value: multa.Porcentagem,
	}

	return &cobranca, nil
}

func MakeInterest(cobranca sdk.PaymentBoleto, juros models.Juros) (*sdk.PaymentBoleto, error) {
	if juros.Porcentagem <= 0 {
		return nil, errors.New("Porcentagem de juros menor igual a 0")
	}

	cobranca.Interest = sdk.PaymentInterest{
		Value: juros.Porcentagem,
	}

	return &cobranca, nil
}

func MakeDiscount(cobranca sdk.PaymentBoleto, Desconto models.Desconto) (*sdk.PaymentBoleto, error) {
	if Desconto.Valor <= 0 {
		return nil, errors.New("Porcentagem do Desconto menor igual a 0")
	}

	if Desconto.DescontarAteVencimento {
		cobranca.Discount = sdk.PaymentDiscount{
			Value:            Desconto.Valor,
			DueDateLimitDays: 0,
		}
	} else {
		cobranca.Discount = sdk.PaymentDiscount{
			Value:            Desconto.Valor,
			DueDateLimitDays: Desconto.PrazoAteVenciemnto,
		}
	}

	return &cobranca, nil
}
