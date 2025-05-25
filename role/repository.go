package role

import "gorm.io/gorm"

type RepositoryInterface interface {
	FindByID(id int) (Roles, error)
	SaveHasRole(roleId int, uuid string) bool
	//FindHasRole(uuid string) (Roles, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) SaveHasRole(roleId int, uuid string) bool {
	hasRole := HasRole{}
	hasRole.RoleID = roleId
	hasRole.UUID = uuid
	err := r.db.Table("hasRole").Create(hasRole).Error
	if err != nil {
		return false
	}
	return true
}

//func (r *Repository) FindHasRole(uuid string) (Roles, error) {
//	var roles Roles
//	err := r.db.Table("hasRole").Where("uuid = ?", uuid).Error
//	if err != nil {
//		return Roles{}, err
//	}
//	return
//}

func (r *Repository) FindByID(id int) (Roles, error) {
	var role Roles
	err := r.db.Table("roles").First(&role, "id = ?", id).Error
	if err != nil {
		return role, err
	}
	return role, nil
}
