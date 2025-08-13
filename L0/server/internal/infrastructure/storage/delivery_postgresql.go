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

type DeliveryRepository struct {
	dbClient postgresql.Client
	sb       squirrel.StatementBuilderType
}

func NewDeliveryRepository(dbClient postgresql.Client) *DeliveryRepository {
	return &DeliveryRepository{
		dbClient: dbClient,
		sb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *DeliveryRepository) Save(ctx context.Context, delivery *domain.Delivery) (int, error) {
	const mark = "storage.DeliveryRepository.Save"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Insert("deliveries").
		Columns("name", "phone", "zip", "city", "address", "region", "email").
		Values(
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return 0, fmt.Errorf("build query: %w", err)
	}

	logger.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	var id int
	err = r.dbClient.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		logger.Error("Failed to execute query", slogext.Err(err))
		return 0, fmt.Errorf("execute query: %w", err)
	}

	return id, nil
}

func (r *DeliveryRepository) Get(ctx context.Context, id int) (*domain.Delivery, error) {
	const mark = "storage.DeliveryRepository.Get"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Select("name", "phone", "zip", "city", "address", "region", "email").
		From("deliveries").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return nil, fmt.Errorf("build query: %w", err)
	}

	logger.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	var delivery domain.Delivery
	err = r.dbClient.QueryRow(ctx, query, args...).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("Delivery not found", slog.Int("delivery_id", id))
			return nil, nil
		}
		logger.Error("Failed to execute query", slogext.Err(err))
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &delivery, nil
}
