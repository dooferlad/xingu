package session

import (
	"github.com/adrg/xdg"
	"testing"
	"time"
)

func newTestingOTP() *OTP {
	return &OTP{
		now:      time.Now().Unix,
		cacheDir: xdg.CacheFile,
	}
}

func TestOTPNoCached(t *testing.T) {
	otp := newTestingOTP()
	...
}
