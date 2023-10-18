package helper

import "github.com/tyler-smith/go-bip39"

func GenerateSecureSeedPhrase() string {
	entropy, _ := bip39.NewEntropy(128) // 128 bits for a 12-word mnemonic

	seedPhrase, _ := bip39.NewMnemonic(entropy)
	return seedPhrase
}
