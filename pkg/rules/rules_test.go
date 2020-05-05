// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package rules

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/thanos-io/thanos/pkg/rules/rulespb"
	"github.com/thanos-io/thanos/pkg/store/storepb"
	"github.com/thanos-io/thanos/pkg/testutil"
)

// testRulesAgainstExamples tests against alerts.yaml and rules.yaml examples.
func testRulesAgainstExamples(t *testing.T, dir string, server rulespb.RulesServer) {
	t.Helper()

	defer leaktest.CheckTimeout(t, 10*time.Second)()

	// We don't test internals, just if groups are expected.
	someAlert := &rulespb.Rule{Result: &rulespb.Rule_Alert{Alert: &rulespb.Alert{Name: "some"}}}
	someRecording := &rulespb.Rule{Result: &rulespb.Rule_Recording{Recording: &rulespb.RecordingRule{Name: "some"}}}

	alerts := []*rulespb.RuleGroup{
		{
			Name:                              "thanos-bucket-replicate.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-compact.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-component-absent.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert, someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-query.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert, someAlert, someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-receive.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-rule.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert, someAlert, someAlert, someAlert, someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-sidecar.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-store.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{someAlert, someAlert, someAlert, someAlert},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-bucket-replicate.rules",
			File:                              filepath.Join(dir, "alerts.yaml"),
			Rules:                             []*rulespb.Rule{},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
	}
	rules := []*rulespb.RuleGroup{
		{
			Name:                              "thanos-query.rules",
			File:                              filepath.Join(dir, "rules.yaml"),
			Rules:                             []*rulespb.Rule{someRecording, someRecording, someRecording, someRecording, someRecording},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-receive.rules",
			File:                              filepath.Join(dir, "rules.yaml"),
			Rules:                             []*rulespb.Rule{someRecording, someRecording, someRecording, someRecording, someRecording, someRecording},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
		{
			Name:                              "thanos-store.rules",
			File:                              filepath.Join(dir, "rules.yaml"),
			Rules:                             []*rulespb.Rule{someRecording, someRecording, someRecording, someRecording},
			Interval:                          60,
			DeprecatedPartialResponseStrategy: storepb.PartialResponseStrategy_WARN, PartialResponseStrategy: storepb.PartialResponseStrategy_WARN,
		},
	}

	for _, tcase := range []struct {
		requestedType rulespb.RulesRequest_Type

		expected    []*rulespb.RuleGroup
		expectedErr error
	}{
		{
			requestedType: rulespb.RulesRequest_ALL,
			expected:      append(append([]*rulespb.RuleGroup{}, alerts...), rules...),
		},
		{
			requestedType: rulespb.RulesRequest_ALERTING,
			expected:      append([]*rulespb.RuleGroup{}, alerts...),
		},
		{
			requestedType: rulespb.RulesRequest_RECORDING,
			expected:      append([]*rulespb.RuleGroup{}, rules...),
		},
	} {
		t.Run("", func(t *testing.T) {
			srv := &rulesServer{ctx: context.Background()}
			err := server.Rules(&rulespb.RulesRequest{}, srv)
			if tcase.expectedErr != nil {
				testutil.NotOk(t, err)
				testutil.Equals(t, tcase.expectedErr.Error(), err.Error())
				return
			}

			// We don't want to be picky, just check what number and types of rules within group are.
			got := srv.groups
			for i, g := range got {
				for j, r := range g.Rules {
					if r.GetAlert() != nil {
						got[i].Rules[j] = someAlert
						continue
					}
					if r.GetRecording() != nil {
						got[i].Rules[j] = someRecording
						continue
					}
					t.Fatalf("Found rule in group %s that is neither recording not alert.", g.Name)
				}
			}

			testutil.Ok(t, err)
			testutil.Equals(t, []error(nil), srv.warnings)
			testutil.Equals(t, tcase.expected, srv.groups)
		})
	}
}
