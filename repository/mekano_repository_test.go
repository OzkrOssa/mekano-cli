package repository

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/OzkrOssa/mekano-cli/config"
)

func TestMekanoPayment(t *testing.T) {
	file := "../../test_files/payment_test.xlsx"
	dr, err := NewDatabaseRepository("root:root@tcp(localhost:3306)/mekano_test")
	if err != nil {
		t.Fatalf("Error al inicial la base de datos: %v", err)
	}

	mekano := NewMekanoRepository(dr)
	paymentData, err := mekano.Payment(file)
	if err != nil {
		if err != nil {
			t.Fatalf("Error al procesar los archivos de pagos: %v", err)
		}
	}

	d, err := NewDatabaseRepository("root:root@tcp(localhost:3306)/mekano_test")

	if err != nil {
		if err != nil {
			t.Fatalf("Error initializing the database: %v", err)
		}
	}

	c, err := d.GetPayment(context.Background())
	if err != nil {
		if err != nil {
			t.Fatalf("Error to get payment: %v", err)
		}
	}

	expectedData := []MekanoDataStruct{
		{
			Tipo:          "RC",
			Prefijo:       "_",
			Numero:        strconv.Itoa(c.Consecutive),
			Secuencia:     "",
			Fecha:         "01/07/2023",
			Cuenta:        "13050501",
			Terceros:      "1060536367",
			CentroCostos:  "C1",
			Nota:          "RECAUDO POR VENTA SERVICIOS",
			Debito:        "0",
			Credito:       "75000",
			Base:          "0",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: "XIOMARA DURANGO GOEZ",
			NombreCentro:  "CENTRO DE COSTOS GENERAL",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		{
			Tipo:          "RC",
			Prefijo:       "_",
			Numero:        strconv.Itoa(c.Consecutive),
			Secuencia:     "",
			Fecha:         "01/07/2023",
			Cuenta:        "13452505",
			Terceros:      "1060536367",
			CentroCostos:  "C1",
			Nota:          "RECAUDO POR VENTA SERVICIOS",
			Debito:        "75000",
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
			NombreTercero: "XIOMARA DURANGO GOEZ",
			NombreCentro:  "CENTRO DE COSTOS GENERAL",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
	}

	filePath := filepath.Join(config.MekanoExportPath, "CONTABLE.txt")

	_, err = os.Stat(filePath)

	if err != nil {
		t.Fatalf("No se pudo encontrar el archivo CONTABLE.txt: %v", err)
	}

	if !reflect.DeepEqual(paymentData, expectedData) {
		t.Errorf("La salida no coincide con el resultado esperado.")
	}
}

func TestMekanoBilling(t *testing.T) {
	file := "../../test_files/billing_test.xlsx"
	extras := "../../test_files/extras.xlsx"

	dr, err := NewDatabaseRepository("root:root@tcp(localhost:3306)/mekano_test")

	if err != nil {
		t.Fatalf("Error al iniciar la base de datos: %v", err)
	}

	mekano := NewMekanoRepository(dr)

	billingData, err := mekano.Billing(file, extras)
	if err != nil {
		t.Fatalf("Error al procesar los archivos de facturacion: %v", err)
	}

	filePath := filepath.Join(config.MekanoExportPath, "CONTABLE.txt")

	_, err = os.Stat(filePath)

	if err != nil {
		t.Fatalf("No se pudo encontrar el archivo CONTABLE.txt: %v", err)
	}

	expectedData := []MekanoDataStruct{
		//normal
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66137",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "41457070",
			Terceros:      "159122542",
			CentroCostos:  "101",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "0",
			Credito:       "63950.000000",
			Base:          "0",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: "GERMAN DE JESUS ESCOBAR LOAIZA",
			NombreCentro:  "RIOSUCIO",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		// Factura 2
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66137",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "41459030",
			Terceros:      "159122542",
			CentroCostos:  "101",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "0",
			Credito:       "23025.000000",
			Base:          "0",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: "GERMAN DE JESUS ESCOBAR LOAIZA",
			NombreCentro:  "RIOSUCIO",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66137",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "24080505",
			Terceros:      "159122542",
			CentroCostos:  "101",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "0",
			Credito:       "16525.000000",
			Base:          "86975.000000",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: "GERMAN DE JESUS ESCOBAR LOAIZA",
			NombreCentro:  "RIOSUCIO",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66137",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "13050501",
			Terceros:      "159122542",
			CentroCostos:  "101",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "103500.000000",
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
			NombreTercero: "GERMAN DE JESUS ESCOBAR LOAIZA",
			NombreCentro:  "RIOSUCIO",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		/////////////////////////////////////////////////////////////////
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66138",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "41457070",
			Terceros:      "797339211",
			CentroCostos:  "102",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "0",
			Credito:       "63025.000000",
			Base:          "0",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: "JOSE DAVID PARRA SILVA",
			NombreCentro:  "SUPIA",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66138",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "24080505",
			Terceros:      "797339211",
			CentroCostos:  "102",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "0",
			Credito:       "11975.000000",
			Base:          "63025.000000",
			Aplica:        "",
			TipoAnexo:     "",
			PrefijoAnexo:  "",
			NumeroAnexo:   "",
			Usuario:       "SUPERVISOR",
			Signo:         "",
			CuentaCobrar:  "",
			CuentaPagar:   "",
			NombreTercero: "JOSE DAVID PARRA SILVA",
			NombreCentro:  "SUPIA",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
		{
			Tipo:          "FVE",
			Prefijo:       "_",
			Numero:        "66138",
			Secuencia:     "",
			Fecha:         "27/06/2023",
			Cuenta:        "13050501",
			Terceros:      "797339211",
			CentroCostos:  "102",
			Nota:          "FACTURA ELECTRÓNICA DE VENTA",
			Debito:        "75000.000000",
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
			NombreTercero: "JOSE DAVID PARRA SILVA",
			NombreCentro:  "SUPIA",
			Interface:     time.Now().Format("02/01/2006 15:04"),
		},
	}

	if !reflect.DeepEqual(billingData, expectedData) {
		t.Errorf("La salida no coincide con el resultado esperado.")
	}
}

func TestMekanoExporterFile(t *testing.T) {
	mekanoData := []MekanoDataStruct{
		{Tipo: "Tipo1", Prefijo: "Prefijo1", Numero: "Num1", Secuencia: "Sec1", Fecha: "Fecha1", Cuenta: "Cuenta1", Terceros: "Tercero1", CentroCostos: "CentroCostos1", Nota: "Nota1", Debito: "Debito1", Credito: "Credito1", Base: "Base1", Aplica: "Aplica1", TipoAnexo: "TipoAnexo1", PrefijoAnexo: "PrefijoAnexo1", NumeroAnexo: "NumeroAnexo1", Usuario: "Usuario1", Signo: "Signo1", CuentaCobrar: "CuentaCobrar1", CuentaPagar: "CuentaPagar1", NombreTercero: "NombreTercero1", NombreCentro: "NombreCentro1", Interface: "Interface1"},
	}

	// Directorio temporal para el archivo de prueba

	// Ejecutar la función de prueba
	exporterFile(mekanoData)

	// Comprobar si el archivo ha sido creado
	filePath := filepath.Join(config.MekanoExportPath, "CONTABLE.txt")
	_, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("No se pudo encontrar el archivo CONTABLE.txt: %v", err)
	}

	// Leer el contenido del archivo generado
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("No se pudo abrir el archivo CONTABLE.txt: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Error al leer el contenido del archivo CONTABLE.txt: %v", err)
	}

	// Comprobar si el número de filas coincide con los datos de prueba
	if len(rows) != len(mekanoData) {
		t.Fatalf("El número de filas en el archivo no coincide con los datos de prueba. Esperado: %d, Obtenido: %d", len(mekanoData), len(rows))
	}

	// Comprobar si los valores de las filas coinciden con los datos de prueba
	for i, row := range rows {
		expectedRow := []string{
			mekanoData[i].Tipo,
			mekanoData[i].Prefijo,
			mekanoData[i].Numero,
			mekanoData[i].Secuencia,
			mekanoData[i].Fecha,
			mekanoData[i].Cuenta,
			mekanoData[i].Terceros,
			mekanoData[i].CentroCostos,
			mekanoData[i].Nota,
			mekanoData[i].Debito,
			mekanoData[i].Credito,
			mekanoData[i].Base,
			mekanoData[i].Aplica,
			mekanoData[i].TipoAnexo,
			mekanoData[i].PrefijoAnexo,
			mekanoData[i].NumeroAnexo,
			mekanoData[i].Usuario,
			mekanoData[i].Signo,
			mekanoData[i].CuentaCobrar,
			mekanoData[i].CuentaPagar,
			mekanoData[i].NombreTercero,
			mekanoData[i].NombreCentro,
			mekanoData[i].Interface,
		}

		// Convertir los datos a cadena para hacer la comparación más sencilla
		expectedRowStr := strings.Join(expectedRow, ",")
		actualRowStr := strings.Join(row, ",")

		if expectedRowStr != actualRowStr {
			t.Errorf("Los valores en la fila %d no coinciden con los datos de prueba. Esperado: %s, Obtenido: %s", i, expectedRowStr, actualRowStr)
		}
	}
}

func TestPaymentStatistics(t *testing.T) {

}

func TestBillingStatistics(t *testing.T) {

}
