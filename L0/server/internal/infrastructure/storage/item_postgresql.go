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

type ItemRepository struct {
	dbClient postgresql.Client
	sb       squirrel.StatementBuilderType
}

func NewItemRepository(dbClient postgresql.Client) *ItemRepository {
	return &ItemRepository{
		dbClient: dbClient,
		sb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ItemRepository) Save(ctx context.Context, item *domain.Item) error {
	const mark = "storage.ItemRepository.Save"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Insert("items").
		Columns(
			"chrt_id",
			"track_number",
			"price",
			"rid",
			"name",
			"sale",
			"size",
			"total_price",
			"nm_id",
			"brand",
			"status",
		).
		Values(
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		).
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.dbClient.Exec(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to execute query", slogext.Err(err))
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *ItemRepository) Get(ctx context.Context, chrtID int) (*domain.Item, error) {
	const mark = "storage.ItemRepository.Get"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Select(
		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status",
	).
		From("items").
		Where(squirrel.Eq{"chrt_id": chrtID}).
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return nil, fmt.Errorf("build query: %w", err)
	}

	var item domain.Item
	err = r.dbClient.QueryRow(ctx, query, args...).Scan(
		&item.ChrtID,
		&item.TrackNumber,
		&item.Price,
		&item.RID,
		&item.Name,
		&item.Sale,
		&item.Size,
		&item.TotalPrice,
		&item.NMID,
		&item.Brand,
		&item.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error("Failed to execute query", slogext.Err(err))
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) GetByOrder(ctx context.Context, orderUID string) ([]*domain.Item, error) {
	const mark = "storage.ItemRepository.GetByOrder"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Select(
		"i.chrt_id",
		"i.track_number",
		"i.price",
		"i.rid",
		"i.name",
		"i.sale",
		"i.size",
		"i.total_price",
		"i.nm_id",
		"i.brand",
		"i.status",
	).
		From("items i").
		Join("order_items oi ON i.chrt_id = oi.item_chrt_id").
		Where(squirrel.Eq{"oi.order_uid": orderUID}).
		ToSql()
	if err != nil {
		logger.Error("Failed to build query", slogext.Err(err))
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.dbClient.Query(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to execute query", slogext.Err(err))
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		var item domain.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			logger.Error("Failed to scan row", slogext.Err(err))
			return nil, fmt.Errorf("scan row: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Failed to get rows", slogext.Err(err))
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return items, nil
}
