package hasher

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

type service struct {
	iterCount  uint32
	tagLength  uint32
	saltLength uint32
	memory     uint32
	parallel   uint8
}

func New() service {
	return service{
		iterCount:  8,
		tagLength:  128,
		saltLength: 16,
		memory:     64 * 1024,
		parallel:   2,
	}
}

func (s service) Hash(_ context.Context, v string) (string, error) {
	salt, err := s.generateSalt()
	if err != nil {
		return "", errors.WithMessage(err, "generate salt")
	}
	hash := s.hash(v, salt)
	return s.joinSaltToHash(hash, salt), nil
}

func (s service) VerifyHash(ctx context.Context, v string, encodedHash string) (bool, error) {
	expHash, salt, err := s.splitHashAndSalt(encodedHash)
	if err != nil {
		return false, errors.WithMessage(err, "split hash and salt")
	}
	computedHash := s.hash(v, salt)
	return subtle.ConstantTimeCompare(expHash, computedHash) == 1, nil
}

func (s service) generateSalt() ([]byte, error) {
	salt := make([]byte, s.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, errors.WithMessage(err, "read rand")
	}
	return salt, nil
}

func (s service) hash(v string, salt []byte) []byte {
	return argon2.IDKey([]byte(v), salt, s.iterCount, s.memory, s.parallel, s.tagLength)
}

const (
	sep = "$$"
)

func (s service) joinSaltToHash(hash []byte, salt []byte) string {
	return fmt.Sprintf("%s%s%s",
		base64.StdEncoding.EncodeToString(hash),
		sep,
		base64.StdEncoding.EncodeToString(salt),
	)
}

func (s service) splitHashAndSalt(encodedHash string) (hash []byte, salt []byte, err error) {
	hashAndSalt := strings.Split(encodedHash, sep)
	if len(hashAndSalt) != 2 {
		return nil, nil, errors.Errorf("unexpected count of encoded hash parts: %d", len(hashAndSalt))
	}

	hash, err = base64.StdEncoding.DecodeString(hashAndSalt[0])
	if err != nil {
		return nil, nil, errors.WithMessage(err, "decode base64 hash")
	}
	salt, err = base64.StdEncoding.DecodeString(hashAndSalt[1])
	if err != nil {
		return nil, nil, errors.WithMessage(err, "decode base64 salt")
	}

	return hash, salt, nil
}
