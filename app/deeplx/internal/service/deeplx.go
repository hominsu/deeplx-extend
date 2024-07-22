package service

import (
	"context"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
)

func (s *DeepLXService) Translate(ctx context.Context, req *v1.TranslateRequest) (*v1.TranslationResult, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, v1.ErrorInternal("no deadline in context")
	}

	return s.ts.TranslateByDeepLX(
		req.GetSourceLang(),
		req.GetTargetLang(),
		req.GetText(),
		"",
		s.pool.Get(),
		deadline,
	)
}
