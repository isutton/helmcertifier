/*
 * Copyright (C) 04/01/2021, 06:58, igors
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
	"strings"

	"github.com/pkg/errors"
)

const APIVersion2 = "v2"
const NotHelm3Reason = "API version is not V2 used in Helm 3"
const Helm3Reason = "API version is V2 used in Helm 3"
const TestTemplatePrefix = "templates/tests/"
const ChartTestFilesExists = "Chart test files exists"
const ChartTestFilesDoesNotExists = "Chart test files does not exists"

func notImplemented() (Result, error) {
	return Result{Ok: false}, errors.New("not implemented")
}

func IsHelmV3(uri string) (Result, error) {
	c, err := loadChartFromURI(uri)
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
	c, err := loadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	hasReadme := false
	for _, f := range c.Files {
		if f.Name == "README.md" {
			hasReadme = true
			break
		}
	}

	return Result{Ok: hasReadme}, nil
}

func ContainsTest(uri string) (Result, error) {
	c, err := loadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	containTest := false
	reason := ChartTestFilesDoesNotExists
	for _, f := range c.Templates {
		if strings.HasPrefix(f.Name, TestTemplatePrefix) && strings.HasSuffix(f.Name, ".yaml") {
			containTest = true
			reason = ChartTestFilesExists
			break
		}
	}

	return Result{Ok: containTest, Reason: reason}, nil

}

func ReadmeContainsValuesSchema(uri string) (Result, error) {
	return notImplemented()
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
	return notImplemented()
}

func NotContainsCRDs(uri string) (Result, error) {
	return notImplemented()
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
