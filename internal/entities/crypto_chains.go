package entities

import (
	"time"

	"github.com/google/uuid"
)

type CryptoChain struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	Name      string    `gorm:"type:varchar(255);unique;not null;column:name"`
	Symbol    string    `gorm:"type:varchar(255);unique;not null;column:symbol"`
	ImagePath *string   `gorm:"type:varchar(255);column:image_path"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`
}

func (CryptoChain) TableName() string {
	return "crypto_chains"
}

type CryptoChainResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	ImagePath *string   `json:"imagePath"`
	CreatedAt time.Time `json:"createdAt"`
}

type CryptoChainCreateRequest struct {
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	ImagePath *string `json:"imagePath,omitempty"`
}

type CryptoChainUpdateRequest struct {
	Name      *string `json:"name,omitempty"`
	Symbol    *string `json:"symbol,omitempty"`
	ImagePath *string `json:"imagePath,omitempty"`
}
