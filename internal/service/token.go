package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// newPingToken gera o segredo da URL de ping: 32 bytes aleatórios (256 bits)
// em base64url, seguro para usar direto no caminho /ping/{token}.
func newPingToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate ping token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}
