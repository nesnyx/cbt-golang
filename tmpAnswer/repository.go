package tmpAnswer

import "gorm.io/gorm"

type RepositoryInterface interface {
	Save(tmpAnswer TmpAnswer) (TmpAnswer, error)
	FindAll() ([]TmpAnswer, error)
	FindByQuestionID(questionID string, uuid string) (TmpAnswer, error)
	FindByUUID(uuid string) (TmpAnswer, error)
	Update(questionID string, uuid string, answer string, time int64) (TmpAnswer, error)
	FindAnswersByUUID(uuid string) ([]TmpAnswer, error)
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) FindAnswersByUUID(uuid string) ([]TmpAnswer, error) {
	var tmpAnswer []TmpAnswer
	err := r.db.Raw("SELECT answers, peserta_uuid,questions_id FROM tmp_answer  WHERE peserta_uuid = ?", uuid).Scan(&tmpAnswer).Error
	if err != nil {
		return tmpAnswer, err
	}
	return tmpAnswer, nil
}

func (r *Repository) FindByUUID(uuid string) (TmpAnswer, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindAll() ([]TmpAnswer, error) {
	var tmpAnswer []TmpAnswer
	err := r.db.Table("tmp_answer").Find(&tmpAnswer).Error
	if err != nil {
		return tmpAnswer, err
	}
	return tmpAnswer, nil
}

func (r *Repository) Save(tmpAnswer TmpAnswer) (TmpAnswer, error) {
	err := r.db.Table("tmp_answer").Create(&tmpAnswer).Error
	if err != nil {
		return tmpAnswer, err
	}
	return tmpAnswer, nil
}

func (r *Repository) FindByQuestionID(questionID string, uuid string) (TmpAnswer, error) {
	var tmpAnswer TmpAnswer
	err := r.db.Table("tmp_answer").Where("questions_id = ?", questionID).Where("peserta_uuid = ?", uuid).First(&tmpAnswer).Error
	if err != nil {
		return tmpAnswer, err
	}
	return tmpAnswer, nil
}

func (r *Repository) Update(questionID string, uuid string, answer string, time int64) (TmpAnswer, error) {
	var tmpAnswer TmpAnswer
	err := r.db.Table("tmp_answer").Raw("UPDATE tmp_answer SET answers = ?, time = ? WHERE questions_id = ? AND peserta_uuid = ?", answer, time, questionID, uuid).Scan(&tmpAnswer).Error
	if err != nil {
		return tmpAnswer, err
	}
	return tmpAnswer, nil
}
