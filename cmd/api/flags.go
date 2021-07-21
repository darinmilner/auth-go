package main

import "flag"

func getFlags() {
	flag.IntVar(&cfg.port, "port", 8000, "Server port")
	flag.StringVar(&cfg.env, "env", "dev", "App environment")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "60f19c8d43291f29e9585a6c6f272d33e6a83f7a209e9061b2b85afa5358ad49", "secret")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://root:password@localhost:5432/auth?sslmode=disable", "Postgres connection string")
	//flag.StringVar(&cfg.db.dsn, "dsn", "postgres://root:password@postgres-movies:5432/auth?sslmode=disable", "Postgres DOCKER connection string")
}
