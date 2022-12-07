package utils

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
)

func GetSHA256Digest(input []byte) string {
	hash := sha256.New()
	_, err := hash.Write(input)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func GetNewId(input []byte) string {
	return uuid.New().String() + "|" + GetSHA256Digest(input)
}
