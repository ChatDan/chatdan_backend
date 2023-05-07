package utilsx

import (
	"math/rand"
	"strings"
	"time"
)

var rander = rand.New(rand.NewSource(time.Now().UnixNano()))

func SecretGenerator(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var secretBuilder strings.Builder
	for i := 0; i < length; i++ {
		secretBuilder.WriteByte(alphabet[rander.Intn(len(alphabet))])
	}
	return secretBuilder.String()
}
