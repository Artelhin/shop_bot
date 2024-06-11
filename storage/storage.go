package storage

import (
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
