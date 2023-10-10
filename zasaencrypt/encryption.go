package zasaencrypt

import (
	aes256 "github.com/mervick/aes-everywhere/go/aes256"
)

type ZasaEnryption struct{}

var SECRET = "QarA3pCp833HzNLEUJYgTb7PpZWB6Uy18h1kVaBCqA2Twgq8egnvXf05rpiDu08f"

func (c ZasaEnryption) SecureData(secureData string, encrypt bool) string {
	var data string
	if encrypt {
		data = aes256.Encrypt(secureData, SECRET)
	} else {
		data = aes256.Decrypt(secureData, SECRET)
	}
	return data
}
