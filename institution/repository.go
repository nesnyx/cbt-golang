package institution

import "gorm.io/gorm"

type RepositoryInterface interface {
	Save(institution Institution) (Institution, error)
	Update(id int, institution Institution) (InstitutionResponse, error)
	FindByID(id int) (InstitutionResponse, error)
	Delete(id int) (bool, error)
	FindAll() ([]Institution, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindAll() ([]Institution, error) {
	var institutions []Institution
	err := r.db.Find(&institutions).Error
	if err != nil {
		return institutions, err
	}
	return institutions, nil
}

func (r *Repository) Save(institution Institution) (Institution, error) {
	err := r.db.Table("institutions").Create(&institution).Error
	if err != nil {
		return institution, err
	}
	return institution, nil
}

func (r *Repository) Update(id int, institution Institution) (InstitutionResponse, error) {
	err := r.db.Model(&institution).Where("id = ?", id).UpdateColumn("name", institution.Name).UpdateColumn("address", institution.Address).Error
	if err != nil {
		return InstitutionResponse{}, err
	}
	return InstitutionResponse{institution.Name, institution.Address}, nil
}

func (r *Repository) FindByID(id int) (InstitutionResponse, error) {
	var institution Institution
	err := r.db.Table("institutions").Where("id= ?", id).First(&institution).Error
	if err != nil {
		return InstitutionResponse{}, err
	}
	return InstitutionResponse{institution.Name, institution.Address}, nil
}

func (r *Repository) Delete(id int) (bool, error) {
	var institution Institution
	err := r.db.Delete(&institution, "id = ?", id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
