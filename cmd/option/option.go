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

package option

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/yshngg/kle/pkg/client"
	fakeclient "github.com/yshngg/kle/pkg/client/fake"
	"github.com/yshngg/kle/pkg/leaderelection"
	"github.com/yshngg/kle/pkg/middleware"
	"k8s.io/apiserver/pkg/server/healthz"
	clientset "k8s.io/client-go/kubernetes"
	componentbaseconfig "k8s.io/component-base/config"
	componentbaseoptions "k8s.io/component-base/config/options"
	"k8s.io/klog/v2"
)

const ServerShutdownTimeout = 10 * time.Second

type KLEServer struct {
	Addr   string
	DryRun bool

	LeaderElection   componentbaseconfig.LeaderElectionConfiguration
	ClientConnection componentbaseconfig.ClientConnectionConfiguration
}

func NewKLEServer() *KLEServer {
	return &KLEServer{
		LeaderElection: *leaderelection.DefaultLeaderElectionConfig(),
	}
}

func (ks *KLEServer) Apply() error {
	return nil
}

// AddFlags adds flags for a specific KLEServer to the specified FlagSet
func (ks *KLEServer) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&ks.Addr, "addr", ks.Addr, "The address kel server binds to.")

	fs.BoolVar(&ks.DryRun, "dry-run", ks.DryRun, "Execute kle in dry run mode.")
	fs.StringVar(&ks.ClientConnection.Kubeconfig, "kubeconfig", ks.ClientConnection.Kubeconfig, "File with kube configuration. Deprecated, use client-connection-kubeconfig instead.")
	fs.StringVar(&ks.ClientConnection.Kubeconfig, "client-connection-kubeconfig", ks.ClientConnection.Kubeconfig, "File path to kube configuration for interacting with kubernetes apiserver.")
	fs.Float32Var(&ks.ClientConnection.QPS, "client-connection-qps", ks.ClientConnection.QPS, "QPS to use for interacting with kubernetes apiserver.")
	fs.Int32Var(&ks.ClientConnection.Burst, "client-connection-burst", ks.ClientConnection.Burst, "Burst to use for interacting with kubernetes apiserver.")

	componentbaseoptions.BindLeaderElectionFlags(&ks.LeaderElection, fs)
}

func (ks *KLEServer) Run(ctx context.Context) (err error) {
	ctx, done := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer done()

	mux := http.NewServeMux()
	http.Handle("/", mux)
	healthz.InstallHandler(mux)
	healthz.InstallLivezHandler(mux)

	go func() {
		err := httpServer(ctx, ks.Addr)
		if err != nil {
			klog.Errorf("http server, err: %v", err)
		}
	}()

	run := func() {
		healthz.InstallReadyzHandler(mux)
		run(ctx)
	}

	var kubeClient clientset.Interface
	if ks.DryRun {
		klog.Warning("dry run mode")
		kubeClient, err = fakeclient.Kubernetes()
		if err != nil {
			return fmt.Errorf("create kubernetes client, err: %w", err)
		}
	} else {
		kubeClient, err = client.Kubernetes(ks.ClientConnection)
		if err != nil {
			return fmt.Errorf("create kubernetes client, err: %w", err)
		}
	}

	if ks.LeaderElection.LeaderElect {
		if err = leaderelection.NewLeaderElection(run, kubeClient, &ks.LeaderElection, ctx); err != nil {
			return fmt.Errorf("create leader election, err: %w", err)
		}
		return nil
	}

	run()
	return nil
}

func httpServer(ctx context.Context, addr string) error {
	srv := http.Server{Addr: addr}
	serverErr := make(chan error, 1)
	go func() {
		// Capture ListenAndServe errors such as "port already in use".
		// However, when a server is gracefully shutdown, it is safe to ignore errors
		// returned from this method (given the select logic below), because
		// Shutdown causes ListenAndServe to always return http.ErrServerClosed.
		klog.Info("Starting http service...")
		if len(addr) == 0 {
			addr = ":80"
		}
		klog.Infof("Listening on %s", addr)
		serverErr <- srv.ListenAndServe()
	}()
	var err error
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownTimeout)
		defer cancel()
		klog.Info("Shutting down http service...")
		err = srv.Shutdown(ctx)
	case err = <-serverErr:
	}
	return err
}

func run(ctx context.Context) {
	pingCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ping_request_count",
			Help: "No of request handled by Ping handler",
		},
	)

	registry := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		pingCounter,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		pingCounter.Inc()
		fmt.Fprintf(w, "pong")
	})

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle(
		"/metrics",
		middleware.New(registry, nil).
			WrapHandler("/metrics", promhttp.HandlerFor(
				registry,
				promhttp.HandlerOpts{},
			)),
	)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			klog.Info("tick...")
		case <-ctx.Done():
			return
		}
	}
}
