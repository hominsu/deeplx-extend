package service

import (
	"context"
	"net"

	"github.com/go-kratos/kratos/v2/metadata"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/pkgs/middleware"
)

func (s *DeepLXService) Translate(ctx context.Context, req *v1.TranslateRequest) (*v1.TranslationResult, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, v1.ErrorInternal("no deadline in context")
	}

	token := s.cs.GetAuth().GetToken()
	if token != "" {
		tokenInQuery := req.GetToken()
		var tokenInHeader string
		if md, ok := metadata.FromClientContext(ctx); ok {
			tokenInHeader = md.Get(string(middleware.ContextKeyAuthToken))
		}

		if tokenInHeader == "" && tokenInQuery == "" {
			return nil, v1.ErrorUnauthenticated("Missing access token")
		}

		if tokenInHeader != token && tokenInQuery != token {
			return nil, v1.ErrorPermissionDenied("Invalid access token")
		}
	}

	var remoteAddr string
	if md, ok := metadata.FromClientContext(ctx); ok {
		remoteAddr = md.Get(string(middleware.ContextKeyRemoteAddr))
	}

	record, err := s.mmdb.Country(net.ParseIP(remoteAddr))
	if err != nil {
		return nil, v1.ErrorInternal("missing ip record")
	}

	if _, err := s.au.Create(ctx, &v1.AccessLog{
		Ip:          remoteAddr,
		CountryName: record.Country.Names["en"],
		CountryCode: record.Country.IsoCode,
	}); err != nil {
		return nil, v1.ErrorInternal("write access log failed")
	}

	payload := req.GetPayload()
	return s.ts.TranslateByDeepLX(
		payload.GetSourceLang(),
		payload.GetTargetLang(),
		payload.GetText(),
		"",
		s.pool.Get(),
		deadline,
	)
}
