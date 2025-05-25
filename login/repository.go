package login

import "gorm.io/gorm"

type RepositoryInterface interface {
	Login(user Login) (Login, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Login(user Login) (Login, error) {
	err := r.db.Raw("SELECT * FROM login WHERE name = ?", user.Name).Scan(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
