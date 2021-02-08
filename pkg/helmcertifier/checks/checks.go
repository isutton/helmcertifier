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
	"fmt"
	"helm.sh/helm/v3/pkg/lint"
	"path"
	"strings"

	"github.com/pkg/errors"
)

const (
	APIVersion2                  = "v2"
	ReadmeExist                  = "Chart has README"
	ReadmeDoesNotExist           = "Chart has not README"
	NotHelm3Reason               = "API version is not V2 used in Helm 3"
	Helm3Reason                  = "API version is V2 used in Helm 3"
	TestTemplatePrefix           = "templates/tests/"
	ChartTestFilesExist          = "Chart test files exist"
	ChartTestFilesDoesNotExist   = "Chart test files does not exist"
	MinKuberVersionSpecified     = "Minimum Kubernetes version specified"
	MinKuberVersionNotSpecified  = "Minimum Kubernetes version not specified"
	ValuesSchemaFileExist        = "Values schema file exist"
	ValuesSchemaFileDoesNotExist = "Values schema file does not exist"
	ValuesFileExist              = "Values file exist"
	ValuesFileDoesNotExist       = "Values file does not exist"
	ChartContainCRDs             = "Chart contain CRDs"
	ChartDoesNotContainCRDs      = "Chart does not contain CRDs"
	HelmLintSuccessful           = "Helm lint successful"
	HelmLintHasFailedPrefix      = "Helm lint has failed: "
)

func notImplemented() (Result, error) {
	return Result{Ok: false}, errors.New("not implemented")
}

func IsHelmV3(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}
	isHelmV3 := c.Metadata.APIVersion == APIVersion2

	reason := NotHelm3Reason
	if isHelmV3 {
		reason = Helm3Reason
	}
	return Result{Ok: isHelmV3, Reason: reason}, nil
}

func HasReadme(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := Result{Reason: ReadmeDoesNotExist}
	for _, f := range c.Files {
		if f.Name == "README.md" {
			r.Ok = true
			r.Reason = ReadmeExist
		}
	}

	return r, nil
}

func ContainsTest(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := Result{Reason: ChartTestFilesDoesNotExist}
	for _, f := range c.Templates {
		if strings.HasPrefix(f.Name, TestTemplatePrefix) && strings.HasSuffix(f.Name, ".yaml") {
			r.Reason = ChartTestFilesExist
			r.Ok = true
			break
		}
	}

	return r, nil

}

func ContainsValues(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := Result{Reason: ValuesFileDoesNotExist}

	if len(c.Values) > 0 {
		r.Reason = ValuesFileExist
		r.Ok = true
	}

	return r, nil
}

func ContainsValuesSchema(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := Result{Reason: ValuesSchemaFileDoesNotExist}

	if len(c.Schema) > 0 {
		r.Reason = ValuesSchemaFileExist
		r.Ok = true
	}

	return r, nil
}

func KeywordsAreOpenshiftCategories(uri string) (Result, error) {
	return notImplemented()
}

func IsCommercialChart(uri string) (Result, error) {
	return notImplemented()
}

func IsCommunityChart(uri string) (Result, error) {
	return notImplemented()
}

func HasMinKubeVersion(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := Result{Reason: MinKuberVersionNotSpecified}

	if c.Metadata.KubeVersion != "" {
		r.Ok = true
		r.Reason = MinKuberVersionSpecified
	}

	return r, nil
}

func NotContainCRDs(uri string) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := Result{Ok: true, Reason: ChartDoesNotContainCRDs}

	if len(c.CRDObjects()) > 0 {
		r.Ok = false
		r.Reason = ChartContainCRDs
	}

	return r, nil
}

func HelmLint(uri string) (Result, error) {
	c, p, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}
	r := Result{Ok: true, Reason: HelmLintSuccessful}
	p = path.Join(p, c.Name())
	linter := lint.All(p, map[string]interface{}{}, "default", false)
	if len(linter.Messages) > 0 {
		reason := ""
		for _, m := range linter.Messages {
			reason = reason + m.Error() + "\n"
		}
		r = Result{Ok: false, Reason: fmt.Sprintf("%s %s", HelmLintHasFailedPrefix, reason)}
	}
	return r, nil
}

func NotContainsInfraPluginsAndDrivers(uri string) (Result, error) {
	return notImplemented()
}

func CanBeInstalledWithoutManualPreRequisites(uri string) (Result, error) {
	return notImplemented()
}

func CanBeInstalledWithoutClusterAdminPrivileges(uri string) (Result, error) {
	return notImplemented()
}
