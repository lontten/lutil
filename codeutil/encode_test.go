package codeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSHA256HashCode(t *testing.T) {
	as := assert.New(t)
	h1 := GetSHA256HashCode([]byte("abc"), "salt")
	h2 := GetSHA256HashCode([]byte("abc"), "salt")
	as.Equal(h1, h2)
	as.Len(h1, 64)
	as.NotEqual(h1, GetSHA256HashCode([]byte("abc"), "other"))
}

func TestBase64RoundTrip(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)
	enc := Base64Encode("hello")
	dec, err := Base64Decode(enc)
	req.NoError(err)
	as.Equal("hello", dec)
}

func TestMD5(t *testing.T) {
	as := assert.New(t)
	as.Equal("098f6bcd4621d373cade4e832627b4f6", MD5("test"))
}

func TestRandomStr(t *testing.T) {
	as := assert.New(t)
	s := RandomStr(16)
	as.Len(s, 16)
	as.Panics(func() { RandomStr(0) })
}

func TestRandomNumAndCaptcha(t *testing.T) {
	as := assert.New(t)
	as.Len(RandomNum(6), 6)
	as.Len(GenCaptcha(4), 4)
}

func TestHashPasswordVerifyPassword(t *testing.T) {
	as := assert.New(t)
	req := require.New(t)
	hash, err := HashPassword("secret")
	req.NoError(err)
	as.True(VerifyPassword("secret", hash))
	as.False(VerifyPassword("wrong", hash))
	as.False(VerifyPassword("secret", "not-a-bcrypt-hash"))
}

func TestEnPwdCheckPassword(t *testing.T) {
	as := assert.New(t)
	cipher := EnPwd("secret")
	as.True(CheckPassword("secret", cipher))
	as.False(CheckPassword("wrong", cipher))
	as.False(CheckPassword("secret", "short"))
}

func TestRandomTimeID(t *testing.T) {
	as := assert.New(t)
	id := RandomTimeID32()
	as.Len(id, 32)
	numID := RandomTimeNumberID32()
	as.Len(numID, 32)
}
