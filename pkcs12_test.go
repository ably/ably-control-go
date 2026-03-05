package control

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"software.sslmate.com/src/go-pkcs12"
)

// generateTestP12 creates a self-signed certificate encoded as a valid
// PKCS12 (.p12) file suitable for the APNs upload endpoint.
func generateTestP12(t *testing.T, password string) []byte {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	p12Bytes, err := pkcs12.Legacy.Encode(key, cert, nil, password)
	require.NoError(t, err)

	return p12Bytes
}

func TestUploadPKCS12(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	p12 := generateTestP12(t, "test-password")

	result, err := client.UploadPKCS12(app.ID, p12, "test-password")
	assert.NoError(t, err)
	assert.Equal(t, app.ID, result.ID)
}

func TestUploadPKCS12InvalidFile(t *testing.T) {
	client, _ := newTestClient(t)

	_, err := client.UploadPKCS12("fake-app", []byte("not a p12 file"), "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PKCS12 file")
}

func TestUploadPKCS12WrongPassword(t *testing.T) {
	client, _ := newTestClient(t)

	p12 := generateTestP12(t, "correct-password")

	_, err := client.UploadPKCS12("fake-app", p12, "wrong-password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PKCS12 file")
}

func TestUploadPKCS12MissingBasicConstraints(t *testing.T) {
	client, _ := newTestClient(t)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	p12, err := pkcs12.Legacy.Encode(key, cert, nil, "test-password")
	require.NoError(t, err)

	_, err = client.UploadPKCS12("fake-app", p12, "test-password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing BasicConstraints extension")
}
