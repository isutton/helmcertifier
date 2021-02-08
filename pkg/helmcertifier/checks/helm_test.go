/*
 * Copyright 2021 Red Hat
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package checks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"helmcertifier/pkg/testutil"
)

func TestLoadChartFromURI(t *testing.T) {
	addr := "127.0.0.1:9876"

	type testCase struct {
		description string
		uri         string
	}

	positiveCases := []testCase{
		{
			uri:         "chart-0.1.0-v3.valid.tgz",
			description: "absolute path",
		},
		{
			uri:         "http://" + addr + "/charts/chart-0.1.0-v3.valid.tgz",
			description: "remote path, http",
		},
	}

	negativeCases := []testCase{
		{
			uri:         "chart-0.1.0-v3.non-existing.tgz",
			description: "non existing file",
		},
		{
			uri:         "http://" + addr + "/charts/chart-0.1.0-v3.non-existing.tgz",
			description: "non existing remote file",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go testutil.ServeCharts(ctx, addr, "./")

	for _, tc := range positiveCases {
		t.Run(tc.description, func(t *testing.T) {
			c, _, err := LoadChartFromURI(tc.uri)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}

	for _, tc := range negativeCases {
		t.Run(tc.description, func(t *testing.T) {
			c, _, err := LoadChartFromURI(tc.uri)
			require.Error(t, err)
			require.True(t, IsChartNotFound(err))
			require.Equal(t, "chart not found: "+tc.uri, err.Error())
			require.Nil(t, c)
		})
	}

	cancel()
}
