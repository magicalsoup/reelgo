package auth

import (
	"crypto/rand"
	"math/big"
)



func Generate6DigitAuthCode() (string) {
	var code string = ""
	for i := 0; i < 6; i++ {
		d, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += d.String()
	}
	return code
}
