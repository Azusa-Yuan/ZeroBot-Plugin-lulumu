package pcrjjctw

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	// MessagePack 是一个轻量级的、速度快的二进制序列化格式

	"github.com/vmihailenco/msgpack"
)

type pcrclient struct {
	viewer_id  string
	short_udid string
	udid       string
	header     interface{}
	proxy      string
	platform   string
}

func (p pcrclient) makemd5(str string) string {
	salt := "r!I@nt8e5i="

	// Concatenate the string and the salt
	concatenated := str + salt

	// Compute MD5 hash
	hasher := md5.New()
	hasher.Write([]byte(concatenated))
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashHex := hex.EncodeToString(hashBytes)

	return hashHex
}

func (p pcrclient) createkey() []byte {
	hexChars := "0123456789abcdef"
	randomBytes := make([]byte, 32)

	for i := 0; i < 32; i++ {
		randomBytes[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return randomBytes
}

func (p pcrclient) getiv() []byte {
	cleanedUdid := strings.ReplaceAll(p.udid, "-", "")
	shortenedUdid := cleanedUdid[:16]
	return []byte(shortenedUdid)
}

func (p pcrclient) pack(data map[string]interface{}, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesIV := p.getiv()
	aesBlockMode := cipher.NewCBCEncrypter(block, aesIV)

	var packedData bytes.Buffer
	// 将数据编码为 MessagePack 格式
	encoder := msgpack.NewEncoder(&packedData).UseCompactEncoding(true)
	if err := encoder.Encode(data); err != nil {
		return nil, nil, err
	}

	packed := packedData.Bytes()

	padLen := aes.BlockSize - len(data)%aes.BlockSize
	padded := bytes.Repeat([]byte{byte(padLen)}, padLen)

	encrypted := make([]byte, len(padded))
	aesBlockMode.CryptBlocks(encrypted, padded)

	return packed, append(encrypted, key...), nil
}

func (p pcrclient) encrypt(data string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesIV := p.getiv()
	aesBlockMode := cipher.NewCBCEncrypter(block, aesIV)

	// Pad the data as per PKCS#7 standard
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	paddedData := bytes.Repeat([]byte{byte(padLen)}, padLen)

	encrypted := make([]byte, len(paddedData))
	aesBlockMode.CryptBlocks(encrypted, paddedData)

	return append(encrypted, key...), nil
}

func (p pcrclient) decrypt(data []byte) ([]byte, []byte, error) {
	key := data[len(data)-32:]
	encryptedData := data[:len(data)-32]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesIV := p.getiv()
	aesBlockMode := cipher.NewCBCEncrypter(block, aesIV)

	decrypted := make([]byte, len(encryptedData))
	aesBlockMode.CryptBlocks(decrypted, encryptedData)

	return decrypted, key, nil
}

func (p pcrclient) unpack(data []byte) ([]byte, []byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, nil, err
	}

	aesKey := decodedData[len(decodedData)-32:]
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, err
	}

	IV := p.getiv()
	mode := cipher.NewCBCEncrypter(block, IV)
	mode.CryptBlocks(decodedData[:len(decodedData)-32], decodedData[:len(decodedData)-32])

	// 去除填充
	padding := int(data[len(data)-1])
	decryptedData := data[:len(data)-padding]
}

func (p pcrclient) CallApi(ApiUrl string, request map[string]interface{}, delay float32) string {
	if delay > 1 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	key := p.createkey()

}
