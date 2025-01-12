package threecx

import (
	"time"

	"github.com/pquerna/otp/totp"
)

func getOTP(secret string) (string, error) {
	if len(secret) == 0 {
		return "", nil
	}
	if s, err := totp.GenerateCode(secret, time.Now()); err != nil {
		return "", err
	} else {
		return s, nil
	}
}
