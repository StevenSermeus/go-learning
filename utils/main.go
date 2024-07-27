package utils

import (
	"crypto/rand"
	"encoding/binary"
	"unicode/utf8"
) 

func secureRandomInt(min int, max int) (int, error) {
	//Using the crypto/rand package to generate a random number between min and max
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	randomNumber := int(binary.BigEndian.Uint64(b[:]) % uint64(max-min) + uint64(min))
	return randomNumber, nil
}

func secureRandomUppercase() (rune, error) {
	randomNumber, err := secureRandomInt(1, 26)
	if err != nil {
		return 0, err
	}
	return rune(randomNumber % 26 + 65), nil
}

func GeneratePassPhrase(passPhrase ...string) (string, error) {
	if(len(passPhrase) < 1) {
		return GeneratePassPhrase("")
	}
	pass := passPhrase[0]
	// Base case for the recursion
	if utf8.RuneCountInString(pass) >= 24 {
		return pass, nil
	}

	//Not the best way to do this but it works
	if utf8.RuneCountInString(pass) == 4 || utf8.RuneCountInString(pass) == 9 || utf8.RuneCountInString(pass) == 14 || utf8.RuneCountInString(pass) == 19 {
		return GeneratePassPhrase(pass + "-")
	}

	char, err := secureRandomUppercase()
	if err != nil {
		return "", err
	}

	return GeneratePassPhrase(pass + string(char))
}

func GenerateApiKey(api_key ...string) (string, error) {
	if(len(api_key) < 1) {
		return GenerateApiKey("")
	}
	api := api_key[0]
	// Base case for the recursion
	if utf8.RuneCountInString(api) >= 64 {
		return api, nil
	}
	
	char, err := secureRandomUppercase()
	if err != nil {
		return "", err
	}

	return GenerateApiKey(api + string(char))
}

