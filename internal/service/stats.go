package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// statsRanges: duração da janela e tamanho do balde (s) de cada período.
// Todas as durações são múltiplas do balde, então a série tem tamanho fixo.
var statsRanges = map[string]struct {
	d      time.Duration
	bucket int
}{
	"24h": {24 * time.Hour, 600},
	"7d":  {7 * 24 * time.Hour, 3600},
	"14d": {14 * 24 * time.Hour, 10800},
	"30d": {30 * 24 * time.Hour, 21600},
}

// Stats agrega os checks e incidentes do monitor na janela `rng` terminando em now.
// O balde corrente entra parcial; since e until ficam alinhados ao balde para os
// pontos não "andarem" a cada refresh.
func (s *MonitorService) Stats(ctx context.Context, id, rng string, now time.Time) (*domain.MonitorStats, error) {
	r, ok := statsRanges[rng]
	if !ok {
		return nil, fmt.Errorf("%w: range must be 24h, 7d, 14d or 30d", ErrValidation)
	}
	if _, err := s.GetByID(ctx, id); err != nil {
		return nil, err
	}

	now = now.UTC()
	step := time.Duration(r.bucket) * time.Second
	until := now.Truncate(step).Add(step)
	since := until.Add(-r.d)

	st := &domain.MonitorStats{
		Range:         rng,
		Since:         since,
		Until:         until,
		BucketSeconds: r.bucket,
		Series:        make([]domain.StatsBucket, int(r.d/step)),
		Down:          []domain.DownWindow{},
	}
	for i := range st.Series {
		st.Series[i].Start = since.Add(time.Duration(i) * step)
	}

	buckets, err := s.checks.Buckets(ctx, id, since, until, r.bucket)
	if err != nil {
		return nil, err
	}
	var successes int
	var sumLatency int64
	for _, b := range buckets {
		i := int(b.Start.Sub(since) / step)
		if i < 0 || i >= len(st.Series) {
			continue
		}
		st.Series[i] = b
		st.Checks += b.Checks
		successes += b.Successes
		sumLatency += b.SumLatencyMs
	}
	if st.Checks > 0 {
		up := 100 * float64(successes) / float64(st.Checks)
		st.Uptime = &up
	}
	if successes > 0 {
		avg := float64(sumLatency) / float64(successes)
		st.AvgLatencyMs = &avg
	}

	incidents, err := s.incidents.FindInWindow(ctx, id, since)
	if err != nil {
		return nil, err
	}
	// o fim real é now, não until (que fecha o balde corrente no futuro)
	for _, inc := range incidents {
		start, end := inc.StartedAt, now
		if inc.EndedAt != nil {
			end = *inc.EndedAt
		}
		start = maxTime(start, since)
		end = minTime(end, now)
		if !end.After(start) {
			continue
		}
		st.Down = append(st.Down, domain.DownWindow{Start: start, End: end})
		st.DowntimeSeconds += int64(end.Sub(start) / time.Second)
	}
	st.Incidents = len(st.Down)
	return st, nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
