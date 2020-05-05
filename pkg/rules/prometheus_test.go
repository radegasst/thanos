// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package rules

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/thanos-io/thanos/pkg/promclient"
	"github.com/thanos-io/thanos/pkg/testutil"
	"github.com/thanos-io/thanos/pkg/testutil/e2eutil"
)

func TestPrometheus_Rules_e2e(t *testing.T) {
	defer leaktest.CheckTimeout(t, 10*time.Second)()

	p, err := e2eutil.NewPrometheus()
	testutil.Ok(t, err)
	defer func() { testutil.Ok(t, p.Stop()) }()

	curr, err := os.Getwd()
	testutil.Ok(t, err)
	testutil.Ok(t, p.SetConfig(fmt.Sprintf(`
global:
  external_labels:
    region: eu-west

rule_files:
  - %s/../../examples/alerts/alerts.yaml
  - %s/../../examples/alerts/rules.yaml
`, curr, curr)))
	testutil.Ok(t, p.Start())

	u, err := url.Parse(fmt.Sprintf("http://%s", p.Addr()))
	testutil.Ok(t, err)

	promRules := NewPrometheus(u, promclient.NewDefaultClient())
	testRulesAgainstExamples(t, filepath.Join(curr, "examples/alerts"), promRules)
}
