package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juancanchi/products/internal/domain"
)

type CategoryRepository interface {
	Save(ctx context.Context, c *domain.Category) error
	FindAll(ctx context.Context) ([]*domain.Category, error)
	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id string) error
}

type CategoryUsecase struct {
	repo CategoryRepository
}

func NewCategoryUsecase(repo CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (u *CategoryUsecase) Create(ctx context.Context, c *domain.Category) error {
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	return u.repo.Save(ctx, c)
}

func (u *CategoryUsecase) List(ctx context.Context) ([]*domain.Category, error) {
	return u.repo.FindAll(ctx)
}

func (u *CategoryUsecase) Update(ctx context.Context, c *domain.Category) error {
	return u.repo.Update(ctx, c)
}

func (u *CategoryUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
