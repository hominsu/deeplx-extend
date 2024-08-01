package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/task"
	"github.com/oio-network/deeplx-extend/app/deeplx/pkgs/middleware"
	"github.com/oio-network/deeplx-extend/app/deeplx/pkgs/msg"
	"github.com/oio-network/deeplx-extend/pkgs/machinery"
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

	if remoteAddr != "" {
		bytes, err := msg.Marshal(&task.AccessLogParams{
			RemoteAddr: remoteAddr,
			CreatedAt:  timestamppb.Now(),
		})
		if err != nil {
			s.log.Error(err)
		}

		err = s.ms.NewTask(task.LogTaskCreateAccessLog, machinery.WithArgument("[]byte", bytes))
		if err != nil {
			s.log.Error(err)
		}
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
