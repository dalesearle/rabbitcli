package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"sort"
	"strings"
)

/*
The key is generated based on most of the keys in cluster, if any
cluster information changed beside the user and password, the user
and password for each cluster needs to be re encrypted.
 */
func CreateKey() []byte {
	hash := sha256.New()
	h, err := os.Hostname()
	cobra.CheckErr(err)
	hash.Write([]byte(h))
	d, err := os.UserHomeDir()
	cobra.CheckErr(err)
	hash.Write([]byte(d))
	keys := viper.AllKeys()
	sort.Sort(ByKey(keys))
	for _, k := range keys {
		if strings.Contains(k, USER) || strings.Contains(k, PASSWORD) || k == WORKING_CLUSTER{
			continue
		}
		hash.Write([]byte(k))
		hash.Write([]byte(viper.GetString(k)))
	}
	key := hash.Sum(nil)
	return key[:]
}

func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(text []byte) ([]byte, error) {
	key := CreateKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
