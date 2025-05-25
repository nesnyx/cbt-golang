package soal

import (
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Save(soal Soal) (Soal, error)
	FindAll() ([]Soal, error)
	FindByKodeSoal(kodeSoal string) (Soal, error)
	Delete(id int) (bool, error)
	SaveTmpSoal(soal TmpSoal) (TmpSoal, error)
	FindAllTmpAnswer() ([]TmpSoal, error)
	FindByQuestionId(questionId string) (TmpSoal, error)
	Update(id uint, questionText string, questionPicture string, correctAnswer string, tmpSoal TmpSoal) (bool, error)
	FindByID(id uint) (TmpSoal, error)
	FindCorrectAnswerByQuestionID(questionId string) (*TmpSoal, error)
	FindAllTmpSoalByType(typeSoal string) ([]TmpSoal, error)
	FindAllTmpSoalByTypeStudent(typeSoal string) ([]TmpSoal, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAllTmpSoalByTypeStudent(typeSoal string) ([]TmpSoal, error) {
	var tmpSoal []TmpSoal
	err := r.db.Table("tmp_soals").Raw("SELECT * FROM tmp_soals WHERE JSON_UNQUOTE(JSON_EXTRACT(name, '$.jenjang')) = ?", typeSoal).Scan(&tmpSoal).Error
	if err != nil {
		return tmpSoal, err
	}
	return tmpSoal, nil
}

func (r *Repository) FindAllTmpSoalByType(typeSoal string) ([]TmpSoal, error) {
	var tmpSoal []TmpSoal
	err := r.db.Raw("SELECT * FROM tmp_soals WHERE type = ?", typeSoal).Scan(&tmpSoal).Error
	if err != nil {
		return tmpSoal, err
	}
	return tmpSoal, nil
}

func (r *Repository) FindCorrectAnswerByQuestionID(questionId string) (*TmpSoal, error) {
	var tmpSoal *TmpSoal
	err := r.db.Raw("SELECT correct_answer FROM tmp_soals WHERE question_id = ?", questionId).Scan(&tmpSoal).Error
	if err != nil {
		return tmpSoal, err
	}
	return tmpSoal, nil
}

func (r *Repository) Update(id uint, questionText string, questionPicture string, correctAnswer string, tmpSoal TmpSoal) (bool, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := tx.Table("tmp_soals").Raw("UPDATE tmp_soals SET options = ? WHERE id = ?", tmpSoal.Options, id).Scan(&tmpSoal).Error
	if err != nil {
		return false, err
	}
	updateQuestionQuery := `UPDATE tmp_soals SET pertanyaan = JSON_SET(pertanyaan,'$.text',?,'$.gambar',?), correct_answer = ? WHERE id = ?`
	if err := tx.Exec(updateQuestionQuery, questionText, questionPicture, correctAnswer, id).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error

}

func (r *Repository) FindByQuestionId(questionId string) (TmpSoal, error) {
	var tmpSoal TmpSoal
	err := r.db.Where("question_id = ?", questionId).First(&tmpSoal).Error
	if err != nil {
		return tmpSoal, err
	}
	return tmpSoal, nil
}

func (r *Repository) FindAllTmpAnswer() ([]TmpSoal, error) {
	var soal []TmpSoal
	err := r.db.Table("tmp_soals").Find(&soal).Error
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (r *Repository) SaveTmpSoal(soal TmpSoal) (TmpSoal, error) {
	err := r.db.Table("tmp_soals").Create(&soal).Error
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (r *Repository) Save(soal Soal) (Soal, error) {
	err := r.db.Create(&soal).Error
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (r *Repository) FindAll() ([]Soal, error) {
	var soal []Soal
	err := r.db.Find(&soal).Error
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (r *Repository) FindByKodeSoal(kodeSoal string) (Soal, error) {
	var soal Soal
	err := r.db.Where("kode_soal = ?", kodeSoal).First(&soal).Error
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (r *Repository) FindByID(id uint) (TmpSoal, error) {
	var tmpSoal TmpSoal
	err := r.db.Table("tmp_soals").Where("id = ?", id).First(&tmpSoal).Error
	if err != nil {
		return tmpSoal, err

	}
	return tmpSoal, nil
}

func (r *Repository) Delete(id int) (bool, error) {
	var soal TmpSoal
	err := r.db.Table("tmp_soals").Delete(&soal, "id = ?", id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
