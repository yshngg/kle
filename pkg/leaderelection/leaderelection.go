// The MIT License (MIT)
//
// Copyright © 2025 Yusheng Guo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package leaderelection

import (
	"context"
	"fmt"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"
)

// DefaultLeaderElectionConfig returns the default leader election configuration.
func DefaultLeaderElectionConfig() *componentbaseconfig.LeaderElectionConfiguration {
	return &componentbaseconfig.LeaderElectionConfiguration{
		LeaderElect:       false,
		LeaseDuration:     metav1.Duration{Duration: 15 * time.Second},
		RenewDeadline:     metav1.Duration{Duration: 10 * time.Second},
		RetryPeriod:       metav1.Duration{Duration: 2 * time.Second},
		ResourceLock:      "leases",
		ResourceName:      "kle",
		ResourceNamespace: "demo",
	}
}

// NewLeaderElection starts the leader election code loop
func NewLeaderElection(
	run func(),
	client clientset.Interface,
	LeaderElectionConfig *componentbaseconfig.LeaderElectionConfiguration,
	ctx context.Context,
) error {
	var id string

	if hostname, err := os.Hostname(); err != nil {
		// on errors, make sure we're unique
		id = string(uuid.NewUUID())
	} else {
		// add a uniquifier so that two processes on the same host don't accidentally both become active
		id = hostname + "_" + string(uuid.NewUUID())
	}

	klog.V(3).Infof("Assigned unique lease holder id: %s", id)

	if len(LeaderElectionConfig.ResourceNamespace) == 0 {
		return fmt.Errorf("namespace may not be empty")
	}

	if len(LeaderElectionConfig.ResourceName) == 0 {
		return fmt.Errorf("name may not be empty")
	}

	lock, err := resourcelock.New(
		LeaderElectionConfig.ResourceLock,
		LeaderElectionConfig.ResourceNamespace,
		LeaderElectionConfig.ResourceName,
		client.CoreV1(),
		client.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: id,
		},
	)
	if err != nil {
		return fmt.Errorf("create leader election lock, err: %v", err)
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   LeaderElectionConfig.LeaseDuration.Duration,
		RenewDeadline:   LeaderElectionConfig.RenewDeadline.Duration,
		RetryPeriod:     LeaderElectionConfig.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				klog.V(1).Info("Started leading")
				run()
			},
			OnStoppedLeading: func() {
				klog.V(1).Info("Leader lost")
			},
			OnNewLeader: func(identity string) {
				// Just got the lock
				if identity == id {
					return
				}
				klog.V(1).Infof("New leader elected: %v", identity)
			},
		},
	})
	return nil
}
