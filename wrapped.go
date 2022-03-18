package logging

import (
	"go.uber.org/zap/zapcore"
)

type zapWrappedSyncer struct {
	zapcore.WriteSyncer
}

func (r *zapWrappedSyncer) SetOnClose(closeFunc func()) {}
func (r *zapWrappedSyncer) Rotate() (err error)         { return }
func (r *zapWrappedSyncer) Close() (err error)          { return }
func (r *zapWrappedSyncer) Sync() error                 { return r.WriteSyncer.Sync() }
