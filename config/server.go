package config

import "time"

type Server struct {
	Host           string
	Port           int `validate:"required"`
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	JwtIssuer      string `validate:"required"`
	JwtSecretKey   string `validate:"required"`
	JwtExpiresHour int    `validate:"required lte=24"`
}
