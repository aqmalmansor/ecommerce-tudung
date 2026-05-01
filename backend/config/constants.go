package config

import (
	"be/utils"
	"time"
)

var (
	CLIENT_URL             = utils.GetEnvString("CLIENT_URL", "http://localhost:3000")
	PORT                   = utils.GetEnvNumber("PORT", 8080)
	JWT_SECRET             = utils.GetEnvString("JWT_SECRET", "your-secret")
	JWT_CLAIMS_KEY_USER_ID = "userID"
	JWT_CLAIMS_KEY_EMAIL   = "email"
	JWT_EXPIRES_DURATION   = time.Hour
	JWT_EXPIRES_IN         = int(time.Hour.Seconds())
)
