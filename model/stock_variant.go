package model

import "strings"

type StockVariant struct {
	Id          int32  `gorm:"primary_key"`
	Variant     string `gorm:"column:variant"`
	Permissions string `gorm:"column:permissions"`
	// Permission is comma separted string
	Abilities string `gorm:"column:abilities"`
	// Abilities is pipe separated with ? to
	// separated ability with the permission requirement.
	// cred_harvest?read_file,upload_file|....
	// EncryptionCert string
}

func (s StockVariant) GetPermissionList() []string {
	return strings.Split(s.Permissions, ",")
}

func (s StockVariant) GetAbilities() (result map[string][]string) {
	abilities := strings.Split(s.Abilities, "|")
	var i = 0
	result = map[string][]string{}
	for i < len(abilities) {
		ability := strings.Split(abilities[i], "?")
		result[ability[0]] = strings.Split(ability[1], ",")
		i++
	}
	return
}

type NewStockVariant struct {
	Variant        string
	Permissions    string
	Abilities      string
	EncryptionCert string
}

type StockVariantRepository interface {
	Create(nv *StockVariant) (int32, error)
	FindStockByVariant(v string) (*StockVariant, error)
}

type StockVariantUsecase interface {
	RegisterStockVariant(NewStockVariant) error
	GetAllRegisteredStockVariant() ([]StockVariant, error)
}
