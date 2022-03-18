type redis struct{ pool *redisClient.Pool }

func New(host, port, password string) (storage.Service, error) {
	pool := &redisClient.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redisClient.Conn, error) {
			return redisClient.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	return &redis{pool}, nil
}

func (r *redis) isUsed(id uint64) bool {
	conn := r.pool.Get()
	defer conn.Close()

	exists, err := redisClient.Bool(conn.Do("EXISTS", "Shorturl:"+strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (r *redis) Save(url string, expires time.Time) (string, error) {
	
}
// func (r *redis) Save(url string, expires time.Time) (uint64, error) {
// 	conn := r.pool.Get()
// 	defer conn.Close()

// 	id, err := redisClient.Uint64(conn.Do("INCR", "NextId"))
// 	if err != nil {
// 		return 0, err
// 	}

// 	if r.isUsed(id) {
// 		return 0, errors.New("id already used")
// 	}

// 	if _, err := conn.Do("SET", "Shorturl:"+strconv.FormatUint(id, 10), url, "EX", int(expires.Sub(time.Now()).Seconds())); err != nil {
// 		return 0, err
// 	}

// 	return id, nil
// }