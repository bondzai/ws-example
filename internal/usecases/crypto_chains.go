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
	CryptoChainUseCase interface {
		CreateCryptoChain(req entities.CryptoChainCreateRequest) (*entities.CryptoChainResponse, error)
		GetCryptoChains() ([]entities.CryptoChainResponse, error)
		GetCryptoChainById(id string) (*entities.CryptoChainResponse, error)
		UpdateCryptoChain(id string, req *entities.CryptoChainUpdateRequest) (*entities.CryptoChainResponse, error)
	}

	cryptoChainUseCase struct {
		repo repositories.CryptoChainRepository
	}
)

func NewCryptoChainUseCase(repo repositories.CryptoChainRepository) CryptoChainUseCase {
	return &cryptoChainUseCase{
		repo: repo,
	}
}

func (u *cryptoChainUseCase) CreateCryptoChain(req entities.CryptoChainCreateRequest) (*entities.CryptoChainResponse, error) {
	if req.Name == "" || req.Symbol == "" {
		return nil, errs.NewBadRequestError("name and symbol are required")
	}

	createdChain, err := u.repo.CreateCryptoChain(
		entities.CryptoChain{
			Id:        uuid.New(),
			Name:      req.Name,
			Symbol:    req.Symbol,
			ImagePath: req.ImagePath,
			CreatedAt: time.Now(),
		},
	)
	if err != nil {
		log.Println("failed to create crypto chain: ", err)
		return nil, errs.HandleError(err)
	}

	response := convertToCryptoChainResponse(*createdChain)
	return &response, nil
}

func (u *cryptoChainUseCase) GetCryptoChains() ([]entities.CryptoChainResponse, error) {
	cryptoChains, err := u.repo.GetCryptoChains()
	if err != nil {
		log.Println("failed to get crypto chains: ", err)
		return nil, errs.HandleError(err)
	}

	var cryptoChainResponses []entities.CryptoChainResponse
	for _, m := range cryptoChains {
		cryptoChainResponses = append(cryptoChainResponses, convertToCryptoChainResponse(m))
	}

	return cryptoChainResponses, nil
}

func (u *cryptoChainUseCase) GetCryptoChainById(id string) (*entities.CryptoChainResponse, error) {
	cryptoChain, err := u.repo.GetCryptoChainById(id)
	if err != nil {
		log.Println("failed to get crypto chain by ID: ", err)
		return nil, errs.HandleError(err)
	}

	response := convertToCryptoChainResponse(*cryptoChain)
	return &response, nil
}

func (u *cryptoChainUseCase) UpdateCryptoChain(id string, req *entities.CryptoChainUpdateRequest) (*entities.CryptoChainResponse, error) {
	if req == nil {
		return nil, errs.NewBadRequestError("update request cannot be nil")
	}

	existingChain, err := u.repo.GetCryptoChainById(id)
	if err != nil {
		log.Println("failed to get crypto chain: ", err)
		return nil, errs.HandleError(err)
	}

	if err := utils.ApplyUpdates(existingChain, req); err != nil {
		log.Println("failed to apply update crypto chain: ", err)
		return nil, errs.HandleError(err)
	}

	updatedChain, err := u.repo.UpdateCryptoChain(id, existingChain)
	if err != nil {
		log.Println("failed to update crypto chain: ", err)
		return nil, errs.HandleError(err)
	}

	response := convertToCryptoChainResponse(*updatedChain)
	return &response, nil
}

func convertToCryptoChainResponse(data entities.CryptoChain) entities.CryptoChainResponse {
	return entities.CryptoChainResponse{
		Id:        data.Id,
		Name:      data.Name,
		Symbol:    data.Symbol,
		ImagePath: data.ImagePath,
		CreatedAt: data.CreatedAt,
	}
}
