package dto

// Embrulham as saídas das ferramentas por dimensão. O MCP valida o esquema
// contra um objeto de topo; os campos de identificação são metadado de
// rastreabilidade da consulta.

type CustomerRelationshipResultDTO struct {
	CustomerID   uint                     `json:"customer_id"`
	Document     string                   `json:"document"`
	Relationship *CustomerRelationshipDTO `json:"relationship"` // null quando não-cliente
}

type ContractedProductsResultDTO struct {
	CustomerID uint                   `json:"customer_id"`
	Document   string                 `json:"document"`
	Items      []ContractedProductDTO `json:"items"`
}

type InternalPaymentRecordsResultDTO struct {
	CustomerID uint                       `json:"customer_id"`
	Document   string                     `json:"document"`
	Items      []InternalPaymentRecordDTO `json:"items"`
}

type PreApprovedLimitsResultDTO struct {
	CustomerID uint                  `json:"customer_id"`
	Document   string                `json:"document"`
	Items      []PreApprovedLimitDTO `json:"items"`
}

type IncomeDeclarationsResultDTO struct {
	CustomerID uint                   `json:"customer_id"`
	Document   string                 `json:"document"`
	Items      []IncomeDeclarationDTO `json:"items"`
}
