/*
Copyright 2020, 2021 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package policy

import (
	"testing"
)

func TestNewAlphabetical(t *testing.T) {
	cases := []struct {
		label     string
		order     string
		expectErr bool
	}{
		{
			label: "With valid empty order",
			order: "",
		},
		{
			label: "With valid asc order",
			order: AlphabeticalOrderAsc,
		},
		{
			label: "With valid desc order",
			order: AlphabeticalOrderDesc,
		},
		{
			label:     "With invalid order",
			order:     "invalid",
			expectErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.label, func(t *testing.T) {
			_, err := NewAlphabetical(tt.order)
			if tt.expectErr && err == nil {
				t.Fatalf("expecting error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("returned unexpected error: %s", err)
			}
		})
	}
}

func TestAlphabetical_Latest(t *testing.T) {
	cases := []struct {
		label           string
		order           string
		versions        []string
		expectedVersion string
		expectErr       bool
	}{
		{
			label:           "With Ubuntu CalVer",
			versions:        []string{"16.04", "16.04.1", "16.10", "20.04", "20.10"},
			expectedVersion: "20.10",
		},
		{
			label:           "With Ubuntu CalVer descending",
			versions:        []string{"16.04", "16.04.1", "16.10", "20.04", "20.10"},
			order:           AlphabeticalOrderDesc,
			expectedVersion: "16.04",
		},
		{
			label:           "With Ubuntu code names",
			versions:        []string{"xenial", "yakkety", "zesty", "artful", "bionic"},
			expectedVersion: "zesty",
		},
		{
			label:           "With Ubuntu code names descending",
			versions:        []string{"xenial", "yakkety", "zesty", "artful", "bionic"},
			order:           AlphabeticalOrderDesc,
			expectedVersion: "artful",
		},
		{
			label:           "With Timestamps",
			versions:        []string{"1606234201", "1606364286", "1606334092", "1606334284", "1606334201"},
			expectedVersion: "1606364286",
		},
		{
			label:           "With Unix Timestamps desc",
			versions:        []string{"1606234201", "1606364286", "1606334092", "1606334284", "1606334201"},
			order:           AlphabeticalOrderDesc,
			expectedVersion: "1606234201",
		},
		{
			label:           "With Unix Timestamps prefix",
			versions:        []string{"rel-1606234201", "rel-1606364286", "rel-1606334092", "rel-1606334284", "rel-1606334201"},
			expectedVersion: "rel-1606364286",
		},
		{
			label:           "With RFC3339",
			versions:        []string{"2021-01-08T21-18-21Z", "2020-05-08T21-18-21Z", "2021-01-08T19-20-00Z", "1990-01-08T00-20-00Z", "2023-05-08T00-20-00Z"},
			expectedVersion: "2023-05-08T00-20-00Z",
		},
		{
			label:           "With RFC3339 desc",
			versions:        []string{"2021-01-08T21-18-21Z", "2020-05-08T21-18-21Z", "2021-01-08T19-20-00Z", "1990-01-08T00-20-00Z", "2023-05-08T00-20-00Z"},
			order:           AlphabeticalOrderDesc,
			expectedVersion: "1990-01-08T00-20-00Z",
		},
		{
			label:     "Empty version list",
			versions:  []string{},
			expectErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.label, func(t *testing.T) {
			policy, err := NewAlphabetical(tt.order)
			if err != nil {
				t.Fatalf("returned unexpected error: %s", err)
			}
			latest, err := policy.Latest(tt.versions)
			if tt.expectErr && err == nil {
				t.Fatalf("expecting error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("returned unexpected error: %s", err)
			}

			if latest != tt.expectedVersion {
				t.Errorf("incorrect computed version returned, got '%s', expected '%s'", latest, tt.expectedVersion)
			}
		})
	}
}
