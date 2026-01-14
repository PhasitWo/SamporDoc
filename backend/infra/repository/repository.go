package repository

import (
	"SamporDoc/backend/infra/model"
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repo struct {
	ctx context.Context
	DB  *gorm.DB
}

func NewRepo(ctx context.Context, DB *gorm.DB) *Repo {
	return &Repo{ctx, DB}
}

func (r *Repo) newError(err error) error {
	return fmt.Errorf("[REPO]: %w", err)
}

func (r *Repo) GetAllCustomers() ([]model.Customer, error) {
	var customers []model.Customer
	customers, err := gorm.G[model.Customer](r.DB).Find(r.ctx)
	if err != nil {
		return customers, r.newError(err)
	}
	return customers, nil
}

func (r *Repo) GetCustomerByName(name string) (model.Customer, error) {
	customer, err := gorm.G[model.Customer](r.DB).Where("name = ?", name).First(r.ctx)
	if err != nil {
		return customer, r.newError(err)
	}
	return customer, nil
}

func (r *Repo) CreateCustomer(customer *model.Customer) error {
	err := gorm.G[model.Customer](r.DB).Create(r.ctx, customer)
	if err != nil {
		return r.newError(err)
	}
	return nil
}

func (r *Repo) UpdateCustomerByID(customer *model.Customer) (retCustomer model.Customer, err error) {
	_, err = gorm.G[model.Customer](r.DB).Where("id = ?", customer.ID).Select("*").Updates(r.ctx, *customer)
	if err != nil {
		return retCustomer, r.newError(err)
	}
	retCustomer, err = gorm.G[model.Customer](r.DB).Where("id = ?", customer.ID).First(r.ctx)
	if err != nil {
		return retCustomer, r.newError(err)
	}
	return retCustomer, nil
}

func (r *Repo) GetAllShops() ([]model.Shop, error) {
	var shops []model.Shop
	shops, err := gorm.G[model.Shop](r.DB).Order("sorting_level asc").Find(r.ctx)
	if err != nil {
		return shops, r.newError(err)
	}
	return shops, nil
}

func (r *Repo) CreateShop(shop *model.Shop) error {
	err := gorm.G[model.Shop](r.DB).Create(r.ctx, shop)
	if err != nil {
		return r.newError(err)
	}
	return err
}

func (r *Repo) UpdateShopBySlug(shop *model.Shop) (retShop model.Shop, err error) {
	_, err = gorm.G[model.Shop](r.DB).Where("slug = ?", shop.Slug).Updates(r.ctx, *shop)
	if err != nil {
		return retShop, r.newError(err)
	}
	retShop, err = gorm.G[model.Shop](r.DB).Where("id = ?", shop.ID).First(r.ctx)
	if err != nil {
		return retShop, r.newError(err)
	}
	return retShop, nil
}

type Status string

const (
	SUCCESS Status = "SUCCESS"
	ERROR   Status = "ERROR"
)

func (r *Repo) CreateLog(action string, status Status, correlationID *int64, data ...any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return r.newError(err)
	}
	dataStr := string(bytes)
	err = gorm.G[model.Log](r.DB).Create(r.ctx, &model.Log{
		Action: action,
		Status: string(status),
		Data:   datatypes.JSON(dataStr),
	})
	return err
}
