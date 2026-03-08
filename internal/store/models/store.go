package models

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;index"                       json:"tenant_id"`
	Name      string     `gorm:"type:varchar(255);not null"                     json:"name"`
	Slug      string     `gorm:"type:varchar(255);uniqueIndex;not null"          json:"slug"`
	Email     *string    `gorm:"type:varchar(255)"                              json:"email"`
	Phone     *string    `gorm:"type:varchar(20)"                               json:"phone"`
	Address   *string    `gorm:"type:text"                                      json:"address"`
	Logo      *string    `gorm:"type:varchar(500)"                              json:"logo"`
	Currency  string     `gorm:"type:varchar(10);not null;default:'EUR'"        json:"currency"`
	Timezone  string     `gorm:"type:varchar(100);not null;default:'UTC'"       json:"timezone"`
	Language  string     `gorm:"type:varchar(10);not null;default:'fr'"         json:"language"`
	TaxNumber *string    `gorm:"type:varchar(100)"                              json:"tax_number"`
	Status    string     `gorm:"type:varchar(20);default:'active';check:status IN ('active','suspended','inactive')" json:"status"`
	CreatedAt time.Time  `gorm:"type:timestamptz;autoCreateTime"                json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;autoUpdateTime"                json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamptz;index"                         json:"-"`
}

func (Store) TableName() string {
	return "stores"
}
