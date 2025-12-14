package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"ivory/src/clients/auth"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

var ErrUrlNotSpecified = errors.New("url is not specified")
var ErrBindDNNotSpecified = errors.New("BindDN is not specified")
var ErrBindPassNotSpecified = errors.New("BindPass is not specified")
var ErrBaseDNNotSpecified = errors.New("BaseDN is not specified")
var ErrFilterNotSpecified = errors.New("filter is not specified")
var ErrSchemeNotLdaps = errors.New("scheme is not ldaps")
var ErrConfigNotConfigured = errors.New("config is not configured")
var ErrUserNotFoundOrTooMany = errors.New("user isn't found or too many entries returned")
var ErrFailedToAppendCACert = errors.New("failed to append CA cert")

// NOTE: validate that is matches interface in compile-time
var _ auth.Provider[Config, Login] = (*Provider)(nil)

type Provider struct {
	config *Config
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Configured() bool {
	return p.config != nil
}

func (p *Provider) SetConfig(config Config) error {
	if config.Url == "" {
		return ErrUrlNotSpecified
	}
	if config.BindDN == "" {
		return ErrBindDNNotSpecified
	}
	if config.BindPass == "" {
		return ErrBindPassNotSpecified
	}
	if config.BaseDN == "" {
		return ErrBaseDNNotSpecified
	}
	if config.Filter == "" {
		return ErrFilterNotSpecified
	}
	if config.Tls != nil {
		if !strings.HasPrefix(config.Url, "ldaps://") {
			return ErrSchemeNotLdaps
		}
	}
	p.config = &config
	return nil
}

func (p *Provider) DeleteConfig() {
	p.config = nil
}

func (p *Provider) Connect(config Config) error {
	conn, err := p.getConnection(config)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return err
}

func (p *Provider) Verify(subject Login) (string, error) {
	if p.config == nil {
		return "", ErrConfigNotConfigured
	}
	conn, err := p.getConnection(*p.config)
	defer func(conn *ldap.Conn) {
		if conn == nil {
			return
		}
		err := conn.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(conn)
	if err != nil {
		return "", err
	}

	searchRequest := ldap.NewSearchRequest(
		p.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(p.config.Filter, subject.Username),
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(result.Entries) != 1 {
		return "", ErrUserNotFoundOrTooMany
	}

	userDn := result.Entries[0].DN
	err = conn.Bind(userDn, subject.Password)
	if err != nil {
		return "", err
	}

	return subject.Username, nil
}

func (p *Provider) getConnection(config Config) (*ldap.Conn, error) {
	dialer := &net.Dialer{Timeout: 3 * time.Second}

	tlsConfig := &tls.Config{}
	if config.Tls != nil && config.Tls.CaCert != "" {
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM([]byte(config.Tls.CaCert)); !ok {
			return nil, ErrFailedToAppendCACert
		}
		tlsConfig.RootCAs = certPool
	} else {
		tlsConfig.InsecureSkipVerify = true
	}

	options := []ldap.DialOpt{
		ldap.DialWithDialer(dialer),
		ldap.DialWithTLSConfig(tlsConfig),
	}

	conn, err := ldap.DialURL(config.Url, options...)
	if err != nil {
		return conn, err
	}
	conn.SetTimeout(3 * time.Second)
	err = conn.Bind(config.BindDN, config.BindPass)
	if err != nil {
		return conn, err
	}
	return conn, nil
}
