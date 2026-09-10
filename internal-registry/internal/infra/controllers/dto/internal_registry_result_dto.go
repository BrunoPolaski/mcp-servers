package dto

// Embrulham as saídas das ferramentas por dimensão. O MCP valida o esquema
// contra um objeto de topo; os campos de identificação são metadado de
// rastreabilidade da consulta.

type CustomerRelationshipResultDTO struct {
	Traceable
	CustomerID   uint                     `json:"customer_id"`
	Document     string                   `json:"document"`
	Relationship *CustomerRelationshipDTO `json:"relationship"` // null quando não-cliente
}

type ContractedProductsResultDTO struct {
	Traceable
	CustomerID uint                   `json:"customer_id"`
	Document   string                 `json:"document"`
	Items      []ContractedProductDTO `json:"items"`
}

type InternalPaymentRecordsResultDTO struct {
	Traceable
	CustomerID uint                       `json:"customer_id"`
	Document   string                     `json:"document"`
	Items      []InternalPaymentRecordDTO `json:"items"`
}

type PreApprovedLimitsResultDTO struct {
	Traceable
	CustomerID uint                  `json:"customer_id"`
	Document   string                `json:"document"`
	Items      []PreApprovedLimitDTO `json:"items"`
}

type IncomeDeclarationsResultDTO struct {
	Traceable
	CustomerID uint                   `json:"customer_id"`
	Document   string                 `json:"document"`
	Items      []IncomeDeclarationDTO `json:"items"`
}
