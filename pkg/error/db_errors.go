package dberrors

import (
	"errors"
	"github.com/xinyi-chong/common-lib/consts"
	apperrors "github.com/xinyi-chong/common-lib/errors"

	"gorm.io/gorm"
)

func WrapDBError(err error, field consts.Field) *apperrors.Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrXNotFound.WithField(field).Wrap(err)
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperrors.ErrXConflict.WithField(field).Wrap(err)
	}

	return apperrors.ErrInternalServerError.Wrap(err)
}
