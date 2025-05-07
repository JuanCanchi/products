package postgres

import (
	"context"
	"github.com/juancanchi/products/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repositorio real usando GORM
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Save(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Conexión a PostgreSQL usando GORM
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
