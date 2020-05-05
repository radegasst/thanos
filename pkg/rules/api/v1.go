// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package v1

import (
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thanos-io/thanos/pkg/rules/rulespb"

	"github.com/go-kit/kit/log"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/route"
	extpromhttp "github.com/thanos-io/thanos/pkg/extprom/http"
	qapi "github.com/thanos-io/thanos/pkg/query/api"
	"github.com/thanos-io/thanos/pkg/rules/manager"
	"github.com/thanos-io/thanos/pkg/tracing"
)

type API struct {
	logger        log.Logger
	now           func() time.Time
	ruleRetriever RulesRetriever
	reg           prometheus.Registerer
}

func NewAPI(
	logger log.Logger,
	reg prometheus.Registerer,
	ruleRetriever RulesRetriever,
) *API {
	return &API{
		logger:        logger,
		now:           time.Now,
		ruleRetriever: ruleRetriever,
		reg:           reg,
	}
}

func (api *API) Register(r *route.Router, tracer opentracing.Tracer, logger log.Logger, ins extpromhttp.InstrumentationMiddleware) {
	instr := func(name string, f qapi.ApiFunc) http.HandlerFunc {
		hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			qapi.SetCORS(w)
			if data, warnings, err := f(r); err != nil {
				qapi.RespondError(w, err, data)
			} else if data != nil {
				qapi.Respond(w, data, warnings)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		})
		return ins.NewHandler(name, tracing.HTTPMiddleware(tracer, name, logger, gziphandler.GzipHandler(hf)))
	}

	r.Get("/alerts", instr("alerts", api.alerts))
	r.Get("/rules", instr("rules", api.rules))
}

type RulesRetriever interface {
	RuleGroups() []manager.Group
	AlertingRules() []manager.AlertingRule
}

func (api *API) rules(*http.Request) (interface{}, []error, *qapi.ApiError) {
	res := &rulespb.RuleGroups{}
	for _, grp := range api.ruleRetriever.RuleGroups() {
		res.Groups = append(res.Groups, grp.ToProto())
	}
	return res, nil, nil
}

func (api *API) alerts(*http.Request) (interface{}, []error, *qapi.ApiError) {
	var alerts []*rulespb.AlertInstance
	for _, alertingRule := range api.ruleRetriever.AlertingRules() {
		alerts = append(alerts, alertingRule.ActiveAlertsToProto()...)
	}
	return struct{ Alerts []*rulespb.AlertInstance }{Alerts: alerts}, nil, nil
}
