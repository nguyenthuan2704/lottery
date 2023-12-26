package storage

import (
	"Lottery/modules/model"
	"context"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *model.PlayerCreation) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
