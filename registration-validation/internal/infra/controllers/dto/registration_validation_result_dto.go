package dto

// Embrulham as saídas das ferramentas por dimensão. O MCP valida o esquema
// contra um objeto de topo; os campos de identificação são metadado de
// rastreabilidade da consulta.

type DocumentValidationsResultDTO struct {
	Traceable
	CustomerID uint                    `json:"customer_id"`
	Document   string                  `json:"document"`
	Items      []DocumentValidationDTO `json:"items"`
}

type FiscalRegularitiesResultDTO struct {
	Traceable
	CustomerID uint                  `json:"customer_id"`
	Document   string                `json:"document"`
	Items      []FiscalRegularityDTO `json:"items"`
}

type EmploymentLinkValidationsResultDTO struct {
	Traceable
	CustomerID uint                          `json:"customer_id"`
	Document   string                        `json:"document"`
	Items      []EmploymentLinkValidationDTO `json:"items"`
}

type ComplianceChecksResultDTO struct {
	Traceable
	CustomerID uint                 `json:"customer_id"`
	Document   string               `json:"document"`
	Items      []ComplianceCheckDTO `json:"items"`
}
