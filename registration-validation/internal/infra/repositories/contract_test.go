package repositories

import "github.com/BrunoPolaski/registration-validation/internal/infra/repositories/interfaces"

var (
	_ interfaces.PersonRepository                   = (*gormPersonRepository)(nil)
	_ interfaces.DocumentValidationRepository       = (*gormDocumentValidationRepository)(nil)
	_ interfaces.FiscalRegularityRepository         = (*gormFiscalRegularityRepository)(nil)
	_ interfaces.EmploymentLinkValidationRepository = (*gormEmploymentLinkValidationRepository)(nil)
	_ interfaces.ComplianceCheckRepository          = (*gormComplianceCheckRepository)(nil)
)
