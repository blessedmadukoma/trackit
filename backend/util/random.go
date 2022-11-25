package util

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates random numbers between minimum and maximum
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomString generates a set of random string with length of length
func RandomString(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		character := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(character)
	}

	return sb.String()

}

// RandomUser generates a random user name
func RandomUser() string {
	return RandomString(5)
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

const numbers = "0123456789"

// RandomMobileNumber generates a random mobile or phone number
func RandomMobileNumber() string {
	var mobile strings.Builder
	mobile.WriteString("081")
	for i := 0; i < 8; i++ {
		number := numbers[rand.Intn(len(numbers))]
		mobile.WriteByte(number)
	}

	return mobile.String()
}

// RandomPassword generates random password
func RandomPassword() string {
	password := RandomString(6)
	fmt.Println("Password before hash:", password)
	hashed, err := HashPassword(password)
	if err != nil {
		log.Fatal("error hashing password:", err)
	}
	return hashed
}
