package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/postgresql"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	dbClient   postgresql.Client
	sb         squirrel.StatementBuilderType
	orderCache domain.IOrderCache
}

func NewOrderRepository(dbClient postgresql.Client, orderCache domain.IOrderCache) *OrderRepository {
	return &OrderRepository{
		dbClient:   dbClient,
		sb:         squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		orderCache: orderCache,
	}
}

func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	// Serialize nested structures to JSON
	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return err
	}

	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return err
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	query, args, err := r.sb.Insert("orders").
		Columns(
			"order_uid",
			"track_number",
			"entry",
			"delivery",
			"payment",
			"items",
			"locale",
			"internal_signature",
			"customer_id",
			"delivery_service",
			"shardkey",
			"sm_id",
			"date_created",
			"oof_shard",
		).
		Values(
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			deliveryJSON,
			paymentJSON,
			itemsJSON,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.Shardkey,
			order.SMID,
			order.DateCreated,
			order.OOFShard,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbClient.Exec(ctx, query, args...)
	// Save to cache if no error
	if err != nil {
		r.orderCache.Save(ctx, order)
	}
	return err
}

func (r *OrderRepository) Get(ctx context.Context, orderID string) (*domain.Order, error) {
	// Try cache first
	cachedOrder, err := r.orderCache.Get(ctx, orderID)
	if err == nil {
		return cachedOrder, nil
	}

	query, args, err := r.sb.Select(
		"order_uid",
		"track_number",
		"entry",
		"delivery",
		"payment",
		"items",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard",
	).
		From("orders").
		Where(squirrel.Eq{"order_uid": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var order domain.Order
	var deliveryJSON, paymentJSON, itemsJSON []byte
	var dateCreated time.Time

	err = r.dbClient.QueryRow(ctx, query, args...).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&deliveryJSON,
		&paymentJSON,
		&itemsJSON,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SMID,
		&dateCreated,
		&order.OOFShard,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // or return custom "not found" error
		}
		return nil, err
	}

	// Parse JSON fields
	order.DateCreated = dateCreated
	if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	query, args, err := r.sb.Select(
		"order_uid",
		"track_number",
		"entry",
		"delivery",
		"payment",
		"items",
		"locale",
		"internal_signature",
		"customer_id",
		"delivery_service",
		"shardkey",
		"sm_id",
		"date_created",
		"oof_shard",
	).
		From("orders").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbClient.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte
		var dateCreated time.Time

		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&deliveryJSON,
			&paymentJSON,
			&itemsJSON,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SMID,
			&dateCreated,
			&order.OOFShard,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		order.DateCreated = dateCreated
		if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
