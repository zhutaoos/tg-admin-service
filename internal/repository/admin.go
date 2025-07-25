package repository

import (
	"app/internal/model"

	"gorm.io/gorm"
)

// AdminRepo 管理员数据访问接口
type AdminRepo interface {
	GetAdmin(admin *model.Admin) (*model.Admin, error)
	GetByAccount(account string) (*model.Admin, error)
	GetByID(id uint) (*model.Admin, error)
	Create(admin *model.Admin) error
	Update(admin *model.Admin) error
	UpdatePassword(id uint, hashedPassword string) error
	Delete(id int64) error
	DeleteBatch(ids []int) error
	GetList(groupId uint) ([]*model.Admin, error)
}

// AdminRepository 管理员数据访问实现
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓储实例
func NewAdminRepository(db *gorm.DB) AdminRepo {
	return &AdminRepository{db: db}
}

// GetAdmin 根据条件查询管理员信息
func (ar *AdminRepository) GetAdmin(admin *model.Admin) (*model.Admin, error) {
	err := ar.db.Where(admin).First(admin).Error
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// GetByAccount 根据账号查询管理员
func (ar *AdminRepository) GetByAccount(account string) (*model.Admin, error) {
	var admin model.Admin
	err := ar.db.Where("account = ?", account).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByID 根据ID查询管理员
func (ar *AdminRepository) GetByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	err := ar.db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// Create 创建管理员
func (ar *AdminRepository) Create(admin *model.Admin) error {
	return ar.db.Create(admin).Error
}

// Update 更新管理员信息
func (ar *AdminRepository) Update(admin *model.Admin) error {
	return ar.db.Model(admin).Updates(admin).Error
}

// UpdatePassword 更新管理员密码
func (ar *AdminRepository) UpdatePassword(id uint, hashedPassword string) error {
	return ar.db.Model(&model.Admin{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// Delete 删除管理员
func (ar *AdminRepository) Delete(id int64) error {
	return ar.db.Delete(&model.Admin{}, id).Error
}

// DeleteBatch 批量删除管理员
func (ar *AdminRepository) DeleteBatch(ids []int) error {
	return ar.db.Delete(&model.Admin{}, ids).Error
}

// GetList 获取管理员列表
func (ar *AdminRepository) GetList(groupId uint) ([]*model.Admin, error) {
	var admins []*model.Admin
	query := ar.db.Model(&model.Admin{})

	if groupId > 0 {
		query = query.Where("group_id = ?", groupId)
	}

	err := query.Find(&admins).Error
	return admins, err
}
