package checker

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// newCertServer sobe um servidor HTTPS cujo certificado (autoassinado, para
// 127.0.0.1) vence em notAfter, e devolve um checker que confia nele.
func newCertServer(t *testing.T, notAfter time.Time, status int) (*httptest.Server, *HTTPChecker) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "pingou-test"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              notAfter,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	}))
	srv.TLS = &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}}}
	srv.StartTLS()
	t.Cleanup(srv.Close)

	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	c := NewHTTPChecker(5, 10)
	c.client.Transport = &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}
	return srv, c
}

func TestHTTPCheckerReadsCertExpiry(t *testing.T) {
	notAfter := time.Now().Add(20 * 24 * time.Hour).Truncate(time.Second)

	for _, status := range []int{http.StatusOK, http.StatusInternalServerError} {
		srv, c := newCertServer(t, notAfter, status)
		m := &domain.Monitor{Type: domain.TypeHTTP, URL: srv.URL, TimeoutSeconds: 2}

		res := c.Check(context.Background(), m)
		if res.Success != (status == http.StatusOK) {
			t.Fatalf("status %d: success=%v", status, res.Success)
		}
		if res.CertNotAfter == nil || !res.CertNotAfter.Equal(notAfter) {
			t.Fatalf("status %d: CertNotAfter=%v, want %v", status, res.CertNotAfter, notAfter)
		}
		if res.CertNotAfter.Location() != time.UTC {
			t.Fatalf("status %d: CertNotAfter fora de UTC: %v", status, res.CertNotAfter.Location())
		}
		if res.CertURL != srv.URL {
			t.Fatalf("status %d: CertURL=%q, want %q", status, res.CertURL, srv.URL)
		}
	}
}

func TestHTTPCheckerNoCertWithoutTLS(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	res := NewHTTPChecker(5, 10).Check(context.Background(), &domain.Monitor{Type: domain.TypeHTTP, URL: srv.URL, TimeoutSeconds: 2})
	if !res.Success || res.CertNotAfter != nil || res.CertURL != "" {
		t.Fatalf("http://: %+v, want sucesso sem certificado", res)
	}
}

// Certificado vencido: o Go recusa o handshake, não há resposta e o campo fica
// vazio. O monitor cai pelo fluxo normal de falha.
func TestHTTPCheckerExpiredCertFails(t *testing.T) {
	srv, c := newCertServer(t, time.Now().Add(-time.Minute), http.StatusOK)

	res := c.Check(context.Background(), &domain.Monitor{Type: domain.TypeHTTP, URL: srv.URL, TimeoutSeconds: 2})
	if res.Success || res.ErrorMessage == nil || res.CertNotAfter != nil {
		t.Fatalf("vencido: %+v, want falha sem certificado", res)
	}
}
