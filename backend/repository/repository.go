package repository

import (
	"SamporDoc/backend/model"
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

func (r *Repo) GetAllSchools() ([]model.School, error) {
	var schools []model.School
	schools, err := gorm.G[model.School](r.DB).Find(r.ctx)
	if err != nil {
		return schools, r.newError(err)
	}
	return schools, nil
}

func (r *Repo) CreateSchool(school *model.School) error {
	err := gorm.G[model.School](r.DB).Create(r.ctx, school)
	if err != nil {
		return r.newError(err)
	}
	return nil
}

func (r *Repo) GetAllShops() ([]model.Shop, error) {
	var shops []model.Shop
	shops, err := gorm.G[model.Shop](r.DB).Find(r.ctx)
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

func (r *Repo) UpdateShopBySlug(shop *model.Shop) error {
	_, err := gorm.G[model.Shop](r.DB).Where("slug = ?", shop.Slug).Updates(r.ctx, *shop)
	if err != nil {
		return r.newError(err)
	}
	return err
}

func (r *Repo) CreateLog(action string, data ...any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return r.newError(err)
	}
	dataStr := string(bytes)
	err = gorm.G[model.Log](r.DB).Create(r.ctx, &model.Log{
		Action: action,
		Data:   datatypes.JSON(dataStr),
	})
	return err
}
