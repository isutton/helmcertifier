/*
 * Copyright (C) 08/01/2021, 01:52, igors
 * This file is part of helmcertifier.
 *
 * helmcertifier is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * helmcertifier is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with helmcertifier.  If not, see <http://www.gnu.org/licenses/>.
 */

package checks

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// loadChartFromRemote attempts to retrieve a Helm chart from the given remote url. Returns an error if the given url
// doesn't contain the 'http' or 'https' schema, or any other error related to retrieving the contents of the chart.
func loadChartFromRemote(url *url.URL) (*chart.Chart, error) {
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, errors.Errorf("only 'http' and 'https' schemes are supported, but got %q", url.Scheme)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ChartNotFoundErr(url.String())
	}

	return loader.LoadArchive(resp.Body)
}

// loadChartFromAbsPath attempts to retrieve a local Helm chart by resolving the maybe relative path into an absolute
// path from the current working directory.
func loadChartFromAbsPath(path string) (*chart.Chart, error) {
	// although filepath.Abs() can return an error according to its signature, this won't happen (as of go 1.15)
	// because the only invalid value it would accept is an empty string, which is internally converted into "."
	// regardless, the error is still being caught and propagated to avoid being bitten by internal changes in the
	// future
	chartPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	c, err := loader.Load(chartPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ChartNotFoundErr(path)
		}
		return nil, err
	}

	return c, nil
}

type ChartCache interface {
	MakeKey(uri string) string
	Add(uri string, chrt *chart.Chart) error
	Get(uri string) (*chart.Chart, bool, error)
}

type chartCacheItem struct {
	Chart *chart.Chart
	Path  string
}

type chartCache struct {
	chartMap map[string]chartCacheItem
}

func newChartCache() *chartCache {
	return &chartCache{
		chartMap: make(map[string]chartCacheItem),
	}
}

func (c *chartCache) MakeKey(uri string) string {
	return regexp.MustCompile("[:/?]").ReplaceAllString(uri, "_")
}

func (c *chartCache) Get(uri string) (*chart.Chart, bool, error) {
	if chrt, ok := c.chartMap[c.MakeKey(uri)]; !ok {
		return nil, false, nil
	} else {
		return chrt.Chart, true, nil
	}
}

func (c *chartCache) Add(uri string, chrt *chart.Chart) error {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	key := c.MakeKey(uri)
	cacheDir := path.Join(userCacheDir, "helmcertifier")
	chartCacheDir := path.Join(cacheDir, key)
	c.chartMap[key] = chartCacheItem{
		Chart: chrt,
		Path:  chartCacheDir,
	}
	return nil
}

var defaultChartCache *chartCache

func init() {
	defaultChartCache = newChartCache()
}

// LoadChartFromURI attempts to retrieve a chart from the given uri string. It accepts "http", "https", "file" schemes,
// and defaults to "file" if there isn't one.
func LoadChartFromURI(uri string) (*chart.Chart, error) {
	var (
		chrt *chart.Chart
		err  error
	)

	if chrt, ok, _ := defaultChartCache.Get(uri); ok {
		return chrt, nil
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		chrt, err = loadChartFromRemote(u)
	case "file", "":
		chrt, err = loadChartFromAbsPath(u.Path)
	default:
		return nil, errors.Errorf("scheme %q not supported", u.Scheme)
	}

	if err != nil {
		return nil, err
	}

	if err = defaultChartCache.Add(uri, chrt); err != nil {
		return nil, err
	}

	return chrt, nil
}

type ChartNotFoundErr string

func (c ChartNotFoundErr) Error() string {
	return "chart not found: " + string(c)
}

func IsChartNotFound(err error) bool {
	_, ok := err.(ChartNotFoundErr)
	return ok
}
