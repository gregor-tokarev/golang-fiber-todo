package config

import "github.com/joho/godotenv"

type Config struct {
	PostgresHost            string
	PostgresPort            string
	PostgresUser            string
	PostgresPassword        string
	PostgresDB              string
	JwtAccessSecret         string
	JwtRefreshSecret        string
	GoogleOauthClientID     string
	GoogleOauthClientSecret string
	FrontendUrl             string
}

var Cfg *Config

func init() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	Cfg = &Config{
		PostgresHost:            env["POSTGRES_HOST"],
		PostgresUser:            env["POSTGRES_USER"],
		PostgresPassword:        env["POSTGRES_PASSWORD"],
		PostgresDB:              env["POSTGRES_DB"],
		PostgresPort:            env["POSTGRES_PORT"],
		JwtAccessSecret:         env["JWT_ACCESS_SECRET"],
		JwtRefreshSecret:        env["JWT_REFRESH_SECRET"],
		GoogleOauthClientID:     env["GOOGLE_OAUTH_CLIENT_ID"],
		GoogleOauthClientSecret: env["GOOGLE_OAUTH_CLIENT_SECRET"],
		FrontendUrl:             env["FRONTEND_URL"],
	}
}
