package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"ivory/clients/storage"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

func createTestCertService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "cert-service-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	oldWd, errWd := os.Getwd()
	if errWd != nil {
		t.Fatalf("failed to get working dir: %v", errWd)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(oldWd)
	})

	fileStorage := storage.NewFileStorage("cert", "")
	repository := NewRepository(storage.NewDbBucket[Cert](db, "Cert"), fileStorage)

	return NewService(repository)
}

// generateCertKeyPair writes a self-signed certificate/key pair to files under dir
// and returns their paths, so GetTLSConfigCertificates/GetTLSConfigRootCA have real
// PEM content to parse.
func generateCertKeyPair(t *testing.T, dir string) (certPath string, keyPath string) {
	t.Helper()

	priv, errKey := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if errKey != nil {
		t.Fatalf("failed to generate key: %v", errKey)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:         true,
	}

	der, errCreate := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if errCreate != nil {
		t.Fatalf("failed to create certificate: %v", errCreate)
	}

	certPath = filepath.Join(dir, "test.crt")
	certOut, errCertFile := os.Create(certPath)
	if errCertFile != nil {
		t.Fatalf("failed to create cert file: %v", errCertFile)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		t.Fatalf("failed to encode cert: %v", err)
	}
	certOut.Close()

	keyBytes, errMarshal := x509.MarshalECPrivateKey(priv)
	if errMarshal != nil {
		t.Fatalf("failed to marshal key: %v", errMarshal)
	}
	keyPath = filepath.Join(dir, "test.key")
	keyOut, errKeyFile := os.Create(keyPath)
	if errKeyFile != nil {
		t.Fatalf("failed to create key file: %v", errKeyFile)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}); err != nil {
		t.Fatalf("failed to encode key: %v", err)
	}
	keyOut.Close()

	return certPath, keyPath
}

func TestServiceCreateAndGet(t *testing.T) {
	s := createTestCertService(t)
	dir := t.TempDir()
	certPath, _ := generateCertKeyPair(t, dir)

	t.Run("creates a cert by path and can get it back", func(t *testing.T) {
		created, err := s.Create(certPath, CLIENT_CA, PATH)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if created.Path != certPath {
			t.Fatalf("expected path %s, got %s", certPath, created.Path)
		}

		list, errList := s.List()
		if errList != nil {
			t.Fatalf("expected no error, got %v", errList)
		}
		var found uuid.UUID
		for k, v := range list {
			if v.Path == certPath {
				found = uuid.MustParse(k)
			}
		}
		got, errGet := s.Get(found)
		if errGet != nil {
			t.Fatalf("expected no error, got %v", errGet)
		}
		if got.Type != CLIENT_CA {
			t.Fatalf("expected type CLIENT_CA, got %v", got.Type)
		}
	})

	t.Run("unsupported file format is rejected", func(t *testing.T) {
		_, err := s.Create(filepath.Join(dir, "test.pem"), CLIENT_CA, PATH)
		if err == nil {
			t.Fatalf("expected an error for unsupported format")
		}
	})

	t.Run("nonexistent path is rejected", func(t *testing.T) {
		_, err := s.Create(filepath.Join(dir, "missing.crt"), CLIENT_CA, PATH)
		if err != ErrNoSuchFile {
			t.Fatalf("expected ErrNoSuchFile, got %v", err)
		}
	})

	t.Run("uploaded file is created via file storage", func(t *testing.T) {
		created, err := s.Create("anything.crt", CLIENT_CERT, UPLOAD)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, errStat := os.Stat(created.Path); errStat != nil {
			t.Fatalf("expected uploaded file to exist at %s: %v", created.Path, errStat)
		}
	})
}

func TestServiceListByType(t *testing.T) {
	s := createTestCertService(t)
	dir := t.TempDir()
	certPath, keyPath := generateCertKeyPair(t, dir)

	if _, err := s.Create(certPath, CLIENT_CA, PATH); err != nil {
		t.Fatalf("failed to create ca: %v", err)
	}
	if _, err := s.Create(keyPath, CLIENT_KEY, PATH); err != nil {
		t.Fatalf("failed to create key: %v", err)
	}

	list, err := s.ListByType(CLIENT_KEY)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
}

func TestServiceDeleteAndDeleteAll(t *testing.T) {
	s := createTestCertService(t)
	dir := t.TempDir()
	certPath, _ := generateCertKeyPair(t, dir)

	created, err := s.Create(certPath, CLIENT_CA, PATH)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}
	list, _ := s.List()
	var id uuid.UUID
	for k := range list {
		id = uuid.MustParse(k)
	}

	t.Run("delete removes the entry", func(t *testing.T) {
		if err := s.Delete(id); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, errGet := s.Get(id); errGet == nil {
			t.Fatalf("expected error getting deleted cert")
		}
	})

	t.Run("delete all clears everything", func(t *testing.T) {
		if _, err := s.Create(created.Path, CLIENT_CA, PATH); err != nil {
			t.Fatalf("failed to recreate: %v", err)
		}
		if err := s.DeleteAll(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		list, errList := s.List()
		if errList != nil {
			t.Fatalf("expected no error, got %v", errList)
		}
		if len(list) != 0 {
			t.Fatalf("expected empty list, got %v", list)
		}
	})
}

