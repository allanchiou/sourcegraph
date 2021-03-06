package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sourcegraph/sourcegraph/internal/logging"
	"github.com/sourcegraph/sourcegraph/internal/metrics"
	"github.com/sourcegraph/sourcegraph/internal/trace"
)

// ObservedSource returns a decorator that wraps a Source
// with error logging, Prometheus metrics and tracing.
func ObservedSource(l logging.ErrorLogger, m SourceMetrics) func(Source) Source {
	return func(s Source) Source {
		return &observedSource{
			Source:  s,
			metrics: m,
			log:     l,
		}
	}
}

// An observedSource wraps another Source with error logging,
// Prometheus metrics and tracing.
type observedSource struct {
	Source
	metrics SourceMetrics
	log     logging.ErrorLogger
}

// SourceMetrics encapsulates the Prometheus metrics of a Source.
type SourceMetrics struct {
	ListRepos *metrics.OperationMetrics
}

// MustRegister registers all metrics in SourceMetrics in the given
// prometheus.Registerer. It panics in case of failure.
func (sm SourceMetrics) MustRegister(r prometheus.Registerer) {
	r.MustRegister(sm.ListRepos.Count)
	r.MustRegister(sm.ListRepos.Duration)
	r.MustRegister(sm.ListRepos.Errors)
}

// NewSourceMetrics returns SourceMetrics that need to be registered
// in a Prometheus registry.
func NewSourceMetrics() SourceMetrics {
	return SourceMetrics{
		ListRepos: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_source_duration_seconds",
				Help: "Time spent sourcing repos",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_source_repos_total",
				Help: "Total number of sourced repositories",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_source_errors_total",
				Help: "Total number of sourcing errors",
			}, []string{}),
		},
	}
}

// ListRepos calls into the inner Source registers the observed results.
func (o *observedSource) ListRepos(ctx context.Context, results chan SourceResult) {
	var (
		err   error
		count float64
	)

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		o.metrics.ListRepos.Observe(secs, count, &err)
		logging.Log(o.log, "source.list-repos", &err)
	}(time.Now())

	uncounted := make(chan SourceResult)
	go func() {
		o.Source.ListRepos(ctx, uncounted)
		close(uncounted)
	}()

	var errs *multierror.Error
	for res := range uncounted {
		results <- res
		if res.Err != nil {
			errs = multierror.Append(errs, res.Err)
		}
		count++
	}
}

// NewObservedStore wraps the given Store with error logging,
// Prometheus metrics and tracing.
func NewObservedStore(
	s Store,
	l logging.ErrorLogger,
	m StoreMetrics,
	t trace.Tracer,
) *ObservedStore {
	return &ObservedStore{
		store:   s,
		log:     l,
		metrics: m,
		tracer:  t,
	}
}

// An ObservedStore wraps another Store with error logging,
// Prometheus metrics and tracing.
type ObservedStore struct {
	store   Store
	log     logging.ErrorLogger
	metrics StoreMetrics
	tracer  trace.Tracer
	txtrace *trace.Trace
	txctx   context.Context
}

// StoreMetrics encapsulates the Prometheus metrics of a Store.
type StoreMetrics struct {
	Transact               *metrics.OperationMetrics
	Done                   *metrics.OperationMetrics
	UpsertRepos            *metrics.OperationMetrics
	ListRepos              *metrics.OperationMetrics
	UpsertExternalServices *metrics.OperationMetrics
	ListExternalServices   *metrics.OperationMetrics
	SetClonedRepos         *metrics.OperationMetrics
	CountNotClonedRepos    *metrics.OperationMetrics
}

// MustRegister registers all metrics in StoreMetrics in the given
// prometheus.Registerer. It panics in case of failure.
func (sm StoreMetrics) MustRegister(r prometheus.Registerer) {
	for _, om := range []*metrics.OperationMetrics{
		sm.Transact,
		sm.Done,
		sm.ListRepos,
		sm.UpsertRepos,
		sm.ListExternalServices,
		sm.UpsertExternalServices,
		sm.SetClonedRepos,
	} {
		r.MustRegister(om.Count)
		r.MustRegister(om.Duration)
		r.MustRegister(om.Errors)
	}
}

