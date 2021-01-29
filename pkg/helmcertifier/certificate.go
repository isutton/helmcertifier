/*
 * Copyright (C) 29/12/2020, 15:35 igors
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

package helmcertifier

import (
	"encoding/json"
)

type chartMetadata struct {
	Name    string
	Version string
}

type metadata struct {
	ChartMetadata chartMetadata
}

func newMetadata(name, version string) *metadata {
	return &metadata{
		ChartMetadata: chartMetadata{
			Name:    name,
			Version: version,
		},
	}
}

type certificate struct {
	Ok             bool
	Metadata       *metadata
	CheckResultMap checkResultMap
}

type checkResultMap map[string]checkResult

type checkResult struct {
	Ok     bool   `json:"ok"`
	Reason string `json:"reason"`
}

func (c *certificate) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"metadata": map[string]interface{}{
			"chart": map[string]interface{}{
				"name":    c.Metadata.ChartMetadata.Name,
				"version": c.Metadata.ChartMetadata.Version,
			},
			"results": c.CheckResultMap,
		},
	}
	return json.Marshal(m)
}

func newCertificate(name, version string, ok bool, resultMap checkResultMap) Certificate {
	return &certificate{
		Metadata:       newMetadata(name, version),
		Ok:             ok,
		CheckResultMap: resultMap,
	}
}

func (c *certificate) IsOk() bool {
	return c.Ok
}

func (c *certificate) String() string {
	return "<CERTIFICATION OUTPUT>"
}
