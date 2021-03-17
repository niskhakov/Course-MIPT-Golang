package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"strings"
	"time"
)

type SignMethod string

const (
	HS256 SignMethod = "HS256"
	HS512 SignMethod = "HS512"
)

var (
	ErrInvalidSignMethod      = errors.New("invalid sign method")
	ErrSignatureInvalid       = errors.New("signature invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrSignMethodMismatched   = errors.New("sign method mismatched")
	ErrConfigurationMalformed = errors.New("configuration malformed")
	ErrInvalidToken           = errors.New("invalid token")
)

type header struct {
	Alg string `json:"alg"`
	Type string `json:"typ"`
}

type payload struct {
	Data interface{} `json:"d"`
	Exp *int64 `json:"exp,omitempty"`
}

func int64pointer(x int64) *int64 {
	return &x
}

func int64value(x *int64) int64 {
	if x == nil {
		return 0
	}
	return *x
}

func getCryptoMethod(signMethod SignMethod) (_ func() hash.Hash, error bool) {
	var cryptoMethod func() hash.Hash
	switch signMethod {
	case HS256:
		cryptoMethod = sha256.New
	case HS512:
		cryptoMethod = sha512.New
	default:
		return nil, true
	}
	return cryptoMethod, false
}

func getControlSum(b []byte, c *config) ([]byte, error) {
	cryptoMethod, err := getCryptoMethod(c.SignMethod)
	if err {
		return nil, ErrInvalidSignMethod
	}

	h := hmac.New(cryptoMethod, []byte(c.Key))
	h.Write(b)
	bytesum := h.Sum(nil)
	sum := make([]byte, base64.RawStdEncoding.EncodedLen(len(bytesum)))
	base64.RawURLEncoding.Encode(sum, bytesum)

	return sum, nil
}

func marshalAndEncodePart(d interface{}) string {
	dataJson, _ := json.Marshal(d)
	return base64.RawURLEncoding.EncodeToString(dataJson)
}

func decodeAndUnmarshalPart(part string, d interface{}) error {
	partJson, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return ErrInvalidToken
	}

	err = json.Unmarshal(partJson, &d)
	if err != nil {
		return ErrInvalidToken
	}

	return nil
}

func checkControlSum(token []byte, c *config) error {
	lastDotIndex := bytes.LastIndex(token, []byte("."))

	actualSum := token[lastDotIndex + 1:]
	infoToVerify := token[:lastDotIndex]
	expectedSum, err := getControlSum(infoToVerify, c)

	if err != nil {
		return err
	}

	if len(actualSum) != len(expectedSum) {
		return ErrSignMethodMismatched
	}

	if !bytes.Equal(actualSum, expectedSum) {
		return ErrSignatureInvalid
	}

	return nil
}

func Encode(data interface{}, opts ...Option) ([]byte, error) {
	var c config

	// Применение опций
	for _, opt := range opts {
			opt(&c)
	}

	// Формирование заголовка
	hdrEncoded := marshalAndEncodePart(header{
		Alg: string(c.SignMethod),
		Type: "JWT",
	})
	
	// Формирование тела
	var exp *int64
	if c.Expires != nil {
		exp = int64pointer(c.Expires.Unix())
	}

	if c.TTL != nil {
		exp = int64pointer(timeFunc().Add(*c.TTL).Unix())
	}

	if c.TTL != nil && c.Expires != nil {
		return nil, ErrConfigurationMalformed
	}

	if c.Expires != nil && c.Expires.Before(timeFunc()) {
		return nil, ErrConfigurationMalformed
	}
	
	pldEncoded := marshalAndEncodePart(payload{
		Data: data,
		Exp: exp,
	})

	// Формирование контрольной суммы
	var b bytes.Buffer
	b.WriteString(hdrEncoded)
	b.WriteString(".")
	b.WriteString(pldEncoded)

	shasum, err := getControlSum(b.Bytes(), &c)
	if err != nil {
		return nil, err
	}

	b.WriteString(".")
	b.Write(shasum)

	return b.Bytes(), nil
}

func Decode(token []byte, data interface{}, opts ...Option) error {

	var c config
	for _, opt := range opts {
		opt(&c)
	}

	parts := strings.Split(string(token), ".")
	
	if len(parts) != 3 {
		return ErrInvalidToken
	}

	var hdr header
	if err := decodeAndUnmarshalPart(parts[0], &hdr); err != nil {
		return err
	}
	
	var pld payload
	if err := decodeAndUnmarshalPart(parts[1], &pld); err != nil {
		return err
	}

	if err := checkControlSum(token, &c); err != nil {
		return err
	}

	if pld.Exp != nil && timeFunc().After(time.Unix(int64value(pld.Exp), 0)) {
		return ErrTokenExpired
	}

	// Bad code: need to be revised - writing to interface() variable
	dataJson, _ := json.Marshal(pld.Data)
	json.Unmarshal(dataJson, data)

	return nil
}

// To mock time in tests
var timeFunc = time.Now
