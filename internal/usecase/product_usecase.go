package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juancanchi/jujuy-market/products/internal/domain"
)

type ProductRepository interface {
	Save(ctx context.Context, p *domain.Product) error
	FindAll(ctx context.Context) ([]*domain.Product, error)
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
