package models

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type CAEntry struct {
	Id     string `json:"id" gorm:"primaryKey"`
	Parent string `json:"parent"`

	Children []CAEntry `json:"children" gorm:"foreignKey:Parent"`

	Certificate string `json:"certificate"`

	// We do not want to expose our privatekeys to even the administrator.
	PrivateKey string `json:"-"`

	Certificates []string `json:"certificates" gorm:"type:string[]"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CAChain struct {
	Id     string `json:"id" gorm:"primaryKey"`
	Parent string

	Children []CAEntry `json:"children" gorm:"foreignKey:Parent"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getCAEntryTree(ctx context.Context, db *gorm.DB, id string) (*CAEntry, error) {
	var entry CAEntry
	if err := db.WithContext(ctx).Preload("Children").Where("id = ?", id).First(&entry).Error; err != nil {
		return nil, err
	}
	for i := range entry.Children {
		child, err := getCAEntryTree(ctx, db, entry.Children[i].Id)
		if err != nil {
			return nil, err
		}
		entry.Children[i] = *child
	}
	return &entry, nil
}

func GetCAChain(ctx context.Context, db *gorm.DB) (*CAChain, error) {
	var chain CAChain
	if err := db.WithContext(ctx).Preload("Children").First(&chain).Error; err != nil {
		return nil, err
	}

	for i := range chain.Children {
		child, err := getCAEntryTree(ctx, db, chain.Children[i].Id)
		if err != nil {
			return nil, err
		}
		chain.Children[i] = *child
	}

	return &chain, nil
}
