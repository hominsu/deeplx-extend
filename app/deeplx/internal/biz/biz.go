package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewAccessLogUsecase,
)

const (
	// MaxPageSize is the maximum page size that can be returned by a List call. Requesting page sizes larger than
	// this value will return, at most, MaxPageSize entries.
	MaxPageSize = 1024

	// MaxBatchCreateSize is the maximum number of entries that can be created by a single BatchCreate call. Requests
	// exceeding this batch size will return an error.
	MaxBatchCreateSize = 1024

	// MaxBatchUpdateSize is the maximum number of entries that can be updated by a single BatchUpdate call. Requests
	// exceeding this batch size will return an error.
	MaxBatchUpdateSize = 1024
)
