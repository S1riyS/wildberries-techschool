package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/postgresql"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository struct {
	dbClient postgresql.Client
	sb       squirrel.StatementBuilderType
}

func NewPaymentRepository(dbClient postgresql.Client) *PaymentRepository {
	return &PaymentRepository{
		dbClient: dbClient,
		sb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	const mark = "storage.PaymentRepository.Save"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Insert("payments").
		Columns(
			"transaction",
			"request_id",
			"currency",
			"provider",
			"amount",
			"payment_dt",
			"bank",
			"delivery_cost",
			"goods_total",
			"custom_fee",
		).
		Values(
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDT,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
		).
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return fmt.Errorf("build query: %w", err)
	}

	logger.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	_, err = r.dbClient.Exec(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to execute query", slogext.Err(err))
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *PaymentRepository) Get(ctx context.Context, transaction string) (*domain.Payment, error) {
	const mark = "storage.PaymentRepository.Get"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Select(
		"transaction",
		"request_id",
		"currency",
		"provider",
		"amount",
		"payment_dt",
		"bank",
		"delivery_cost",
		"goods_total",
		"custom_fee",
	).
		From("payments").
		Where(squirrel.Eq{"transaction": transaction}).
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return nil, fmt.Errorf("build query: %w", err)
	}

	logger.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	var payment domain.Payment
	err = r.dbClient.QueryRow(ctx, query, args...).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("Payment not found", slog.String("transaction", transaction))
			return nil, nil
		}
		logger.Error("Failed to execute query", slogext.Err(err))
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &payment, nil
}
