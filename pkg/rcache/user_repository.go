package rcache

//Redis缓存管理, 文章之类的 不细分
type UserRedisRepository struct {
	BaseRedisRepository
}

func NewUserRedisRepository() *UserRedisRepository {
	return &UserRedisRepository{BaseRedisRepository: BaseRedisRepository{}}
}
