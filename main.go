package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/OzkrOssa/mekano-cli/repository"
)

type arguments struct {
	paymentFile string
	billingFile string
	extrasFile  string
}

func main() {

	var args arguments
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	d, err := repository.NewDatabaseRepository(dns)
	if err != nil {
		log.Println(err)
	}

	mekano := repository.NewMekanoRepository(d)

	// Definir los flags
	flag.StringVar(&args.paymentFile, "p", "", "Ruta del archivo de pagos")
	flag.StringVar(&args.billingFile, "b", "", "Ruta del archivo de facturación")
	flag.StringVar(&args.extrasFile, "e", "", "Ruta del archivo de extras (opcional)")

	// Parsear los flags
	flag.Parse()

	// Verificar que se haya especificado una de las opciones (-p o -b)
	if args.paymentFile == "" && args.billingFile == "" {
		fmt.Println("Debes especificar al menos una opción (-p o -b)")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Procesar la opción de pagos (-p)
	if args.paymentFile != "" {
		mekano.Payment(args.paymentFile)
	}

	// Procesar la opción de facturación (-b)
	if args.billingFile != "" {
		if args.extrasFile != "" {
			mekano.Billing(args.billingFile, args.extrasFile)
		} else {
			fmt.Println("Debes especificar el parametro (-e)")
		}
	}
}
