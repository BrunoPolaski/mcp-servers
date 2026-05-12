package thirdparty

import (
	"sync"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/cache"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/database"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/mailer"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ThirdPartyFactory struct {
	db         *gorm.DB
	redis      *redis.Client
	jwtAdapter jwt.JWT
	mailer     mailer.Mailer
	dbOnce     sync.Once
	redisOnce  sync.Once
	jwtOnce    sync.Once
	mailerOnce sync.Once
	queueOnce  sync.Once
}

func NewThirdPartyFactory() *ThirdPartyFactory {
	return &ThirdPartyFactory{}
}

func (f *ThirdPartyFactory) DB() *gorm.DB {
	f.dbOnce.Do(func() {
		f.db = database.NewGormAdapter().Connect()
	})
	return f.db
}

func (f *ThirdPartyFactory) Redis() *redis.Client {
	f.redisOnce.Do(func() {
		f.redis = cache.NewRedisAdapter().Connect()
	})
	return f.redis
}

func (f *ThirdPartyFactory) JWT() jwt.JWT {
	f.jwtOnce.Do(func() {
		f.jwtAdapter = jwt.NewJWTAdapter()
	})
	return f.jwtAdapter
}

func (f *ThirdPartyFactory) Mailer() mailer.Mailer {
	f.mailerOnce.Do(func() {
		f.mailer = mailer.NewGoMailAdapter()
	})
	return f.mailer
}
