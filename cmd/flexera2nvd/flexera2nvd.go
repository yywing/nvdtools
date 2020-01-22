// Copyright (c) Facebook, Inc. and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/facebookincubator/nvdtools/providers/flexera/api"
	"github.com/facebookincubator/nvdtools/providers/flexera/schema"
	"github.com/facebookincubator/nvdtools/providers/lib/client"
	"github.com/facebookincubator/nvdtools/providers/lib/runner"
)

func Read(r io.Reader, c chan runner.Convertible) error {
	var vulns map[string]*schema.Advisory
	if err := json.NewDecoder(r).Decode(&vulns); err != nil {
		return fmt.Errorf("can't decode into vulns: %v", err)
	}

	for _, vuln := range vulns {
		c <- vuln
	}

	return nil
}

func FetchSince(ctx context.Context, c client.Client, baseURL string, since int64) (<-chan runner.Convertible, error) {
	apiKey := os.Getenv("FLEXERA_TOKEN")
	if apiKey == "" {
		return nil, fmt.Errorf("please set FLEXERA_TOKEN in environment")
	}

	client := api.NewClient(c, baseURL, apiKey)
	return client.FetchAllVulnerabilities(ctx, since)
}

func main() {
	r := runner.Runner{
		Config: runner.Config{
			BaseURL: "https://api.app.secunia.com",
			ClientConfig: client.Config{
				UserAgent: "flexera2nvd",
			},
		},
		FetchSince: FetchSince,
		Read:       Read,
	}

	if err := r.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
