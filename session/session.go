package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/creachadair/otp"
	"github.com/spf13/viper"
)

type Session struct {
	*session.Session
}

func NewOTP() *OTP {
	return &OTP{
		now:      time.Now().Unix,
		cacheDir: xdg.CacheFile,
	}
}

type OTP struct {
	now      func() int64
	cacheDir func(string) (string, error)
}

func (o *OTP) GetOTP() (string, error) {
	otpSecret := viper.GetString("otp")
	var cfg otp.Config
	if err := cfg.ParseKey(otpSecret); err != nil {
		return "", err
	}

	fn, _ := xdg.CacheFile("xingu/aws-otp-used-ts.txt")
	b, err := ioutil.ReadFile(fn)
	var oldTS int
	if err == nil {
		oldTS, _ = strconv.Atoi(string(b))
	}

	ts := time.Now().Unix() / 30

	var printed bool
	for ts == int64(oldTS) {
		if !printed {
			fmt.Println("Waiting for OTP token (a new one is obtained every 30 seconds and can't be reused)")
			printed = true
		}

		progress := int(time.Now().Unix() % 30)
		fmt.Printf("[%-29s]\r", strings.Repeat("#", progress))

		time.Sleep(time.Second)
		ts = time.Now().Unix() / 30
	}

	if printed {
		fmt.Printf("\n")
	}

	ioutil.WriteFile(fn, []byte(fmt.Sprintf("%d", ts)), 0600)

	s := cfg.HOTP(uint64(ts))
	return s, nil
}

func New() (*Session, error) {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	otpSecret := viper.GetString("otp")

	if otpSecret != "" {
		otp := NewOTP()
		options.AssumeRoleTokenProvider = otp.GetOTP
	}

	creds, err := Load()
	if err == nil {
		options.Config.Credentials = credentials.NewStaticCredentialsFromCreds(*creds)
	}

	s, err := session.NewSessionWithOptions(options)
	return &Session{s}, err
}

type CachedCredential struct {
	credentials.Value
	ExpiresAt time.Time
}

func fileName() (string, error) {
	profile := os.Getenv("AWS_PROFILE")
	path := fmt.Sprintf("xingu/aws-cred-%s.json", profile)
	return xdg.CacheFile(path)
}

func Load() (*credentials.Value, error) {
	var c *CachedCredential

	fn, err := fileName()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	if c.ExpiresAt.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("credential expired")
	}

	return &c.Value, nil
}

func (s Session) SaveCreds() error {
	fn, err := fileName()
	if err != nil {
		return err
	}

	v, err := s.Config.Credentials.Get()
	if err != nil {
		return err
	}

	expiresAt, err := s.Config.Credentials.ExpiresAt()
	if err != nil {
		return err
	}

	creds, err := json.Marshal(CachedCredential{
		Value:     v,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fn, creds, 0600)
}
