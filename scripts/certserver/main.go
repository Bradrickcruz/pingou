// certserver sobe um servidor HTTPS local com um certificado de validade
// escolhida, para testar o aviso de expiração do certificado SSL sem esperar
// 30 dias. Cria uma CA de teste e grava o certificado dela em PEM, para o
// Pingou confiar nela (no Linux: SSL_CERT_FILE=<ca.pem>).
//
//	go run ./scripts/certserver --days 25 --port 9443
//
// A CA é gravada em --out (padrão .bench/ca.pem) e reaproveitada nas próximas
// execuções, então dá para reiniciar com outro --days sem reiniciar o Pingou.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	days := flag.Float64("days", 25, "dias até o certificado do site expirar (aceita fração, ex.: 0.5)")
	port := flag.Int("port", 9443, "porta HTTPS")
	out := flag.String("out", ".bench/ca.pem", "onde gravar (ou ler) o certificado da CA de teste")
	status := flag.Int("status", http.StatusOK, "status HTTP das respostas")
	flag.Parse()

	ca, caKey, err := loadOrCreateCA(*out)
	if err != nil {
		log.Fatal(err)
	}

	notAfter := time.Now().Add(time.Duration(*days * float64(24*time.Hour))).Truncate(time.Second)
	leaf, err := siteCert(ca, caKey, notAfter)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:      fmt.Sprintf(":%d", *port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{leaf}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(*status)
			fmt.Fprintf(w, "certserver: certificado expira em %s\n", notAfter.UTC().Format(time.RFC3339))
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("https://localhost:%d — certificado expira em %s (%.1f dias)", *port, notAfter.UTC().Format(time.RFC3339), *days)
	log.Printf("CA de teste: %s (use SSL_CERT_FILE=%s no Pingou)", *out, *out)
	log.Fatal(srv.ListenAndServeTLS("", ""))
}

// loadOrCreateCA lê a CA de path (certificado) e path+".key" (chave), ou cria
// uma nova válida por 1 ano.
func loadOrCreateCA(path string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	keyPath := path + ".key"
	certPEM, errCert := os.ReadFile(path)
	keyPEM, errKey := os.ReadFile(keyPath)
	if errCert == nil && errKey == nil {
		cb, _ := pem.Decode(certPEM)
		kb, _ := pem.Decode(keyPEM)
		if cb == nil || kb == nil {
			return nil, nil, errors.New("CA inválida em " + path)
		}
		cert, err := x509.ParseCertificate(cb.Bytes)
		if err != nil {
			return nil, nil, err
		}
		key, err := x509.ParseECPrivateKey(kb.Bytes)
		if err != nil {
			return nil, nil, err
		}
		return cert, key, nil
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{CommonName: "Pingou certserver CA (teste)"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, nil, err
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, nil, err
	}
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o644); err != nil {
		return nil, nil, err
	}
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600); err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

// siteCert emite o certificado de localhost/127.0.0.1 assinado pela CA.
func siteCert(ca *x509.Certificate, caKey *ecdsa.PrivateKey, notAfter time.Time) (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: "localhost"},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.IPv6loopback},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca, &key.PublicKey, caKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, nil
}
