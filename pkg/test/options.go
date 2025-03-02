/*
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

package test

import (
	"fmt"
	"time"

	"github.com/imdario/mergo"
	"github.com/samber/lo"

	"github.com/aws/karpenter-core/pkg/operator/options"
)

type OptionsFields struct {
	// Vendor Neutral
	ServiceName          *string
	DisableWebhook       *bool
	WebhookPort          *int
	MetricsPort          *int
	WebhookMetricsPort   *int
	HealthProbePort      *int
	KubeClientQPS        *int
	KubeClientBurst      *int
	EnableProfiling      *bool
	EnableLeaderElection *bool
	MemoryLimit          *int64
	LogLevel             *string
	BatchMaxDuration     *time.Duration
	BatchIdleDuration    *time.Duration
	FeatureGates         FeatureGates
}

type FeatureGates struct {
	Drift *bool
}

func Options(overrides ...OptionsFields) *options.Options {
	opts := OptionsFields{}
	for _, override := range overrides {
		if err := mergo.Merge(&opts, override, mergo.WithOverride); err != nil {
			panic(fmt.Sprintf("Failed to merge pod options: %s", err))
		}
	}

	return &options.Options{
		ServiceName:          lo.FromPtrOr(opts.ServiceName, ""),
		DisableWebhook:       lo.FromPtrOr(opts.DisableWebhook, false),
		WebhookPort:          lo.FromPtrOr(opts.WebhookPort, 8443),
		MetricsPort:          lo.FromPtrOr(opts.MetricsPort, 8000),
		WebhookMetricsPort:   lo.FromPtrOr(opts.WebhookMetricsPort, 8001),
		HealthProbePort:      lo.FromPtrOr(opts.HealthProbePort, 8081),
		KubeClientQPS:        lo.FromPtrOr(opts.KubeClientQPS, 200),
		KubeClientBurst:      lo.FromPtrOr(opts.KubeClientBurst, 300),
		EnableProfiling:      lo.FromPtrOr(opts.EnableProfiling, false),
		EnableLeaderElection: lo.FromPtrOr(opts.EnableLeaderElection, true),
		MemoryLimit:          lo.FromPtrOr(opts.MemoryLimit, -1),
		LogLevel:             lo.FromPtrOr(opts.LogLevel, ""),
		BatchMaxDuration:     lo.FromPtrOr(opts.BatchMaxDuration, 10*time.Second),
		BatchIdleDuration:    lo.FromPtrOr(opts.BatchIdleDuration, time.Second),
		FeatureGates: options.FeatureGates{
			Drift: lo.FromPtrOr(opts.FeatureGates.Drift, false),
		},
	}
}
