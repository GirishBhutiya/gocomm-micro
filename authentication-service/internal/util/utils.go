package util

import (
	"authentication/internal/urlsigner"
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetVeifyEmailLink generates a link which send via email to verify your email.
func GetVeifyEmailLink(email, frontEndDomain string) string {
	return fmt.Sprintf("%s/verify-email?email=%s", frontEndDomain, email)
}
func GetFullVerifyEmailLink(email, frontEndDomain, hashSecretKeyVerifyEmail string) string {
	link := fmt.Sprintf("%s/verify-email?email=%s", frontEndDomain, email)
	sign := urlsigner.Signer{
		Secret: []byte(fmt.Sprintf("%s%s", hashSecretKeyVerifyEmail, email)),
	}
	return sign.GenerateTokenFromString(link)
}

// GetForgotPasswordLink generates a link which send via email to reset your password.
func GetForgotPasswordLink(email, frontEndDomain string) string {
	return fmt.Sprintf("%s/reset-password?email=%s", frontEndDomain, email)
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandStringBytesMaskImprSrcUnsafe(6))
}

// RandStringBytesMaskImprSrcUnsafe generates a random string with provided length
func RandStringBytesMaskImprSrcUnsafe(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
