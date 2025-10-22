package ldap

import (
	"errors"
	"fmt"
	"ivory/src/clients/auth"
	"log/slog"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

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
		return errors.New("url is not specified")
	}
	if config.BindDN == "" {
		return errors.New("BindDN is not specified")
	}
	if config.BindPass == "" {
		return errors.New("BindPass is not specified")
	}
	if config.BaseDN == "" {
		return errors.New("BaseDN is not specified")
	}
	if config.Filter == "" {
		return errors.New("filter is not specified")
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
		return "", errors.New("config is not configured")
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
		return "", errors.New("user isn't found or too many entries returned")
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
	conn, err := ldap.DialURL(config.Url, ldap.DialWithDialer(dialer))
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
