package models

import (
	"chatdan_backend/config"
	"gorm.io/gorm"
	"testing"
)

// TestMigrate is a test function for migrate from database to meilisearch
func TestMigrate(t *testing.T) {
	config.InitConfig()
	InitDB()

	// migrate box
	var boxes []Box
	result := DB.FindInBatches(&boxes, 1000, func(tx *gorm.DB, batch int) error {
		var boxSearchModels []BoxSearchModel
		for _, box := range boxes {
			boxSearchModels = append(boxSearchModels, box.ToBoxSearchModel())
		}
		err := SearchAddOrReplaceInBatch(boxSearchModels)
		if err != nil {
			t.Error(err)
		}
		return nil
	})
	if result.Error != nil {
		t.Error(result.Error)
	}
}