// NewStoreMetrics returns StoreMetrics that need to be registered
// in a Prometheus registry.
func NewStoreMetrics() StoreMetrics {
	return StoreMetrics{
		Transact: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_transact_duration_seconds",
				Help: "Time spent opening a transaction",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_transact_total",
				Help: "Total number of opened transactions",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_transact_errors_total",
				Help: "Total number of errors when opening a transaction",
			}, []string{}),
		},
		Done: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_done_duration_seconds",
				Help: "Time spent closing a transaction",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_done_total",
				Help: "Total number of closed transactions",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_done_errors_total",
				Help: "Total number of errors when closing a transaction",
			}, []string{}),
		},
		UpsertRepos: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_upsert_repos_duration_seconds",
				Help: "Time spent upserting repos",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_upsert_repos_total",
				Help: "Total number of upserted repositories",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_upsert_repos_errors_total",
				Help: "Total number of errors when upserting repos",
			}, []string{}),
		},
		ListRepos: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_list_repos_duration_seconds",
				Help: "Time spent listing repos",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_list_repos_total",
				Help: "Total number of listed repositories",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_list_repos_errors_total",
				Help: "Total number of errors when listing repos",
			}, []string{}),
		},
		UpsertExternalServices: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_external_serviceupdater_store_upsert_external_services_duration_seconds",
				Help: "Time spent upserting external_services",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_external_serviceupdater_store_upsert_external_services_total",
				Help: "Total number of upserted external_servicesitories",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_external_serviceupdater_store_upsert_external_services_errors_total",
				Help: "Total number of errors when upserting external_services",
			}, []string{}),
		},
		ListExternalServices: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_list_external_services_duration_seconds",
				Help: "Time spent listing external_services",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_list_external_services_total",
				Help: "Total number of listed external_servicesitories",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_list_external_services_errors_total",
				Help: "Total number of errors when listing external_services",
			}, []string{}),
		},
		SetClonedRepos: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_set_cloned_repos_duration_seconds",
				Help: "Time spent setting cloned repos",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_set_cloned_repos_total",
				Help: "Total number of set cloned repos calls",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_set_cloned_repos_errors_total",
				Help: "Total number of errors when setting cloned repos",
			}, []string{}),
		},
		CountNotClonedRepos: &metrics.OperationMetrics{
			Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "src_repoupdater_store_count_not_cloned_repos_duration_seconds",
				Help: "Time spent counting not-cloned repos",
			}, []string{}),
			Count: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_count_not_cloned_repos_total",
				Help: "Total number of count not-cloned repos calls",
			}, []string{}),
			Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "src_repoupdater_store_count_not_cloned_repos_errors_total",
				Help: "Total number of errors when counting not-cloned repos",
			}, []string{}),
		},
	}
}

// Transact calls into the inner Store Transact method and
// returns an observed TxStore.
func (o *ObservedStore) Transact(ctx context.Context) (s TxStore, err error) {
	tr, ctx := o.trace(ctx, "Store.Transact")

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		o.metrics.Transact.Observe(secs, 1, &err)
		logging.Log(o.log, "store.transact", &err)
		if err != nil {
			tr.SetError(err)
			// Finish is called in Done in the non-error case
			tr.Finish()
		}
	}(time.Now())

	s, err = o.store.(Transactor).Transact(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "observed store")
	}

	return &ObservedStore{
		store:   s,
		log:     o.log,
		metrics: o.metrics,
		tracer:  o.tracer,
		txtrace: tr,
		txctx:   ctx,
	}, nil
}

// Done calls into the inner Store Done method.
func (o *ObservedStore) Done(errs ...*error) {
	tr := o.txtrace
	tr.LogFields(otlog.String("event", "Store.Done"))

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		done := false

		for _, err := range errs {
			if err != nil && *err != nil {
				done = true
				tr.SetError(*err)
				o.metrics.Done.Observe(secs, 1, err)
				logging.Log(o.log, "store.done", err)
			}
		}

		if !done {
			o.metrics.Done.Observe(secs, 1, nil)
		}

		tr.Finish()
	}(time.Now())

	o.store.(TxStore).Done(errs...)
}

// ListExternalServices calls into the inner Store and registers the observed results.
func (o *ObservedStore) ListExternalServices(ctx context.Context, args StoreListExternalServicesArgs) (es []*ExternalService, err error) {
	tr, ctx := o.trace(ctx, "Store.ListExternalServices")
	tr.LogFields(
		otlog.Object("args.ids", args.IDs),
		otlog.Object("args.kinds", args.Kinds),
	)

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		count := float64(len(es))

		o.metrics.ListExternalServices.Observe(secs, count, &err)
		logging.Log(o.log, "store.list-external-services", &err,
			"args", fmt.Sprintf("%+v", args),
			"count", len(es),
		)

		tr.LogFields(
			otlog.Int("count", len(es)),
			otlog.Object("names", ExternalServices(es).DisplayNames()),
			otlog.Object("urns", ExternalServices(es).URNs()),
		)
		tr.SetError(err)
		tr.Finish()
	}(time.Now())

	return o.store.ListExternalServices(ctx, args)
}

