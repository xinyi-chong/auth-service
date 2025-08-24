package user

import (
	"auth-service/pkg/filters"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *Filter) ([]User, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
	UsernameOrEmailExists(ctx context.Context, username *string, email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) filteredQuery(ctx context.Context, filter *Filter) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&User{})

	if filter != nil {
		if filter.Email != nil {
			query = query.Where("email = ?", *filter.Email)
		}

		if filter.Username != nil {
			query = query.Where("username = ?", *filter.Username)
		}

		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}

		query = filters.PaginateQuery(query, &filter.Pagination, []string{"username", "email", "is_active"})
	}

	return query
}

func (r *repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).First(&user, id)
	return &user, result.Error
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user)
	return &user, result.Error
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, user *User) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) List(ctx context.Context, filter *Filter) ([]User, error) {
	var users []User
	result := r.filteredQuery(ctx, filter).Find(&users)
	return users, result.Error
}

func (r *repository) Count(ctx context.Context, filter *Filter) (int64, error) {
	var count int64
	result := r.filteredQuery(ctx, filter).Count(&count)
	return count, result.Error
}

func (r *repository) UsernameOrEmailExists(ctx context.Context, username *string, email string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&User{}).
		Or("email = ?", email)

	if username != nil && *username != "" {
		query = query.Or("username = ?", *username)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
