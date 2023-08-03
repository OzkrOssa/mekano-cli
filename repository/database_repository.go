package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Payment struct {
	Consecutive int
	CreateAt    string
	FileName    string
}

type Billing struct {
	Debit    int
	Credit   int
	Base     int
	FileName string
	CreateAt string
}

type DatabaseRepositoryInterface interface {
	GetPayment(ctx context.Context) (Payment, error)
	SavePayment(ctx context.Context, payment Payment) error
	SaveBilling(ctx context.Context, billing Billing) error
}

type DatabaseRepository struct {
	db *sql.DB
}

func NewDatabaseRepository(dns string) (DatabaseRepositoryInterface, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &DatabaseRepository{
		db,
	}, nil
}

func (r *DatabaseRepository) GetPayment(ctx context.Context) (Payment, error) {
	query := "SELECT consecutive, create_at, file_name FROM mekanopayments ORDER BY id DESC LIMIT 1;"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return Payment{}, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return Payment{}, err
	}
	defer rows.Close()

	// Lista para almacenar los pagos obtenidos
	var payments []Payment

	// Itera sobre los resultados
	for rows.Next() {
		var payment Payment

		// Escanea los valores de la fila en la estructura Payment
		if err := rows.Scan(&payment.Consecutive, &payment.CreateAt, &payment.FileName); err != nil {
			return Payment{}, err
		}

		// Agrega el pago a la lista
		payments = append(payments, payment)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		return Payment{}, err
	}

	return payments[0], nil
}

func (r *DatabaseRepository) SavePayment(ctx context.Context, payment Payment) error {
	insertSQL := "INSERT INTO mekanopayments (consecutive, create_at, file_name) VALUES (?, ?, ?)"
	stmt, err := r.db.PrepareContext(ctx, insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Ejecuta la consulta con los valores de la estructura Payment
	_, err = stmt.ExecContext(ctx, payment.Consecutive, payment.CreateAt, payment.FileName)
	if err != nil {
		return err
	}
	return nil
}

func (r *DatabaseRepository) SaveBilling(ctx context.Context, billing Billing) error {
	// Implementa la lógica para guardar datos de facturación en la base de datos
	insertSQL := "INSERT INTO mekanobilling (debit, credit, base, file_name, create_at) VALUES (?,?,?,?,?)"

	stmt, err := r.db.PrepareContext(ctx, insertSQL)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, billing.Debit, billing.Credit, billing.Base, billing.FileName, billing.CreateAt)
	if err != nil {
		return err
	}
	return nil
}
