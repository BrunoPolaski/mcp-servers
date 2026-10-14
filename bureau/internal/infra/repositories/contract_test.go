package repositories

import "github.com/BrunoPolaski/bureau/internal/infra/repositories/interfaces"

var (
	_ interfaces.PersonRepository  = (*gormPersonRepository)(nil)
	_ interfaces.UserRepository    = (*gormUserRepository)(nil)
	_ interfaces.AddressRepository = (*gormAddressRepository)(nil)
	_ interfaces.AnalystRepository = (*gormAnalystRepository)(nil)
	_ interfaces.ApiKeyRepository  = (*gormApiKeyRepository)(nil)
	_ interfaces.AdminRepository   = (*gormAdminRepository)(nil)
	_ interfaces.SessionRepository = (*gormSessionRepository)(nil)
)
