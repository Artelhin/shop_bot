package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"shop_bot/models"
	"time"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %s", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) GetUserByID(id int64) (*models.User, error) {
	rows, err := s.db.Query(`select * from users where id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from users: %s", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Username,
			&user.ChatID,
			&user.AccessHash)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &(users[0]), nil
}

func (s *Storage) CreateUser(user *models.User) error {
	now := sql.NullTime{Time: time.Now(), Valid: true}
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := s.db.Exec(`
		insert into users (
        	id, chat_id, username, access_hash, created_at, updated_at
        ) values (
            $1, $2, $3, $4, $5, $6
        );`, user.ID, user.ChatID, user.Username, user.AccessHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("can't insert user: %s", err)
	}
	return nil
}

func (s *Storage) UpdateUser(user *models.User) error {
	user.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	_, err := s.db.Exec(`
		update users set
			chat_id = $1,
			username = $2,
			access_hash = $3,
			updated_at = $4
		where id = $5;
		`, user.ChatID, user.Username, user.AccessHash, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("can't update user: %s", err)
	}
	return nil
}

func (s *Storage) GetTopLevelCategories() ([]models.Category, error) {
	rows, err := s.db.Query(`select * from categories where parent_id is null and deleted_at is null`)
	if err != nil {
		return nil, fmt.Errorf("can't select from categories: %s", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
			&category.ParentID,
			&category.Name)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		categories = append(categories, category)
	}
	if len(categories) == 0 {
		return nil, nil
	}
	return categories, nil
}

func (s *Storage) GetSubcategoriesByCategoryID(id int64) ([]models.Category, error) {
	rows, err := s.db.Query(`select * from categories where parent_id = $1 and deleted_at is null`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from categories: %s", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
			&category.ParentID,
			&category.Name)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (s *Storage) GetCategoryByID(id int64) (*models.Category, error) {
	rows, err := s.db.Query(`select * from categories where id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from categories: %s", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		category := models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
			&category.ParentID,
			&category.Name)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		categories = append(categories, category)
	}
	if len(categories) == 0 {
		return nil, nil
	}
	return &(categories[0]), nil
}

func (s *Storage) GetItemsByCategoryID(id int64) ([]models.Item, error) {
	rows, err := s.db.Query(`select * from items where items.category_id = $1 and deleted_at is null`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from categories: %s", err)
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
			&item.Name,
			&item.Description,
			&item.CategoryID,
			&item.Image)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Storage) GetItemByID(id int64) (*models.Item, error) {
	rows, err := s.db.Query(`select * from items where id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from categories: %s", err)
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
			&item.Name,
			&item.Description,
			&item.CategoryID,
			&item.Image)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &(items[0]), nil
}

func (s *Storage) GetItemsInStoresByItemID(id int64) ([]models.ItemToStorage, error) {
	rows, err := s.db.Query(`select * from item_to_storage where item_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from item_to_storage: %s", err)
	}
	defer rows.Close()

	items := make([]models.ItemToStorage, 0)
	for rows.Next() {
		item := models.ItemToStorage{}
		err = rows.Scan(
			&item.ItemID,
			&item.StorageID,
			&item.Count)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items, nil
}

func (s *Storage) GetStoragesForItemID(id int64) ([]models.Storage, error) {
	rows, err := s.db.Query(`select s.* from storages s left join item_to_storage its on its.storage_id = s.id where its.item_id = $1 and its.count > 0`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from item_to_storage: %s", err)
	}
	defer rows.Close()

	storages := make([]models.Storage, 0)
	for rows.Next() {
		storage := models.Storage{}
		err = rows.Scan(
			&storage.ID,
			&storage.CreatedAt,
			&storage.UpdatedAt,
			&storage.DeletedAt,
			&storage.Name,
			&storage.Address)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		storages = append(storages, storage)
	}
	return storages, nil
}

func (s *Storage) GetStorageByID(id int64) (*models.Storage, error) {
	rows, err := s.db.Query(`select * from storages where id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't select from storages: %s", err)
	}
	defer rows.Close()

	storages := make([]models.Storage, 0)
	for rows.Next() {
		storage := models.Storage{}
		err = rows.Scan(
			&storage.ID,
			&storage.CreatedAt,
			&storage.UpdatedAt,
			&storage.DeletedAt,
			&storage.Name,
			&storage.Address)
		if err != nil {
			return nil, fmt.Errorf("can't convert from rows: %s", err)
		}
		storages = append(storages, storage)
	}
	if len(storages) == 0 {
		return nil, nil
	}
	return &(storages[0]), nil
}

const (
	OrderResultError = iota
	OrderResultNotInStock
	OrderResultSuccess
)

func (s *Storage) CreateOrder(ctx context.Context, order *models.Order) (int, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault})
	if err != nil {
		return OrderResultError, fmt.Errorf("can't create transaction")
	}

	rows, err := tx.Query(`select * from item_to_storage where item_id = $1 and storage_id = $2`, order.ItemID, order.StorageID)
	if err != nil {
		return OrderResultError, fmt.Errorf("can't check item to storage: %s", err)
	}
	if !rows.Next() {
		return OrderResultNotInStock, nil
	}
	its := &models.ItemToStorage{}
	err = rows.Scan(&its.ItemID, &its.StorageID, &its.Count)
	if err != nil {
		return OrderResultError, fmt.Errorf("can't scan from rows: %s", err)
	}
	rows.Close()
	if its.Count < 1 {
		return OrderResultNotInStock, nil
	}

	_, err = tx.Exec(`update item_to_storage set count = count - 1 where item_id = $1 and storage_id = $2`, order.ItemID, order.StorageID)
	if err != nil {
		return OrderResultError, fmt.Errorf("can't order from item to storate: %s", err)
	}

	_, err = tx.Exec(`
		insert into orders (
        	created_at, updated_at, deleted_at, user_id, item_id, storage_id, active, code
    	) values (now(), now(), null, $1, $2, $3, $4, $5)`,
		order.UserID, order.ItemID, order.StorageID, order.Active, order.Code)
	if err != nil {
		return OrderResultError, fmt.Errorf("can't create new order: %s", err)
	}
	if err = tx.Commit(); err != nil {
		return OrderResultError, fmt.Errorf("can't commit new order: %s", err)
	}
	return OrderResultSuccess, nil
}

func (s *Storage) DeactivateOrders(ctx context.Context) (int, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault})
	if err != nil {
		return 0, fmt.Errorf("can't create tx: %s", err)
	}

	rows, err := tx.Query(`select * from orders where age(created_at, now()) > interval '24' hour and active`)
	if err != nil {
		return 0, fmt.Errorf("can't select orders for deactivation: %s", err)
	}
	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
			&order.UserID,
			&order.ItemID,
			&order.StorageID,
			&order.Active,
			&order.Code)
		if err != nil {
			return 0, fmt.Errorf("can't scan from rows: %s", err)
		}
		orders = append(orders, order)
	}
	rows.Close()

	for i := range orders {
		_, err := tx.Exec(`update item_to_storage set count = count + 1 where item_id = $1 and storage_id = $2`,
			orders[i].ItemID, orders[i].StorageID)
		if err != nil {
			return 0, fmt.Errorf("can't add item back to stock: %s", err)
		}
		_, err = tx.Exec(`update orders set active = false where id = $1`, orders[i].ID)
		if err != nil {
			return 0, fmt.Errorf("can't deactivate order: %s", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("can't commit tx: %s", err)
	}
	return len(orders), nil
}
