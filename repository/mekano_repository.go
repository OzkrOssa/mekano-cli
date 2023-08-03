package repository

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/OzkrOssa/mekano-cli/config"
	"github.com/mozillazg/go-unidecode"
	"github.com/xuri/excelize/v2"
)

type MekanoDataStruct struct {
	Tipo          string
	Prefijo       string
	Numero        string
	Secuencia     string
	Fecha         string
	Cuenta        string
	Terceros      string
	CentroCostos  string
	Nota          string
	Debito        string
	Credito       string
	Base          string
	Aplica        string
	TipoAnexo     string
	PrefijoAnexo  string
	NumeroAnexo   string
	Usuario       string
	Signo         string
	CuentaCobrar  string
	CuentaPagar   string
	NombreTercero string
	NombreCentro  string
	Interface     string
}

type mekanoInterface interface {
	Payment(file string) ([]MekanoDataStruct, error)
	Billing(file string, extras string) ([]MekanoDataStruct, error)
}

type mekanoRepository struct {
	dr DatabaseRepositoryInterface
}

func NewMekanoRepository(dr DatabaseRepositoryInterface) mekanoInterface {

	return &mekanoRepository{
		dr,
	}
}

func (mr *mekanoRepository) Payment(file string) ([]MekanoDataStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var paymentDataSlice []MekanoDataStruct
	var consecutive int = 0

	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	excelRows, err := xlsx.GetRows(xlsx.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	c, err := mr.dr.GetPayment(ctx)
	if err != nil {
		return nil, err
	}

	for i, row := range excelRows[1:] {

		consecutive = c.Consecutive + i + 1

		paymentData := MekanoDataStruct{
			Tipo:          "RC",
			Prefijo:       "_",
			Numero:        strconv.Itoa(consecutive + i),
			Secuencia:     "",
			Fecha:         row[4],
			Cuenta:        "13050501",
			Terceros:      row[1],
			CentroCostos:  "C1",
			Nota:          "RECAUDO POR VENTA SERVICIOS",
			Debito:        "0",
			Credito:       row[5],
			Base:          "0",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: row[2],
			NombreCentro:  "CENTRO DE COSTOS GENERAL",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		}
		paymentDataSlice = append(paymentDataSlice, paymentData)

		paymentData2 := MekanoDataStruct{
			Tipo:          "RC",
			Prefijo:       "_",
			Numero:        strconv.Itoa(consecutive),
			Secuencia:     "",
			Fecha:         row[4],
			Cuenta:        config.Cashier[row[9]],
			Terceros:      row[1],
			CentroCostos:  "C1",
			Nota:          "RECAUDO POR VENTA SERVICIOS",
			Debito:        row[5],
			Credito:       "0",
			Base:          "0",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: row[2],
			NombreCentro:  "CENTRO DE COSTOS GENERAL",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		}
		paymentDataSlice = append(paymentDataSlice, paymentData2)
	}
	exporterFile(paymentDataSlice)

	mr.dr.SavePayment(ctx, Payment{Consecutive: consecutive, CreateAt: time.Now().Format("2006-01-02"), FileName: file})

	PaymentStatistics(file, paymentDataSlice, c.Consecutive, consecutive)
	return paymentDataSlice, nil
}

func (mr *mekanoRepository) Billing(file string, extras string) ([]MekanoDataStruct, error) {

	var montoBaseFinal float64
	var montoIvaFinal float64
	var montoDebitoFinal float64
	var itemIvaBaseFinal float64

	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	billingFile, err := xlsx.GetRows(xlsx.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ivaXlsx, err := excelize.OpenFile(extras)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	itemsIvaFile, err := ivaXlsx.GetRows(ivaXlsx.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var BillingDataSheet []MekanoDataStruct

	if err != nil {
		log.Println(err, "itemsIvaFile")
	}

	for _, bRow := range billingFile[1:] {

		montoDebito, err := strconv.ParseFloat(bRow[14], 64)
		if err != nil {
			log.Println(err, "MontoDebito")
		}
		_, decimalDebito := math.Modf(montoDebito)
		if decimalDebito >= 0.5 {
			montoDebitoFinal = math.Ceil(montoDebito)
		} else {
			montoDebitoFinal = math.Round(montoDebito)
		}

		montoBase, err := strconv.ParseFloat(bRow[12], 64)
		if err != nil {
			log.Println(err, "MontoBase")
		}
		_, decimalBase := math.Modf(montoBase)
		if decimalBase >= 0.5 {
			montoBaseFinal = math.Ceil(montoBase)
		} else {
			montoBaseFinal = math.Round(montoBase)
		}

		montoIva, err := strconv.ParseFloat(strings.TrimSpace(bRow[13]), 64)
		if err != nil {
			log.Println(err, "MontoIva")
		}
		_, decimalIva := math.Modf(montoIva)

		if decimalIva >= 0.5 {
			montoIvaFinal = math.Ceil(montoIva)
		} else {
			montoIvaFinal = math.Round(montoIva)
		}

		if !strings.Contains(bRow[21], ",") {
			_, ok := config.Accounts[bRow[21]]
			if !ok {
				log.Println("Cuenta no existe en la base de datos: ", bRow[21])
			}
			billingNormal := MekanoDataStruct{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Secuencia:     "",
				Fecha:         bRow[9],
				Cuenta:        config.Accounts[bRow[21]],
				Terceros:      bRow[1],
				CentroCostos:  config.CostCenter[unidecode.Unidecode(bRow[17])],
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        "0",
				Credito:       fmt.Sprintf("%f", montoBaseFinal),
				Base:          "0",
				Aplica:        "",
				TipoAnexo:     "",
				PrefijoAnexo:  "",
				NumeroAnexo:   "",
				Usuario:       "SUPERVISOR",
				Signo:         "",
				CuentaCobrar:  "",
				CuentaPagar:   "",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
				Interface:     time.Now().Format("02/01/2006 15:04"),
			}

			BillingDataSheet = append(BillingDataSheet, billingNormal)

			billingIva := MekanoDataStruct{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Secuencia:     "",
				Fecha:         bRow[9],
				Cuenta:        "24080505",
				Terceros:      bRow[1],
				CentroCostos:  config.CostCenter[unidecode.Unidecode(bRow[17])],
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        "0",
				Credito:       fmt.Sprintf("%f", montoIvaFinal),
				Base:          fmt.Sprintf("%f", montoBaseFinal),
				Aplica:        "",
				TipoAnexo:     "",
				PrefijoAnexo:  "",
				NumeroAnexo:   "",
				Usuario:       "SUPERVISOR",
				Signo:         "",
				CuentaCobrar:  "",
				CuentaPagar:   "",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
				Interface:     time.Now().Format("02/01/2006 15:04"),
			}

			BillingDataSheet = append(BillingDataSheet, billingIva)

			billingCxC := MekanoDataStruct{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Secuencia:     "",
				Fecha:         bRow[9],
				Cuenta:        "13050501",
				Terceros:      bRow[1],
				CentroCostos:  config.CostCenter[unidecode.Unidecode(bRow[17])],
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        fmt.Sprintf("%f", montoDebitoFinal),
				Credito:       "0",
				Base:          "0",
				Aplica:        "",
				TipoAnexo:     "",
				PrefijoAnexo:  "",
				NumeroAnexo:   "",
				Usuario:       "SUPERVISOR",
				Signo:         "",
				CuentaCobrar:  "",
				CuentaPagar:   "",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
				Interface:     time.Now().Format("02/01/2006 15:04"),
			}

			BillingDataSheet = append(BillingDataSheet, billingCxC)
		} else {
			splitBillingItems := strings.Split(bRow[21], ",")
			for _, item := range splitBillingItems {
				for _, itemIva := range itemsIvaFile[1:] {

					if itemIva[1] == strings.TrimSpace(item) && itemIva[0] == bRow[0] {
						itemIvaBase, _ := strconv.ParseFloat(itemIva[2], 64)
						_, decimalIvaBase := math.Modf(itemIvaBase)

						if decimalIvaBase >= 0.5 {
							itemIvaBaseFinal = math.Ceil(itemIvaBase)
						} else {
							itemIvaBaseFinal = math.Round(itemIvaBase)
						}
						_, ok := config.Accounts[unidecode.Unidecode(strings.TrimSpace(item))]

						if !ok {
							log.Println("Cuenta no existe en la base de datos: ", unidecode.Unidecode(strings.TrimSpace(item)))

						}

						billingNormalPlus := MekanoDataStruct{
							Tipo:          "FVE",
							Prefijo:       "_",
							Numero:        bRow[8],
							Secuencia:     "",
							Fecha:         bRow[9],
							Cuenta:        config.Accounts[unidecode.Unidecode(strings.TrimSpace(item))],
							Terceros:      bRow[1],
							CentroCostos:  config.CostCenter[unidecode.Unidecode(bRow[17])],
							Nota:          "FACTURA ELECTRÓNICA DE VENTA",
							Debito:        "0",
							Credito:       fmt.Sprintf("%f", itemIvaBaseFinal),
							Base:          "0",
							Aplica:        "",
							TipoAnexo:     "",
							PrefijoAnexo:  "",
							NumeroAnexo:   "",
							Usuario:       "SUPERVISOR",
							Signo:         "",
							CuentaCobrar:  "",
							CuentaPagar:   "",
							NombreTercero: bRow[2],
							NombreCentro:  bRow[17],
							Interface:     time.Now().Format("02/01/2006 15:04"),
						}
						BillingDataSheet = append(BillingDataSheet, billingNormalPlus)
					}
				}
			}
			billingIvaPlus := MekanoDataStruct{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Secuencia:     "",
				Fecha:         bRow[9],
				Cuenta:        "24080505",
				Terceros:      bRow[1],
				CentroCostos:  config.CostCenter[unidecode.Unidecode(bRow[17])],
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        "0",
				Credito:       fmt.Sprintf("%f", montoIvaFinal),
				Base:          fmt.Sprintf("%f", montoBaseFinal),
				Aplica:        "",
				TipoAnexo:     "",
				PrefijoAnexo:  "",
				NumeroAnexo:   "",
				Usuario:       "SUPERVISOR",
				Signo:         "",
				CuentaCobrar:  "",
				CuentaPagar:   "",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
				Interface:     time.Now().Format("02/01/2006 15:04"),
			}

			BillingDataSheet = append(BillingDataSheet, billingIvaPlus)

			billingCxCPlus := MekanoDataStruct{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Secuencia:     "",
				Fecha:         bRow[9],
				Cuenta:        "13050501",
				Terceros:      bRow[1],
				CentroCostos:  config.CostCenter[unidecode.Unidecode(bRow[17])],
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        fmt.Sprintf("%f", montoDebitoFinal),
				Credito:       "0",
				Base:          "0",
				Aplica:        "",
				TipoAnexo:     "",
				PrefijoAnexo:  "",
				NumeroAnexo:   "",
				Usuario:       "SUPERVISOR",
				Signo:         "",
				CuentaCobrar:  "",
				CuentaPagar:   "",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
				Interface:     time.Now().Format("02/01/2006 15:04"),
			}

			BillingDataSheet = append(BillingDataSheet, billingCxCPlus)
		}
	}

	exporterFile(BillingDataSheet)
	BillingStatistics(BillingDataSheet)
	return BillingDataSheet, nil
}

func exporterFile(mekanoData []MekanoDataStruct) {
	txtFile, err := os.Create(filepath.Join(config.MekanoExportPath, "CONTABLE.txt"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer txtFile.Close()

	writer := csv.NewWriter(txtFile)
	writer.Comma = ','

	for _, data := range mekanoData {
		row := []string{
			data.Tipo,
			data.Prefijo,
			data.Numero,
			data.Secuencia,
			data.Fecha,
			data.Cuenta,
			data.Terceros,
			data.CentroCostos,
			data.Nota,
			data.Debito,
			data.Credito,
			data.Base,
			data.Aplica,
			data.TipoAnexo,
			data.PrefijoAnexo,
			data.NumeroAnexo,
			data.Usuario,
			data.Signo,
			data.CuentaCobrar,
			data.CuentaPagar,
			data.NombreTercero,
			data.NombreCentro,
			data.Interface,
		}
		writer.Write(row)
	}
	writer.Flush()
}

type paymentStatistics struct {
	FileName    string `json:"archivo"`
	RangoRC     string `json:"rango-rc"`
	Bancolombia int    `json:"bancolombia"`
	Davivienda  int    `json:"davivienda"`
	Susuerte    int    `json:"susuerte"`
	PayU        int    `json:"payu"`
	Efectivo    int    `json:"efectivo"`
	Total       int    `json:"total"`
}
type billingStatistics struct {
	Debito  float64 `json:"debito"`
	Credito float64 `json:"credito"`
	Base    float64 `json:"base"`
}

func PaymentStatistics(fileName string, data []MekanoDataStruct, initialRC, lastRC int) {

	var efectivo, bancolombia, davivienda, susuerte, payU, total int = 0, 0, 0, 0, 0, 0

	for _, d := range data {
		debito, err := strconv.Atoi(d.Debito)
		total += debito
		if err != nil {
			log.Println(err)
		}
		switch d.Cuenta {
		case "11050501": //Efectivo
			efectivo += debito
		case "11200501": //Bancolombia
			bancolombia += debito
		case "11200510": //Davivienda
			davivienda += debito
		case config.Cashier["SUSUERTE S"]: //Pay U
			susuerte += debito
		case config.Cashier["PAY U"]: //Susuerte
			payU += debito
		}
	}

	s := paymentStatistics{
		FileName:    fileName,
		RangoRC:     fmt.Sprintf("%d-%d", initialRC+1, lastRC),
		Efectivo:    efectivo,
		Bancolombia: bancolombia,
		Davivienda:  davivienda,
		PayU:        payU,
		Susuerte:    susuerte,
		Total:       total,
	}

	result, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(result))

}

var (
	d, c, b float64 = 0, 0, 0
)

func BillingStatistics(data []MekanoDataStruct) {

	for _, row := range data {
		debito, _ := strconv.ParseFloat(row.Debito, 64)
		d += debito
		credito, _ := strconv.ParseFloat(row.Credito, 64)
		c += credito
		base, _ := strconv.ParseFloat(row.Base, 64)
		b += base
	}

	bs := billingStatistics{
		Debito:  d,
		Credito: c,
		Base:    b,
	}

	result, err := json.Marshal(bs)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(result))

}
