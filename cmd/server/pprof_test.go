// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/moov-io/base/admin"
	"github.com/stretchr/testify/require"
)

func TestAdminPprofFromEnv(t *testing.T) {
	t.Setenv("HTTP_ADMIN_PPROF_ENABLED", "")
	t.Setenv("HTTP_ADMIN_PPROF_SECRET", "")

	got := adminPprofFromEnv()
	require.Nil(t, got)

	t.Setenv("HTTP_ADMIN_PPROF_ENABLED", "false")
	got = adminPprofFromEnv()
	require.NotNil(t, got)
	require.False(t, got.Enabled)

	t.Setenv("HTTP_ADMIN_PPROF_ENABLED", "true")
	t.Setenv("HTTP_ADMIN_PPROF_SECRET", "test-pprof-secret")
	got = adminPprofFromEnv()
	require.True(t, got.Enabled)
	require.Equal(t, "test-pprof-secret", got.Secret)
}

func TestAdminServer_PprofExposure(t *testing.T) {
	tests := []struct {
		name    string
		enabled string
		status  int
	}{
		{name: "omitted env keeps historical pprof on", enabled: "", status: http.StatusOK},
		{name: "explicitly disabled", enabled: "false", status: http.StatusNotFound},
		{name: "explicitly enabled", enabled: "true", status: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HTTP_ADMIN_PPROF_ENABLED", tc.enabled)
			server, err := admin.New(admin.Opts{
				Addr:  "127.0.0.1:0",
				Pprof: adminPprofFromEnv(),
			})
			require.NoError(t, err)
			t.Cleanup(server.Shutdown)
			go func() { _ = server.Listen() }()

			for _, path := range []string{"/debug/pprof/", "/debug/pprof/cmdline"} {
				require.Equal(t, tc.status, adminHTTPStatus(t, server.BindAddr(), path), path)
			}

			require.Equal(t, http.StatusOK, adminHTTPStatus(t, server.BindAddr(), "/metrics"))
		})
	}
}

func adminHTTPStatus(t *testing.T, addr, path string) int {
	t.Helper()

	var last error
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err != nil {
			last = err
			time.Sleep(20 * time.Millisecond)
			continue
		}

		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
		_, err = fmt.Fprintf(conn, "GET %s HTTP/1.0\r\nHost: localhost\r\n\r\n", path)
		if err != nil {
			conn.Close()
			last = err
			time.Sleep(20 * time.Millisecond)
			continue
		}

		line, err := bufio.NewReader(conn).ReadString('\n')
		conn.Close()
		if err != nil {
			last = err
			time.Sleep(20 * time.Millisecond)
			continue
		}

		parts := strings.SplitN(strings.TrimSpace(line), " ", 3)
		require.GreaterOrEqual(t, len(parts), 2, line)
		code, err := strconv.Atoi(parts[1])
		require.NoError(t, err, line)
		return code
	}

	require.NoError(t, last)
	return 0
}
