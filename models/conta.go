package models

import (
	"go-cron-check-payments/database"
	"log"
	"strconv"
)

type Conta struct {
	ID                            int           `json:"id"`
	CodigoInterno                 interface{}   `json:"codigo_interno"`
	ModuloID                      interface{}   `json:"modulo_id"`
	ModuloType                    interface{}   `json:"modulo_type"`
	OrigemID                      interface{}   `json:"origem_id"`
	OrigemType                    interface{}   `json:"origem_type"`
	Valor                         float64       `json:"valor"`
	Saldo                         float64       `json:"saldo"`
	NumeroParcela                 int64         `json:"numero_parcela"`
	TotalParcelas                 int64         `json:"total_parcelas"`
	Fatura                        int64         `json:"fatura"`
	GrupoID                       int64         `json:"grupo_id"`
	Multa                         interface{}   `json:"multa"`
	Juro                          interface{}   `json:"juro"`
	Desconto                      interface{}   `json:"desconto"`
	Descricao                     string        `json:"descricao"`
	Tipo                          string        `json:"tipo"`
	DataVencimento                string        `json:"data_vencimento"`
	DataVencimentoOriginal        string        `json:"data_vencimento_original"`
	DataEmissao                   string        `json:"data_emissao"`
	Situacao                      string        `json:"situacao"`
	Cancelamento                  bool          `json:"cancelamento"`
	InicioReferencia              interface{}   `json:"inicio_referencia"`
	FinalReferencia               interface{}   `json:"final_referencia"`
	Caixa_id                      int64         `json:"caixa_id"`
	Cobranca                      interface{}   `json:"cobranca"`
	TaxaCobranca                  bool          `json:"taxa_cobranca"`
	NotaFiscalHabilitada          interface{}   `json:"nota_fiscal_habilitada"`
	TemRepasseAutomatico          bool          `json:"tem_repasse_automatico"`
	ProprietariosRepasse          bool          `json:"proprietarios_repasse"`
	FormaCobranca                 interface{}   `json:"forma_cobranca"`
	RepasseAutomatico             bool          `json:"repasse_automatico"`
	Retroativa                    bool          `json:"retroativa"`
	MotivoRepasseAutomatico       interface{}   `json:"motivo_repasse_automatico"`
	GeracaoAutomaticaNf           int64         `json:"geracao_automatica_nf"`
	MotivoGeracaoAutomaticaNf     interface{}   `json:"motivo_geracao_automatica_nf"`
	GeracaoAutomaticaBoleto       int64         `json:"geracao_automatica_boleto"`
	MotivoGeracaoAutomaticaBoleto interface{}   `json:"motivo_geracao_automatica_boleto"`
	Pagamentos                    []interface{} `json:"pagamentos"`
	PagamentoAgendado             bool          `json:"pagamentoAgendado"`
}

type AllContas struct {
	Data []Conta
}

func GetAllContas(selectQuery string, whereQuery string) (*AllContas, error) {
	var allContas *AllContas
	var conta Conta

	db := database.Connect()
	defer db.Close()

	rows, err := db.Query(selectQuery + "FROM contas" + whereQuery)

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&conta.ID, &conta.Situacao)
		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		} else {
			allContas.Data = append(allContas.Data, conta)
		}
	}

	return allContas, nil
}

func GetConta(id string) (*Conta, error) {
	var conta *Conta

	db := database.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id,situacao FROM contas WHERE id=" + id)

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&conta.ID, &conta.Situacao)
		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		}
	}

	return conta, nil
}

// Nao funciona ainda
func (c *Conta) UpdateConta(params string) (*Conta, error) {
	var conta *Conta

	db := database.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id,situacao FROM contas WHERE id=" + strconv.Itoa(c.ID))

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&conta.ID, &conta.Situacao)
		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		}
	}

	return conta, nil
}
