package storage

import (
	"Lottery/modules/model"
	"context"
)

func (s *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.Player, error) {
	var data model.Player

	if err := s.db.Where(cond).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
