package task

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/biz"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
)

type LogBatch struct {
	Logs []*v1.AccessLog
}

type LogSink struct {
	logChan        chan *v1.AccessLog
	autoCommitChan chan *LogBatch

	conf *conf.Data_Log
	au   *biz.AccessLogUseCase
	log  *log.Helper
}

func NewLogSink(
	conf *conf.Data,
	au *biz.AccessLogUseCase,
	logger log.Logger,
) (*LogSink, func()) {
	sink := &LogSink{
		logChan:        make(chan *v1.AccessLog, biz.MaxBatchCreateSize),
		autoCommitChan: make(chan *LogBatch, biz.MaxBatchCreateSize),
		conf:           conf.Log,
		au:             au,
		log:            log.NewHelper(log.With(logger, "module", "task/log")),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sink.loop(ctx.Done())
	}()

	return sink, cancel
}

func (s *LogSink) Append(log *v1.AccessLog) error {
	select {
	case s.logChan <- log:
	default:
		return v1.ErrorInvalidArgument("append unsupported log struct into sink")
	}
	return nil
}

func (s *LogSink) loop(stop <-chan struct{}) {
	var batch *LogBatch
	var commit *time.Timer

	for {
		select {
		case al := <-s.logChan:
			if batch == nil {
				batch = &LogBatch{}
				commit = time.AfterFunc(s.conf.BatchWriteInterval.AsDuration(), func() {
					s.autoCommitChan <- batch
				})
			}

			batch.Logs = append(batch.Logs, al)

			if len(batch.Logs) >= int(s.conf.BatchWriteSize) {
				if err := s.saveLogs(batch); err != nil {
					s.log.Error(err)
				}
				batch = nil
				commit.Stop()
			}

		case timeoutBatch := <-s.autoCommitChan:
			if timeoutBatch != batch {
				continue
			}
			if err := s.saveLogs(batch); err != nil {
				s.log.Error(err)
			}
			batch = nil

		case <-stop:
			if batch != nil {
				if err := s.saveLogs(batch); err != nil {
					s.log.Error(err)
				}
				batch = nil
			}
			if commit != nil {
				commit.Stop()
			}
		}
	}
}

func (s *LogSink) saveLogs(batch *LogBatch) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.conf.BatchWriteTimeout.AsDuration())
	defer cancel()

	if _, err := s.au.BatchCreate(ctx, batch.Logs); err != nil {
		return err
	}
	return nil
}
