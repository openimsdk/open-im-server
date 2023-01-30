package controller

import (
	"Open_IM/pkg/common/db/relation"
	"context"
	"gorm.io/gorm"
)

type BlackInterface interface {
	Create(ctx context.Context, blacks []*relation.Black) (err error)
	Delete(ctx context.Context, blacks []*relation.Black) (err error)
	UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, blacks []*relation.Black) (err error)
	Find(ctx context.Context, blacks []*relation.Black) (blackList []*relation.Black, err error)
	Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.Black, err error)
	FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*relation.Black, err error)
}

type BlackController struct {
	database BlackDatabaseInterface
}

func NewBlackController(db *gorm.DB) *BlackController {
	return &BlackController{database: NewBlackDatabase(db)}
}
func (f *BlackController) Create(ctx context.Context, blacks []*relation.Black) (err error) {
	return f.database.Create(ctx, blacks)
}
func (f *BlackController) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
	return f.database.Delete(ctx, blacks)
}
func (f *BlackController) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	return f.database.UpdateByMap(ctx, ownerUserID, blockUserID, args)
}
func (f *BlackController) Update(ctx context.Context, blacks []*relation.Black) (err error) {
	return f.database.Update(ctx, blacks)
}
func (f *BlackController) Find(ctx context.Context, blacks []*relation.Black) (blackList []*relation.Black, err error) {
	return f.database.Find(ctx, blacks)
}
func (f *BlackController) Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.Black, err error) {
	return f.database.Take(ctx, ownerUserID, blockUserID)
}
func (f *BlackController) FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*relation.Black, err error) {
	return f.database.FindByOwnerUserID(ctx, ownerUserID)
}

type BlackDatabaseInterface interface {
	Create(ctx context.Context, blacks []*relation.Black) (err error)
	Delete(ctx context.Context, blacks []*relation.Black) (err error)
	UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, blacks []*relation.Black) (err error)
	Find(ctx context.Context, blacks []*relation.Black) (blackList []*relation.Black, err error)
	Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.Black, err error)
	FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*relation.Black, err error)
}

type BlackDatabase struct {
	sqlDB *relation.Black
}

func NewBlackDatabase(db *gorm.DB) *BlackDatabase {
	sqlDB := relation.NewBlack(db)
	database := &BlackDatabase{
		sqlDB: sqlDB,
	}
	return database
}

func (f *BlackDatabase) Create(ctx context.Context, blacks []*relation.Black) (err error) {
	return f.sqlDB.Create(ctx, blacks)
}
func (f *BlackDatabase) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
	return f.sqlDB.Delete(ctx, blacks)
}
func (f *BlackDatabase) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	return f.sqlDB.UpdateByMap(ctx, ownerUserID, blockUserID, args)
}
func (f *BlackDatabase) Update(ctx context.Context, blacks []*relation.Black) (err error) {
	return f.sqlDB.Update(ctx, blacks)
}
func (f *BlackDatabase) Find(ctx context.Context, blacks []*relation.Black) (blackList []*relation.Black, err error) {
	return f.sqlDB.Find(ctx, blacks)
}
func (f *BlackDatabase) Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.Black, err error) {
	return f.sqlDB.Take(ctx, ownerUserID, blockUserID)
}
func (f *BlackDatabase) FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*relation.Black, err error) {
	return f.sqlDB.FindByOwnerUserID(ctx, ownerUserID)
}
