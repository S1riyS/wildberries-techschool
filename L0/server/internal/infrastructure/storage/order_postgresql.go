package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/postgresql"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	dbClient     postgresql.Client
	sb           squirrel.StatementBuilderType
	orderCache   domain.IOrderCache
	deliveryRepo domain.IDeliveryRepository
	paymentRepo  domain.IPaymentRepository
	itemRepo     domain.IItemRepository
}

func NewOrderRepository(
	dbClient postgresql.Client,
	orderCache domain.IOrderCache,
	deliveryRepo domain.IDeliveryRepository,
	paymentRepo domain.IPaymentRepository,
	itemRepo domain.IItemRepository,
) *OrderRepository {
	return &OrderRepository{
		dbClient:     dbClient,
		sb:           squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		orderCache:   orderCache,
		deliveryRepo: deliveryRepo,
		paymentRepo:  paymentRepo,
		itemRepo:     itemRepo,
	}
}

func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	const mark = "storage.OrderRepository.Save"
	logger := slog.With(slog.String("mark", mark))

	// Start transaction
	tx, err := r.dbClient.Begin(ctx)
	if err != nil {
		logger.Error("Failed to begin transaction", slogext.Err(err))
		return fmt.Errorf("begin transaction: %w", err)
	}
	// Ensure the transaction will be rolled back if not committed
	defer tx.Rollback(ctx)

	// Set current time for order
	order.DateCreated = time.Now()

	// Save delivery
	deliveryID, err := r.deliveryRepo.Save(ctx, &order.Delivery)
	if err != nil {
		logger.Error("Failed to save delivery", slogext.Err(err))
		return fmt.Errorf("save delivery: %w", err)
	}

	// Save payment
	if err := r.paymentRepo.Save(ctx, &order.Payment); err != nil {
		logger.Error("Failed to save payment", slogext.Err(err))
		return fmt.Errorf("save payment: %w", err)
	}

	// Save items
	for _, item := range order.Items {
		if err := r.itemRepo.Save(ctx, item); err != nil {
			logger.Error("Failed to save item", slogext.Err(err))
			return fmt.Errorf("save item: %w", err)
		}
	}

	// Save order
	query, args, err := r.sb.Insert("orders").
		Columns(
			"order_uid",
			"track_number",
			"entry",
			"delivery_id",
			"payment_transaction",
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
			deliveryID,
			order.Payment.Transaction,
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
		logger.Error("Failed to build query", slogext.Err(err))
		return fmt.Errorf("build query: %w", err)
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		logger.Error("Failed to execute query", slogext.Err(err))
		return fmt.Errorf("execute query: %w", err)
	}

	// Save order items
	for _, item := range order.Items {
		query, args, err := r.sb.Insert("order_items").
			Columns("order_uid", "item_chrt_id").
			Values(order.OrderUID, item.ChrtID).
			ToSql()
		if err != nil {
			logger.Error("Failed to build query", slogext.Err(err))
			return fmt.Errorf("build query: %w", err)
		}

		if _, err = tx.Exec(ctx, query, args...); err != nil {
			logger.Error("Failed to execute query", slogext.Err(err))
			return fmt.Errorf("execute query: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		logger.Error("Failed to commit transaction", slogext.Err(err))
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, orderID string) (*domain.Order, error) {
	const mark = "storage.OrderRepository.Get"
	logger := slog.With(slog.String("mark", mark))

	// Try cache first
	cachedOrder, err := r.orderCache.Get(ctx, orderID)
	if err == nil {
		logger.Debug("Order found in cache", slog.String("order_uid", orderID))
		return cachedOrder, nil
	}

	// Get order from DB
	query, args, err := r.sb.Select(
		"order_uid",
		"track_number",
		"entry",
		"delivery_id",
		"payment_transaction",
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
		logger.Error("Failed to build query", slogext.Err(err))
		return nil, fmt.Errorf("build query: %w", err)
	}

	var order domain.Order
	var deliveryID int
	var paymentTransaction string
	var dateCreated time.Time

	err = r.dbClient.QueryRow(ctx, query, args...).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&deliveryID,
		&paymentTransaction,
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
			logger.Debug("Order not found in database", slog.String("order_uid", orderID))
			return nil, nil
		}
		logger.Error("Failed to execute query", slogext.Err(err))
		return nil, fmt.Errorf("execute query: %w", err)
	}

	order.DateCreated = dateCreated

	// Get delivery
	delivery, err := r.deliveryRepo.Get(ctx, deliveryID)
	if err != nil {
		logger.Error("Failed to get delivery", slogext.Err(err))
		return nil, fmt.Errorf("get delivery: %w", err)
	}
	order.Delivery = *delivery

	// Get payment
	payment, err := r.paymentRepo.Get(ctx, paymentTransaction)
	if err != nil {
		logger.Error("Failed to get payment", slogext.Err(err))
		return nil, fmt.Errorf("get payment: %w", err)
	}
	order.Payment = *payment

	// Get items
	items, err := r.itemRepo.GetByOrder(ctx, orderID)
	if err != nil {
		logger.Error("Failed to get items", slogext.Err(err))
		return nil, fmt.Errorf("get items: %w", err)
	}
	order.Items = items

	// Save to cache
	if err := r.orderCache.Save(ctx, &order); err != nil {
		logger.Error("Failed to set order in cache", slogext.Err(err))
	}

	return &order, nil
}

func (r *OrderRepository) GetRecentlyCreated(ctx context.Context, limit int) ([]*domain.Order, error) {
	const mark = "storage.OrderRepository.GetRecentlyCreated"
	logger := slog.With(slog.String("mark", mark))

	query, args, err := r.sb.Select(
		"order_uid",
		"track_number",
		"entry",
		"delivery_id",
		"payment_transaction",
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
		OrderBy("date_created DESC").
		Limit(uint64(limit)).
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

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var deliveryID int
		var paymentTransaction string
		var dateCreated time.Time

		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&deliveryID,
			&paymentTransaction,
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
			logger.Error("Failed to scan row", slogext.Err(err), slog.String("order_uid", order.OrderUID))
			return nil, fmt.Errorf("scan row: %w", err)
		}

		order.DateCreated = dateCreated

		// Get delivery
		delivery, err := r.deliveryRepo.Get(ctx, deliveryID)
		if err != nil {
			logger.Error("Failed to get delivery", slogext.Err(err))
			return nil, fmt.Errorf("get delivery: %w", err)
		}
		order.Delivery = *delivery

		// Get payment
		payment, err := r.paymentRepo.Get(ctx, paymentTransaction)
		if err != nil {
			logger.Error("Failed to get payment", slogext.Err(err))
			return nil, fmt.Errorf("get payment: %w", err)
		}
		order.Payment = *payment

		// Get items
		items, err := r.itemRepo.GetByOrder(ctx, order.OrderUID)
		if err != nil {
			logger.Error("Failed to get items", slogext.Err(err))
			return nil, fmt.Errorf("get items: %w", err)
		}
		order.Items = items

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Failed to get rows", slogext.Err(err))
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}
