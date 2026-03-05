package repository

import (
	"SamporDoc/backend/infra/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	ctx context.Context
	DB  *gorm.DB
}

func NewCustomerRepo(ctx context.Context, DB *gorm.DB) *CustomerRepo {
	return &CustomerRepo{ctx, DB}
}

func (r *CustomerRepo) newError(err error) error {
	return fmt.Errorf("[REPO]: %w", err)
}

func (r *CustomerRepo) GetAllCustomers() ([]model.Customer, error) {
	var customers []model.Customer
	customers, err := gorm.G[model.Customer](r.DB).Find(r.ctx)
	if err != nil {
		return customers, r.newError(err)
	}
	return customers, nil
}

func (r *CustomerRepo) GetCustomerByName(name string) (model.Customer, error) {
	customer, err := gorm.G[model.Customer](r.DB).Where("name = ?", name).First(r.ctx)
	if err != nil {
		return customer, r.newError(err)
	}
	return customer, nil
}

func (r *CustomerRepo) CreateCustomer(customer *model.Customer) error {
	err := gorm.G[model.Customer](r.DB).Create(r.ctx, customer)
	if err != nil {
		return r.newError(err)
	}
	return nil
}

func (r *CustomerRepo) UpdateCustomerByID(customer *model.Customer) (retCustomer model.Customer, err error) {
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
