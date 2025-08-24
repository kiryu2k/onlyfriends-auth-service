package hasher

import (
	"context"
	"crypto/rand"
	"encoding/base64"

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
	hash := argon2.IDKey([]byte(v), salt, s.iterCount, s.memory, s.parallel, s.tagLength)
	return base64.StdEncoding.EncodeToString(hash), nil
}

func (s service) VerifyHash(ctx context.Context, v string, hash string) (bool, error) {
	hashToCompare, err := s.Hash(ctx, v)
	if err != nil {
		return false, errors.WithMessage(err, "hash")
	}
	return hash == hashToCompare, nil
}

func (s service) generateSalt() ([]byte, error) {
	salt := make([]byte, s.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, errors.WithMessage(err, "read rand")
	}
	return salt, nil
}
