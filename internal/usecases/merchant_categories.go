package usecases

import (
	"api-gateway/internal/entities"
	"api-gateway/internal/repositories"
	"api-gateway/pkg/errs"
	"log"
)

type (
	MerchantCategoryUseCase interface {
		GetMerchantCategories() ([]entities.MerchantCategoryResponse, error)
		GetMerchantCategoryById(id string) (*entities.MerchantCategoryResponse, error)
		CreateMerchantCategory(req entities.MerchantCategoryCreateRequest) (*entities.MerchantCategoryResponse, error)
		UpdateMerchantCategory(id string, req entities.MerchantCategoryUpdateRequest) (*entities.MerchantCategoryResponse, error)
		DeleteMerchantCategory(id string) error
	}

	merchantCategoryUseCase struct {
		repo repositories.MerchantCategoryRepository
	}
)

func NewMerchantCategoryUseCase(repo repositories.MerchantCategoryRepository) MerchantCategoryUseCase {
	return &merchantCategoryUseCase{
		repo: repo,
	}
}

func (u *merchantCategoryUseCase) GetMerchantCategories() ([]entities.MerchantCategoryResponse, error) {
	categories, err := u.repo.GetMerchantCategories()
	if err != nil {
		log.Println("failed to get merchant categories: ", err)
		return nil, errs.HandleError(err)
	}

	return convertToMerchantCategoryResponses(categories), nil
}

func (u *merchantCategoryUseCase) GetMerchantCategoryById(id string) (*entities.MerchantCategoryResponse, error) {
	category, err := u.repo.GetMerchantCategoryById(id)
	if err != nil {
		log.Println("failed to get merchant category by id: ", err)
		return nil, errs.HandleError(err)
	}

	return convertToMerchantCategoryResponse(category), nil
}

func (u *merchantCategoryUseCase) CreateMerchantCategory(req entities.MerchantCategoryCreateRequest) (*entities.MerchantCategoryResponse, error) {
	category, err := u.repo.CreateMerchantCategory(entities.MerchantCategory{
		Name: req.Name,
		Icon: req.Icon,
	})
	if err != nil {
		log.Println("failed to created merchant category: ", err)
		return nil, errs.HandleError(err)
	}

	return convertToMerchantCategoryResponse(category), nil
}

func (u *merchantCategoryUseCase) UpdateMerchantCategory(id string, req entities.MerchantCategoryUpdateRequest) (*entities.MerchantCategoryResponse, error) {
	category, err := u.repo.UpdateMerchantCategory(id, &entities.MerchantCategory{
		Name: *req.Name,
		Icon: req.Icon,
	})
	if err != nil {
		log.Println("failed to update merchant category: ", err)
		return nil, errs.HandleError(err)
	}

	return convertToMerchantCategoryResponse(category), nil
}

func (u *merchantCategoryUseCase) DeleteMerchantCategory(id string) error {
	if err := u.repo.DeleteMerchantCategory(id); err != nil {
		log.Println("failed to delete merchant category: ", err)
		return errs.HandleError(err)
	}

	return nil
}

func convertToMerchantCategoryResponses(categories []entities.MerchantCategory) []entities.MerchantCategoryResponse {
	responses := make([]entities.MerchantCategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = *convertToMerchantCategoryResponse(&c)
	}

	return responses
}

func convertToMerchantCategoryResponse(mc *entities.MerchantCategory) *entities.MerchantCategoryResponse {
	return &entities.MerchantCategoryResponse{
		Id:   mc.Id,
		Name: mc.Name,
		Icon: mc.Icon,
	}
}
