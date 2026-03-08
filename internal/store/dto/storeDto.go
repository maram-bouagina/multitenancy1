package dto

type CreateStoreRequest struct {
	Name      string  `json:"name"      validate:"required"`
	Slug      string  `json:"slug"      validate:"required"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	Logo      *string `json:"logo"`
	Currency  string  `json:"currency"  validate:"required"`
	Timezone  string  `json:"timezone"  validate:"required"`
	Language  string  `json:"language"  validate:"required"`
	TaxNumber *string `json:"tax_number"`
}

type UpdateStoreRequest struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	Logo      *string `json:"logo"`
	Currency  *string `json:"currency"`
	Timezone  *string `json:"timezone"`
	Language  *string `json:"language"`
	TaxNumber *string `json:"tax_number"`
	Status    *string `json:"status" validate:"omitempty,oneof=active suspended inactive"`
}

type StoreResponse struct {
	ID        string  `json:"id"`
	TenantID  string  `json:"tenant_id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	Logo      *string `json:"logo"`
	Currency  string  `json:"currency"`
	Timezone  string  `json:"timezone"`
	Language  string  `json:"language"`
	TaxNumber *string `json:"tax_number"`
	Status    string  `json:"status"`
}