// UpsertExternalServices calls into the inner Store and registers the observed results.
func (o *ObservedStore) UpsertExternalServices(ctx context.Context, svcs ...*ExternalService) (err error) {
	tr, ctx := o.trace(ctx, "Store.UpsertExternalServices")
	tr.LogFields(
		otlog.Int("count", len(svcs)),
		otlog.Object("names", ExternalServices(svcs).DisplayNames()),
		otlog.Object("urns", ExternalServices(svcs).URNs()),
	)

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		count := float64(len(svcs))

		o.metrics.UpsertExternalServices.Observe(secs, count, &err)
		logging.Log(o.log, "store.upsert-external-services", &err,
			"count", len(svcs),
			"names", ExternalServices(svcs).DisplayNames(),
		)

		tr.SetError(err)
		tr.Finish()
	}(time.Now())

	return o.store.UpsertExternalServices(ctx, svcs...)
}

// ListRepos calls into the inner Store and registers the observed results.
func (o *ObservedStore) ListRepos(ctx context.Context, args StoreListReposArgs) (rs []*Repo, err error) {
	tr, ctx := o.trace(ctx, "Store.ListRepos")
	tr.LogFields(
		otlog.Object("args.names", args.Names),
		otlog.Object("args.ids", args.IDs),
		otlog.Object("args.kinds", args.Kinds),
	)

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		count := float64(len(rs))

		o.metrics.ListRepos.Observe(secs, count, &err)
		logging.Log(o.log, "store.list-repos", &err,
			"args", fmt.Sprintf("%+v", args),
			"count", len(rs),
		)

		tr.LogFields(otlog.Int("count", len(rs)))
		tr.SetError(err)
		tr.Finish()
	}(time.Now())

	return o.store.ListRepos(ctx, args)
}

// UpsertRepos calls into the inner Store and registers the observed results.
func (o *ObservedStore) UpsertRepos(ctx context.Context, repos ...*Repo) (err error) {
	tr, ctx := o.trace(ctx, "Store.UpsertRepos")
	tr.LogFields(otlog.Int("count", len(repos)))

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		count := float64(len(repos))

		o.metrics.UpsertRepos.Observe(secs, count, &err)
		logging.Log(o.log, "store.upsert-repos", &err, "count", len(repos))

		tr.SetError(err)
		tr.Finish()
	}(time.Now())

	return o.store.UpsertRepos(ctx, repos...)
}

// SetClonedRepos calls into the inner Store and registers the observed results.
func (o *ObservedStore) SetClonedRepos(ctx context.Context, repoNames ...string) (err error) {
	tr, ctx := o.trace(ctx, "Store.SetClonedRepos")
	tr.LogFields(otlog.Int("count", len(repoNames)))

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()
		count := float64(len(repoNames))

		o.metrics.SetClonedRepos.Observe(secs, count, &err)
		logging.Log(o.log, "store.set-cloned-repos", &err, "count", len(repoNames))

		tr.SetError(err)
		tr.Finish()
	}(time.Now())

	return o.store.SetClonedRepos(ctx, repoNames...)
}

// CountNotClonedRepos calls into the inner Store and registers the observed results.
func (o *ObservedStore) CountNotClonedRepos(ctx context.Context) (count uint64, err error) {
	tr, ctx := o.trace(ctx, "Store.CountNotClonedRepos")

	defer func(began time.Time) {
		secs := time.Since(began).Seconds()

		o.metrics.CountNotClonedRepos.Observe(secs, float64(count), &err)
		logging.Log(o.log, "store.count-not-cloned-repos", &err, "count", count)

		tr.SetError(err)
		tr.Finish()
	}(time.Now())

	return o.store.CountNotClonedRepos(ctx)
}

func (o *ObservedStore) trace(ctx context.Context, family string) (*trace.Trace, context.Context) {
	txctx := o.txctx
	if txctx == nil {
		txctx = ctx
	}
	tr, _ := o.tracer.New(txctx, family, "")
	return tr, trace.ContextWithTrace(ctx, tr)
}
