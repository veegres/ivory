package ldap

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

type Provider struct {
	config *Config
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) SetConfig(config *Config) error {
	if config == nil {
		return errors.New("config is not configured")
	}
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

	// NOTE: validate connection
	conn, err := p.getConnection(config)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}

	p.config = config
	return nil
}

func (p *Provider) DeleteConfig() {
	p.config = nil
}

func (p *Provider) Verify(login Login) (string, error) {
	if p.config == nil {
		return "", errors.New("config is not configured")
	}
	conn, err := p.getConnection(p.config)
	defer func(conn *ldap.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(conn)
	if err != nil {
		return "", err
	}

	filter := p.config.Filter
	if filter == "" {
		filter = "(uid=%s)"
	}

	searchRequest := ldap.NewSearchRequest(
		p.config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(filter, login.Username),
		[]string{"dn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", errors.New("user isn't found or too many entries returned")
	}

	userDn := sr.Entries[0].DN
	err = conn.Bind(userDn, login.Password)
	if err != nil {
		return "", err
	}

	return login.Username, nil
}

func (p *Provider) getConnection(config *Config) (*ldap.Conn, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := ldap.DialURL(config.Url, ldap.DialWithDialer(dialer))
	if err != nil {
		return nil, err
	}
	//conn.SetTimeout(3 * time.Second)

	err = conn.Bind(config.BindDN, config.BindPass)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
