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

type chartMetadata struct {
	Name    string
	Version string
}

type certificateMetadata struct {
	ChartMetadata chartMetadata
}

func newCertificateMetadata(name, version string) *certificateMetadata {
	return &certificateMetadata{
		ChartMetadata: chartMetadata{
			Name:    name,
			Version: version,
		},
	}
}

type certificate struct {
	CertificateMetadata *certificateMetadata
	Ok                  bool
}

func newCertificate(name, version string, ok bool) Certificate {
	return &certificate{
		CertificateMetadata: newCertificateMetadata(name, version),
		Ok:                  ok,
	}
}

func (r *certificate) IsOk() bool {
	return r.Ok
}

func (r *certificate) String() string {
	return "<CERTIFICATION OUTPUT>"
}
