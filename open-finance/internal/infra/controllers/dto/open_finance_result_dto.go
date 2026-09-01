package dto

// Os DTOs abaixo embrulham as listas devolvidas pelas ferramentas por dimensão.
// O MCP valida o esquema de saída contra um objeto de topo, e os campos de
// identificação servem de metadado de rastreabilidade da consulta.

type BankStatementsResultDTO struct {
	CustomerID uint               `json:"customer_id"`
	Document   string             `json:"document"`
	Items      []BankStatementDTO `json:"items"`
}

type CashFlowAnalysesResultDTO struct {
	CustomerID uint                  `json:"customer_id"`
	Document   string                `json:"document"`
	Items      []CashFlowAnalysisDTO `json:"items"`
}

type RecurringTransactionsResultDTO struct {
	CustomerID uint                      `json:"customer_id"`
	Document   string                    `json:"document"`
	Items      []RecurringTransactionDTO `json:"items"`
}

type DataSharingConsentsResultDTO struct {
	CustomerID uint                    `json:"customer_id"`
	Document   string                  `json:"document"`
	Items      []DataSharingConsentDTO `json:"items"`
}
