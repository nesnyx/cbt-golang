package result

import "gorm.io/gorm"

type RepositoryInterface interface {
	Save(results ResultsEntity) (ResultsEntity, error)
	FindAll(jenjang string) ([]Results, error)
	Delete(id int) (bool, error)
	Update(id int, lulus bool) (ResultsEntity, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Delete(id int) (bool, error) {
	err := r.db.Exec("DELETE FROM results_backup_2 WHERE id = ?", id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) FindAll(jenjang string) ([]Results, error) {
	var results []Results
	// err := r.db.Table("results").Find(&results).Error
	err := r.db.Raw("SELECT r.*,CAST(r.score AS DECIMAL(10,2)) as score,JSON_EXTRACT(p.name, '$.nama') as name,JSON_EXTRACT(p.name, '$.email') as email,JSON_EXTRACT(p.name, '$.no_telp') as no_telp  FROM results_backup_2 r INNER JOIN peserta p ON p.uuid = r.uuid WHERE p.role = ?", "peserta|"+jenjang).Scan(&results).Error
	if err != nil {
		return results, err
	}

	return results, nil
}

func (r *Repository) Save(results ResultsEntity) (ResultsEntity, error) {
	err := r.db.Table("results_backup_2").Save(&results).Error
	if err != nil {
		return ResultsEntity{}, err
	}
	return ResultsEntity{}, nil
}

func (r *Repository) Update(id int, lulus bool) (ResultsEntity, error) {
	err := r.db.Exec("UPDATE results_backup_2 SET lulus = ? WHERE id = ?", lulus, id).Error
	if err != nil {
		return ResultsEntity{}, err
	}
	return ResultsEntity{}, err
}
