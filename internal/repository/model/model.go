package model

type (
	// User модель пользователя для redis.
	UserRedis struct {
		ID        string `redis:"id"`
		Name      string `redis:"name"`
		Email     string `redis:"email"`
		Role      string `redis:"role"`
		Password  string `redis:"password"`
		CreatedAt string `redis:"created_at"`
	}
)