func TestServiceGetTLSConfigRootCA(t *testing.T) {
	s := createTestCertService(t)
	dir := t.TempDir()
	certPath, _ := generateCertKeyPair(t, dir)

	t.Run("nil id returns nil pool without error", func(t *testing.T) {
		pool, err := s.GetTLSConfigRootCA(nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if pool != nil {
			t.Fatalf("expected nil pool, got %v", pool)
		}
	})

	t.Run("valid id returns a populated pool", func(t *testing.T) {
		created, err := s.Create(certPath, CLIENT_CA, PATH)
		if err != nil {
			t.Fatalf("failed to create: %v", err)
		}
		list, _ := s.List()
		var id uuid.UUID
		for k, v := range list {
			if v.Path == created.Path {
				id = uuid.MustParse(k)
			}
		}
		pool, errPool := s.GetTLSConfigRootCA(&id)
		if errPool != nil {
			t.Fatalf("expected no error, got %v", errPool)
		}
		if pool == nil {
			t.Fatalf("expected a populated pool")
		}
	})

	t.Run("unknown id propagates the error", func(t *testing.T) {
		unknown := uuid.New()
		_, err := s.GetTLSConfigRootCA(&unknown)
		if err == nil {
			t.Fatalf("expected an error for unknown id")
		}
	})
}

func TestServiceGetTLSConfigCertificates(t *testing.T) {
	s := createTestCertService(t)
	dir := t.TempDir()
	certPath, keyPath := generateCertKeyPair(t, dir)

	createdCert, err := s.Create(certPath, CLIENT_CERT, PATH)
	if err != nil {
		t.Fatalf("failed to create cert: %v", err)
	}
	createdKey, err := s.Create(keyPath, CLIENT_KEY, PATH)
	if err != nil {
		t.Fatalf("failed to create key: %v", err)
	}
	list, _ := s.List()
	var certId, keyId uuid.UUID
	for k, v := range list {
		if v.Path == createdCert.Path {
			certId = uuid.MustParse(k)
		}
		if v.Path == createdKey.Path {
			keyId = uuid.MustParse(k)
		}
	}

	t.Run("both nil returns an empty slice without error", func(t *testing.T) {
		certs, errCerts := s.GetTLSConfigCertificates(nil, nil)
		if errCerts != nil {
			t.Fatalf("expected no error, got %v", errCerts)
		}
		if len(certs) != 0 {
			t.Fatalf("expected no certificates, got %v", certs)
		}
	})

	t.Run("only cert id set is rejected", func(t *testing.T) {
		_, errCerts := s.GetTLSConfigCertificates(&certId, nil)
		if errCerts != ErrMutualTLSRequiresBothCertAndKey {
			t.Fatalf("expected ErrMutualTLSRequiresBothCertAndKey, got %v", errCerts)
		}
	})

	t.Run("only key id set is rejected", func(t *testing.T) {
		_, errCerts := s.GetTLSConfigCertificates(nil, &keyId)
		if errCerts != ErrMutualTLSRequiresBothCertAndKey {
			t.Fatalf("expected ErrMutualTLSRequiresBothCertAndKey, got %v", errCerts)
		}
	})

	t.Run("both ids set returns a loaded certificate", func(t *testing.T) {
		certs, errCerts := s.GetTLSConfigCertificates(&certId, &keyId)
		if errCerts != nil {
			t.Fatalf("expected no error, got %v", errCerts)
		}
		if len(certs) != 1 {
			t.Fatalf("expected 1 certificate, got %d", len(certs))
		}
	})
}

func TestServiceEnrichTLSConfig(t *testing.T) {
	s := createTestCertService(t)
	dir := t.TempDir()
	certPath, keyPath := generateCertKeyPair(t, dir)

	createdCert, _ := s.Create(certPath, CLIENT_CERT, PATH)
	createdKey, _ := s.Create(keyPath, CLIENT_KEY, PATH)
	list, _ := s.List()
	var certId, keyId uuid.UUID
	for k, v := range list {
		if v.Path == createdCert.Path {
			certId = uuid.MustParse(k)
		}
		if v.Path == createdKey.Path {
			keyId = uuid.MustParse(k)
		}
	}

	t.Run("nil certs clears the config", func(t *testing.T) {
		config := &tls.Config{}
		var configPtr = config
		if err := s.EnrichTLSConfig(&configPtr, nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if configPtr != nil {
			t.Fatalf("expected config to be nil")
		}
	})

	t.Run("populated certs enrich a nil config", func(t *testing.T) {
		var configPtr *tls.Config
		certs := &Certs{ClientCertId: &certId, ClientKeyId: &keyId}
		if err := s.EnrichTLSConfig(&configPtr, certs); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if configPtr == nil {
			t.Fatalf("expected config to be initialized")
		}
		if len(configPtr.Certificates) != 1 {
			t.Fatalf("expected 1 certificate, got %d", len(configPtr.Certificates))
		}
	})

	t.Run("propagates root ca errors", func(t *testing.T) {
		var configPtr *tls.Config
		unknown := uuid.New()
		certs := &Certs{ClientCAId: &unknown}
		if err := s.EnrichTLSConfig(&configPtr, certs); err == nil {
			t.Fatalf("expected an error")
		}
	})
}
