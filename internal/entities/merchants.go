package entities

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	Id         uuid.UUID  `gorm:"type:uuid;primary_key;column:id"`
	CategoryId uuid.UUID  `gorm:"type:uuid;not null;column:category_id"`
	Name       string     `gorm:"type:varchar(255);not null;column:name"`
	Location   string     `gorm:"type:varchar(255);not null;column:location"`
	Latitude   float64    `gorm:"type:decimal(11,8);column:latitude"`
	Longitude  float64    `gorm:"type:decimal(11,8);column:longitude"`
	VerifiedAt *time.Time `gorm:"type:timestamp;column:verified_at"`
	CreatedAt  time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`

	CryptoChains []CryptoChain    `gorm:"many2many:merchant_crypto_chain_map;foreignKey:Id;joinForeignKey:MerchantID;References:Id;joinReferences:ChainId"`
	Category     MerchantCategory `gorm:"foreignKey:CategoryId;references:Id"`
}

func (Merchant) TableName() string {
	return "merchants"
}

type MerchantCategory struct {
	Id    uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	Name  string    `gorm:"type:varchar(255);not null;column:name"`
	Icon  *string   `gorm:"column:icon"`
	Color *string   `gorm:"column:color"`
}

func (MerchantCategory) TableName() string {
	return "merchant_categories"
}

type MerchantCryptoChainMap struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	MerchantId    uuid.UUID `gorm:"type:uuid;not null;index;column:merchant_id"`
	ChainId       uuid.UUID `gorm:"type:uuid;not null;index;column:chain_id"`
	WalletAddress *string   `gorm:"type:varchar(255);unique;not null;column:wallet_address"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`

	Merchant    Merchant    `gorm:"foreignKey:MerchantId;references:Id"`
	CryptoChain CryptoChain `gorm:"foreignKey:ChainId;references:Id"`
}

func (MerchantCryptoChainMap) TableName() string {
	return "merchant_crypto_chain_map"
}

type MerchantCreateRequest struct {
	Name       string  `json:"name" validate:"required"`
	CategoryId string  `json:"categoryId" validate:"required"`
	Location   string  `json:"location" validate:"required"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type MerchantUpdateRequest struct {
	Name         *string                             `json:"name,omitempty"`
	Category     *string                             `json:"category,omitempty"`
	Location     *string                             `json:"location,omitempty"`
	Latitude     *float64                            `json:"latitude,omitempty"`
	Longitude    *float64                            `json:"longitude,omitempty"`
	CryptoChains *[]MerchantCryptoChainUpdateRequest `json:"cryptoChains,omitempty"`
}

type MerchantCryptoChainUpdateRequest struct {
	ChainId       uuid.UUID `json:"chainId" validate:"required"`
	WalletAddress *string   `json:"walletAddress"`
}

type MerchantResponse struct {
	Id           uuid.UUID                `json:"id"`
	Name         string                   `json:"name"`
	Location     string                   `json:"location"`
	Latitude     float64                  `json:"latitude"`
	Longitude    float64                  `json:"longitude"`
	CreatedAt    time.Time                `json:"createdAt"`
	Category     MerchantCategoryResponse `json:"category"`
	CryptoChains []CryptoChainResponse    `json:"cryptoChains"`
}

type MerchantCryptoChainResponse struct {
	Id            uuid.UUID           `json:"id"`
	Merchant      MerchantResponse    `json:"merchant"`
	CryptoChain   CryptoChainResponse `json:"cryptoChain"`
	WalletAddress string              `json:"walletAddress"`
	CreatedAt     time.Time           `json:"createdAt"`
}

type MerchantCategoryCreateRequest struct {
	Name string  `json:"name" validate:"required"`
	Icon *string `json:"icon,omitempty"`
}

type MerchantCategoryUpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Icon *string `json:"icon,omitempty"`
}

type MerchantCategoryResponse struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Icon  *string   `json:"icon"`
	Color *string   `json:"color"`
}
