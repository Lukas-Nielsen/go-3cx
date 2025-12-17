package threecx

import (
	"errors"
	"time"

	"github.com/pquerna/otp/totp"
)

func getOTP(secret string) (string, error) {
	sec := time.Now().Second()
	if sec >= 30 {
		sec = sec - 30
	}
	if sec > 10 {
		time.Sleep(time.Second * time.Duration(31-sec))
	}

	if secret == "" {
		return "", errors.New("secret cannot be empty")
	}
	return totp.GenerateCode(secret, time.Now())
}
