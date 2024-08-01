package task

import (
	"context"
	"net"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
)

const LogTaskCreateAccessLog = "logTask.CreateAccessLog"

func (t *logTask) RegisterLogTask(srv MachineryServer) error {
	if err := srv.HandleFunc(LogTaskCreateAccessLog, t.CreateAccessLog); err != nil {
		return err
	}

	return nil
}

func (t *logTask) CreateAccessLog(remoteAddr string) error {
	ip, _, _ := net.SplitHostPort(remoteAddr)

	log := &v1.AccessLog{
		Ip: ip,
	}

	record, err := t.mmdb.Country(net.ParseIP(remoteAddr))
	if err == nil {
		log.CountryName = record.Country.Names["en"]
		log.CountryCode = record.Country.IsoCode
	}

	if _, err = t.au.Create(context.TODO(), log); err != nil {
		return v1.ErrorInternal("write access log failed")
	}

	return nil
}
