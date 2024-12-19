package verificationcaching

import "sync/atomic"

type InMemoryMetrics struct {
	CredentialVerificationsSaved atomic.Int32
	FromDataVerifyTimeSpentMS    atomic.Int64
	ResultCacheHits              atomic.Int32
	ResultCacheHitsWasted        atomic.Int32
	ResultCacheMisses            atomic.Int32
}

var _ MetricsReporter = (*InMemoryMetrics)(nil)

func (m *InMemoryMetrics) AddCredentialVerificationsSaved(count int) {
	m.CredentialVerificationsSaved.Add(int32(count))
}

func (m *InMemoryMetrics) AddFromDataVerifyTimeSpent(ms int64) {
	m.FromDataVerifyTimeSpentMS.Add(ms)
}

func (m *InMemoryMetrics) AddResultCacheHits(count int) {
	m.ResultCacheHits.Add(int32(count))
}

func (m *InMemoryMetrics) AddResultCacheMisses(count int) {
	m.ResultCacheMisses.Add(int32(count))
}

func (m *InMemoryMetrics) AddResultCacheHitsWasted(count int) {
	m.ResultCacheHitsWasted.Add(int32(count))
}
