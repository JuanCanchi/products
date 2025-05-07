package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/juancanchi/products/internal/domain"
)

type ProductRepository interface {
	Save(ctx context.Context, p *domain.Product) error
	FindAll(ctx context.Context) ([]*domain.Product, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id string) error
}

type ProductUsecase struct {
	repo ProductRepository
}

func NewProductUsecase(repo ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) Create(ctx context.Context, p *domain.Product) error {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return u.repo.Save(ctx, p)
}

func (u *ProductUsecase) List(ctx context.Context) ([]*domain.Product, error) {
	return u.repo.FindAll(ctx)
}

func (u *ProductUsecase) ListByUser(ctx context.Context, userID string) ([]*domain.Product, error) {
	return u.repo.FindByUserID(ctx, userID)
}

func (u *ProductUsecase) Update(ctx context.Context, p *domain.Product, userID string) error {
	// Obtener el producto actual
	existing, err := u.repo.FindByID(ctx, p.ID)
	if err != nil {
		return err
	}

	// Verificar si pertenece al usuario
	if existing.UserID != userID {
		return errors.New("no autorizado")
	}

	// Aplicar cambios permitidos
	existing.Title = p.Title
	existing.Description = p.Description
	existing.Price = p.Price
	existing.ImageURL = p.ImageURL

	return u.repo.Update(ctx, existing)
}

func (u *ProductUsecase) Delete(ctx context.Context, id string, userID string) error {
	product, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if product.UserID != userID {
		return errors.New("no autorizado")
	}

	return u.repo.Delete(ctx, id)
}

func (u *ProductUsecase) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return u.repo.FindByID(ctx, id)
}
