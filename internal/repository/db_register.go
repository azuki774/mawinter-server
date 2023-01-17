package repository

import (
	"errors"
	"mawinter-server/internal/model"
	"mawinter-server/internal/register"

	"gorm.io/gorm"
)

// InsertUniqueCatIDRecord は 同一のカテゴリIDがない場合ときに挿入、既にあればエラーを返す
func (d *DBRepository) InsertUniqueCatIDRecord(req model.Recordstruct) (res model.Recordstruct, err error) {
	err = d.Conn.Table(getRecordTable(req.Datetime)).
		Where("category_id = ?", req.CategoryID).Error

	if err == nil {
		// already recorded
		return model.Recordstruct{}, register.ErrAlreadyRegisted
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// unknown error
		return model.Recordstruct{}, err
	}
	dbres := d.Conn.Table(getRecordTable(req.Datetime)).Create(&req)
	if dbres.Error != nil {
		return model.Recordstruct{}, dbres.Error
	}

	res = req
	return res, err
}
