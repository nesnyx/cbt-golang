package peserta

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Save(peserta Peserta) (Peserta, error)
	SavePeserta(peserta PesertaEntity) (PesertaEntity, error)
	SaveAdmin(peserta PesertaEntity) (CreateNewAdminResponse, error)
	SaveAllowedToken(token string, uuid string) (bool, error)
	SaveLastLogin(time int64, uuid string) (bool, error)
	FindAll(jenjang string) ([]GetAllPesertaInput, error)
	Update(uuid string, peserta Peserta) (Peserta, error)
	FindByUUID(uuid string) (GetPesertaTokenInput, error)
	UpdateByUUID(uuid string, peserta Peserta) (Peserta, error)
	FindByName(name string) (Peserta, error)
	FindByNoPendaftaran(noPendaftaran string) (Peserta, error)
	DeleteAllowedToken(uuid string) (bool, error)
	DeleteLastLogin(uuid string) (bool, error)
	DeletePeserta(uuid string) (bool, error)
	UpdateStartTimeByUUID(uuid string, time int64, isStarted bool, isFinished bool) (Peserta, error)
	ClearExpiredTokens(uuid string, expirationLimit int64) (bool, error)
	UpdatePeserta(uuid string, peserta UpdatePesertaInput) (UpdatePesertaInput, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ClearExpiredTokens(uuid string, expirationLimit int64) (bool, error) {
	// Query untuk mencari user yang last_login lebih lama dari 7200 detik (expired)
	query := "UPDATE users SET token = NULL, last_login = NULL WHERE uuid = ? AND last_login < ?"
	err := r.db.Exec(query, uuid, expirationLimit).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) SaveLastLogin(time int64, uuid string) (bool, error) {
	var peserta Peserta
	err := r.db.Model(&peserta).Where("uuid = ?", uuid).UpdateColumn("expired_time", time).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) SaveAllowedToken(token string, uuid string) (bool, error) {
	var peserta Peserta
	err := r.db.Model(&peserta).Where("uuid = ?", uuid).UpdateColumn("allowed_token", token).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) Save(peserta Peserta) (Peserta, error) {
	err := r.db.Table("peserta").Create(peserta).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) SavePeserta(peserta PesertaEntity) (PesertaEntity, error) {
	err := r.db.Table("peserta").Create(&peserta).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) SaveAdmin(peserta PesertaEntity) (CreateNewAdminResponse, error) {
	err := r.db.Table("peserta").Create(&peserta).Error
	if err != nil {
		return CreateNewAdminResponse{}, err
	}
	return CreateNewAdminResponse{}, nil
}

func (r *repository) FindAll(jenjang string) ([]GetAllPesertaInput, error) {
	var peserta []GetAllPesertaInput
	err := r.db.Table("peserta").Raw("SELECT * FROM peserta WHERE JSON_UNQUOTE(JSON_EXTRACT(name, '$.jenjang')) = ?", jenjang).Scan(&peserta).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) FindByUUID(uuid string) (GetPesertaTokenInput, error) {
	var peserta GetPesertaTokenInput
	fmt.Println(uuid)
	err := r.db.Table("peserta").Where("uuid = ?", uuid).First(&peserta).Error
	if err != nil {
		return GetPesertaTokenInput{}, err
	}
	return peserta, nil
}

func (r *repository) FindByName(name string) (Peserta, error) {
	var peserta Peserta
	err := r.db.Where("name = ?", name).First(&peserta).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) FindByNoPendaftaran(noPendaftaran string) (Peserta, error) {
	var peserta Peserta
	err := r.db.Where("no_pendaftaran = ?", noPendaftaran).First(&peserta).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) Update(uuid string, peserta Peserta) (Peserta, error) {
	err := r.db.Raw("UPDATE peserta SET answers = ? WHERE uuid = ?", peserta.Answers, uuid).Scan(&peserta).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) UpdatePeserta(uuid string, peserta UpdatePesertaInput) (UpdatePesertaInput, error) {

	updatePesertaQuery := `UPDATE peserta SET name = JSON_SET(name, '$.nama', ?, '$.institusi', ?, '$.no_telp', ?, '$.email', ?), no_pendaftaran = ? WHERE uuid = ?`
	err := r.db.Exec(updatePesertaQuery, peserta.Name, peserta.Institusi, peserta.NoTelp, peserta.Email, peserta.NoPendaftaran, uuid).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) UpdateByUUID(uuid string, peserta Peserta) (Peserta, error) {
	err := r.db.Model(&peserta).Where("uuid = ?", uuid).UpdateColumn("answers", peserta.Answers).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) UpdateStartTimeByUUID(uuid string, time int64, isStarted bool, isFinished bool) (Peserta, error) {
	var peserta Peserta
	err := r.db.Model(&peserta).Where("uuid = ?", uuid).UpdateColumn("start_time", time).UpdateColumn("is_started", isStarted).UpdateColumn("is_finish", isFinished).Error
	if err != nil {
		return peserta, err
	}
	return peserta, nil
}

func (r *repository) DeletePeserta(uuid string) (bool, error) {
	err := r.db.Table("peserta").Exec("DELETE FROM peserta WHERE uuid = ?", uuid).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) DeleteAllowedToken(uuid string) (bool, error) {
	var peserta Peserta
	err := r.db.Model(&peserta).Where("uuid = ?", uuid).UpdateColumn("allowed_token", nil).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) DeleteLastLogin(uuid string) (bool, error) {
	var peserta Peserta
	err := r.db.Model(&peserta).Where("uuid = ?", uuid).UpdateColumn("expired_time", nil).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
