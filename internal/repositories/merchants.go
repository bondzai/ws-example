package repositories

import (
	"api-gateway/internal/entities"

	"gorm.io/gorm"
)

type (
	MerchantRepository interface {
		CreateMerchant(merchant entities.Merchant) (*entities.Merchant, error)
		GetMerchants() ([]entities.Merchant, error)
		GetMerchantsByIds(ids []string) ([]entities.Merchant, error)
		GetMerchantById(id string) (*entities.Merchant, error)
		UpdateMerchant(merchant *entities.Merchant) (*entities.Merchant, error)
		DeleteMerchant(id string) error
	}

	merchantRepository struct {
		db *gorm.DB
	}
)

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

func (r *merchantRepository) CreateMerchant(merchant entities.Merchant) (*entities.Merchant, error) {
	err := r.db.Create(&merchant).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) GetMerchants() ([]entities.Merchant, error) {
	var res []entities.Merchant
	err := r.db.
		Preload("CryptoChains").
		Preload("Category").
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *merchantRepository) GetMerchantsByIds(ids []string) ([]entities.Merchant, error) {
	var merchants []entities.Merchant
	err := r.db.
		Preload("CryptoChains").
		Preload("Category").
		Where("id IN ?", ids).
		Find(&merchants).Error
	if err != nil {
		return nil, err
	}
	return merchants, nil
}

func (r *merchantRepository) GetMerchantById(id string) (*entities.Merchant, error) {
	var merchant entities.Merchant
	err := r.db.
		Preload("CryptoChains").
		Preload("Category").
		First(&merchant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *merchantRepository) UpdateMerchant(merchant *entities.Merchant) (*entities.Merchant, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Save the updated merchant (scalar fields and nested structs)
		if err := tx.Save(merchant).Error; err != nil {
			return err
		}

		// Automatically detect and update slice associations using the helper
		if err := UpdateAssociations(tx, merchant); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return merchant, nil
}

func (r *merchantRepository) DeleteMerchant(id string) error {
	if err := r.db.
		Delete(&entities.Merchant{}, "id = ?", id).
		Error; err != nil {
		return err
	}

	return nil
}
