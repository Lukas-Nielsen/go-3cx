package threecx

import (
	"time"

	"github.com/pquerna/otp/totp"
)

func GetOTP(secret string) (string, error) {
	if len(secret) == 0 {
		return "", nil
	}
	if s, err := totp.GenerateCode(secret, time.Now()); err != nil {
		return "", err
	} else {
		return s, nil
	}
}
