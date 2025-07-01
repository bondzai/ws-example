package repositories

import (
	"api-gateway/internal/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MerchantCategoryRepository interface {
		GetMerchantCategories() ([]entities.MerchantCategory, error)
		GetMerchantCategoryById(id string) (*entities.MerchantCategory, error)
		CreateMerchantCategory(merchantCategory entities.MerchantCategory) (*entities.MerchantCategory, error)
		UpdateMerchantCategory(id string, merchantCategory *entities.MerchantCategory) (*entities.MerchantCategory, error)
		DeleteMerchantCategory(id string) error
	}

	merchantCategoryRepository struct {
		db *gorm.DB
	}
)

func NewMerchantCategoryRepository(db *gorm.DB) MerchantCategoryRepository {
	return &merchantCategoryRepository{db: db}
}

func (r *merchantCategoryRepository) CreateMerchantCategory(merchantCategory entities.MerchantCategory) (*entities.MerchantCategory, error) {
	merchantCategory.Id = uuid.New()

	if err := r.db.Create(&merchantCategory).Error; err != nil {
		return nil, err
	}

	return &merchantCategory, nil
}

func (r *merchantCategoryRepository) GetMerchantCategories() ([]entities.MerchantCategory, error) {
	var res []entities.MerchantCategory
	err := r.db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *merchantCategoryRepository) GetMerchantCategoryById(id string) (*entities.MerchantCategory, error) {
	var merchantCategory entities.MerchantCategory
	err := r.db.First(&merchantCategory, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &merchantCategory, nil
}

func (r *merchantCategoryRepository) UpdateMerchantCategory(id string, merchantCategory *entities.MerchantCategory) (*entities.MerchantCategory, error) {
	var existing entities.MerchantCategory
	if err := r.db.
		First(&existing, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Model(&existing).
		Updates(merchantCategory).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (r *merchantCategoryRepository) DeleteMerchantCategory(id string) error {
	if err := r.db.Delete(&entities.MerchantCategory{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
