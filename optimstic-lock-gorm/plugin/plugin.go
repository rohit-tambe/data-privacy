package plugin

import (
	"errors"

	"gorm.io/gorm"
)

type MyModel struct {
	ID      int
	Version int
}

func (m *MyModel) BeforeSave(tx *gorm.DB) (err error) {
	// check if the model has a valid ID and Version field
	if m.ID == 0 || m.Version == 0 {
		return nil
	}

	// retrieve the current version of the model from the database
	var currentVersion int
	tx.Model(m).Select("version").Where("id = ?", m.ID).Scan(&currentVersion)

	// compare the current version with the expected version
	if currentVersion != m.Version {
		return errors.New("optimistic lock failed")
	}

	// increment the version field before saving the model
	m.Version++

	return nil
}
