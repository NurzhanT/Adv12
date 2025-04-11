package usecase

import "inventory/internal/domain"

type ProductUseCase interface {
	Create(p domain.Product)
	GetByID(id string) (domain.Product, error)
	Update(id string, p domain.Product) error
	Delete(id string) error
	List() ([]domain.Product, error)
}

func (p ProductUseCase) GetAll() (any, any) {
	panic("unimplemented")
}

type productUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(r domain.ProductRepository) ProductUseCase {
	return &productUseCase{repo: r}
}

func (uc *productUseCase) Create(p domain.Product) {
	uc.repo.Create(p)
}

func (uc *productUseCase) GetByID(id string) (domain.Product, error) {
	return uc.repo.GetByID(id)
}

func (uc *productUseCase) Update(id string, p domain.Product) error {
	return uc.repo.Update(id, p)
}

func (uc *productUseCase) Delete(id string) error {
	return uc.repo.Delete(id)
}

func (uc *productUseCase) List() ([]domain.Product, error) {
	return uc.repo.List()
}
