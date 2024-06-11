package storage

type Storage struct {
}

func NewStorage(dsn string) (*Storage, error) {
	return &Storage{}, nil
}
