// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"

	"github.com/moov-io/base/admin"
)

// adminPprofFromEnv returns admin pprof options from the environment.
// An unset HTTP_ADMIN_PPROF_ENABLED keeps historical pprof-on behavior (nil Opts.Pprof).
func adminPprofFromEnv() *admin.Pprof {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("HTTP_ADMIN_PPROF_ENABLED")))
	if raw == "" {
		return nil
	}

	p := &admin.Pprof{Enabled: false}
	switch raw {
	case "1", "true", "yes":
		p.Enabled = true
	}
	if secret := os.Getenv("HTTP_ADMIN_PPROF_SECRET"); secret != "" {
		p.Secret = secret
	}
	return p
}
