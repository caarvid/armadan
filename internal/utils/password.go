package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	HashIter    = 2
	HashMemory  = 64 * 1024
	HashThreads = 1
	KeyLen      = 128
	SaltLen     = 128
)

var (
	ErrInvalidFormat       = errors.New("Invalid password format")
	ErrIncompatibleVersion = errors.New("Incompatible argon2 version")
)

type Password struct {
	iter    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
	Hash    []byte
	Salt    []byte
}

func (p *Password) Encode() string {
	h := base64.RawStdEncoding.EncodeToString(p.Hash)
	s := base64.RawStdEncoding.EncodeToString(p.Salt)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iter, p.threads, s, h)
}

func (p *Password) Compare(pass string) (bool, error) {
	h, err := GenerateHash(pass, p.Salt)

	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare(p.Hash, h.Hash) == 1 {
		return true, nil
	}

	return false, nil
}

func GenerateHash(pass string, salt []byte) (*Password, error) {
	var err error

	if salt == nil {
		salt, err = generateSalt(SaltLen)
	}

	if err != nil {
		return nil, err
	}

	saltLen, err := intToUint32(len(salt))
	if err != nil {
		return nil, err
	}

	p := &Password{
		iter:    HashIter,
		memory:  HashMemory,
		threads: HashThreads,
		keyLen:  KeyLen,
		saltLen: saltLen,
	}

	p.Hash = argon2.IDKey([]byte(pass), salt, p.iter, p.memory, p.threads, p.keyLen)
	p.Salt = salt

	return p, nil
}

func DecodeHash(encodedHash string) (*Password, error) {
	vals := strings.Split(encodedHash, "$")

	if len(vals) != 6 {
		return nil, ErrInvalidFormat
	}

	var version int

	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, err
	}

	if version != argon2.Version {
		return nil, ErrIncompatibleVersion
	}

	p := &Password{}

	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iter, &p.threads); err != nil {
		return nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, err
	}

	saltLen, err := intToUint32(len(salt))
	if err != nil {
		return nil, err
	}

	p.Salt = salt
	p.saltLen = saltLen

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, err
	}

	p.Hash = hash
	p.keyLen = saltLen

	return p, nil
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)

	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func intToUint32(val int) (uint32, error) {
	if val < 0 || val > math.MaxUint32 {
		return 0, errors.New("could not convert int to uint32")
	}

	return uint32(val), nil
}
