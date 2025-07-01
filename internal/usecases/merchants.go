package usecases

import (
	"api-gateway/internal/entities"
	"api-gateway/internal/repositories"
	"api-gateway/pkg/errs"
	"api-gateway/pkg/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

type (
	MerchantUseCase interface {
		GetMerchants() ([]entities.MerchantResponse, error)
		GetMerchantById(id string) (*entities.MerchantResponse, error)
		CreateMerchant(req entities.MerchantCreateRequest) (*entities.MerchantResponse, error)
		UpdateMerchant(id string, req entities.MerchantUpdateRequest) (*entities.MerchantResponse, error)
		DeleteMerchant(id string) error
	}

	merchantUseCase struct {
		repo      repositories.MerchantRepository
		cacheRepo repositories.MerchantCacheRepository
	}
)

func NewMerchantUseCase(
	repo repositories.MerchantRepository,
	cacheRepo repositories.MerchantCacheRepository,
) MerchantUseCase {
	return &merchantUseCase{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

func (u *merchantUseCase) GetMerchants() ([]entities.MerchantResponse, error) {
	merchants, err := u.repo.GetMerchants()
	if err != nil {
		log.Println("failed to get merchants: ", err)
		return nil, errs.HandleError(err)
	}

	var merchantResponses []entities.MerchantResponse
	for _, m := range merchants {
		merchantResponses = append(merchantResponses, convertToMerchantResponse(m))
	}

	return merchantResponses, nil
}

func (u *merchantUseCase) GetMerchantById(id string) (*entities.MerchantResponse, error) {
	merchant, err := u.repo.GetMerchantById(id)
	if err != nil {
		log.Println("failed to get merchant by id: ", merchant)
		return nil, errs.HandleError(err)
	}

	response := convertToMerchantResponse(*merchant)
	return &response, nil
}

func (u *merchantUseCase) CreateMerchant(req entities.MerchantCreateRequest) (*entities.MerchantResponse, error) {
	merchant := entities.Merchant{
		Id:         uuid.New(),
		Name:       req.Name,
		CategoryId: uuid.MustParse(req.CategoryId),
		Location:   req.Location,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		CreatedAt:  time.Now().UTC(),
	}

	createdMerchant, err := u.repo.CreateMerchant(merchant)
	if err != nil {
		log.Println("failed to create merchant: ", err)
		return nil, errs.HandleError(err)
	}

	response := convertToMerchantResponse(*createdMerchant)
	return &response, nil
}

func (u *merchantUseCase) UpdateMerchant(id string, req entities.MerchantUpdateRequest) (*entities.MerchantResponse, error) {
	merchant, err := u.repo.GetMerchantById(id)
	if err != nil || merchant == nil {
		log.Println("failed to get merchant: ", err)
		return nil, errs.HandleError(err)
	}

	if err := utils.ApplyUpdates(merchant, &req); err != nil {
		log.Println("failed to apply update merchant:", err)
		return nil, errs.HandleError(err)
	}

	updatedMerchant, err := u.repo.UpdateMerchant(merchant)
	if err != nil {
		log.Println("failed to update merchant: ", err)
		return nil, errs.HandleError(err)
	}

	response := convertToMerchantResponse(*updatedMerchant)
	return &response, nil
}

func (u *merchantUseCase) DeleteMerchant(id string) error {
	if err := u.repo.DeleteMerchant(id); err != nil {
		log.Println("failed to delete merchant: ", err)
		return errs.HandleError(err)
	}

	return nil
}

func convertToMerchantResponse(m entities.Merchant) entities.MerchantResponse {
	cryptoChains := []entities.CryptoChainResponse{}
	for _, c := range m.CryptoChains {
		cryptoChains = append(cryptoChains, entities.CryptoChainResponse{
			Id:        c.Id,
			Name:      c.Name,
			Symbol:    c.Symbol,
			ImagePath: c.ImagePath,
			CreatedAt: c.CreatedAt,
		})
	}

	return entities.MerchantResponse{
		Id:   m.Id,
		Name: m.Name,
		Category: entities.MerchantCategoryResponse{
			Id:    m.Category.Id,
			Name:  m.Category.Name,
			Icon:  m.Category.Icon,
			Color: m.Category.Color,
		},
		Location:     m.Location,
		Latitude:     m.Latitude,
		Longitude:    m.Longitude,
		CreatedAt:    m.CreatedAt,
		CryptoChains: cryptoChains,
	}
}
