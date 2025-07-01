package repositories

import (
	"api-gateway/internal/entities"

	"gorm.io/gorm"
)

type (
	CryptoChainRepository interface {
		GetCryptoChains() ([]entities.CryptoChain, error)
		GetCryptoChainById(id string) (*entities.CryptoChain, error)
		CreateCryptoChain(cryptoChain entities.CryptoChain) (*entities.CryptoChain, error)
		UpdateCryptoChain(id string, cryptoChain *entities.CryptoChain) (*entities.CryptoChain, error)
		DeleteCryptoChain(id string) error
	}

	cryptoChainRepository struct {
		db *gorm.DB
	}
)

func NewCryptoChainRepository(db *gorm.DB) CryptoChainRepository {
	return &cryptoChainRepository{
		db: db,
	}
}

func (r *cryptoChainRepository) GetCryptoChains() ([]entities.CryptoChain, error) {
	var res []entities.CryptoChain
	err := r.db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *cryptoChainRepository) GetCryptoChainById(id string) (*entities.CryptoChain, error) {
	var cryptoChain entities.CryptoChain
	err := r.db.First(&cryptoChain, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &cryptoChain, nil
}

func (r *cryptoChainRepository) CreateCryptoChain(cryptoChain entities.CryptoChain) (*entities.CryptoChain, error) {
	if err := r.db.Create(&cryptoChain).Error; err != nil {
		return nil, err
	}

	return &cryptoChain, nil
}

func (r *cryptoChainRepository) UpdateCryptoChain(id string, cryptoChain *entities.CryptoChain) (*entities.CryptoChain, error) {
	var existing entities.CryptoChain
	if err := r.db.
		First(&existing, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Model(&existing).
		Updates(cryptoChain).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (r *cryptoChainRepository) DeleteCryptoChain(id string) error {
	if err := r.db.Delete(&entities.CryptoChain{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
