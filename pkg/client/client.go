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

package client

import (
	"errors"
	"fmt"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	componentbaseconfig "k8s.io/component-base/config"
)

func Kubernetes(clientConnection componentbaseconfig.ClientConnectionConfiguration) (clientset.Interface, error) {
	cfg, err := createConfig(clientConnection)
	if err != nil {
		return nil, fmt.Errorf("unable to create config: %v", err)
	}

	return clientset.NewForConfig(cfg)
}

func createConfig(clientConnection componentbaseconfig.ClientConnectionConfiguration) (*rest.Config, error) {
	var cfg *rest.Config
	if len(clientConnection.Kubeconfig) != 0 {
		master, err := getMasterFromKubeconfig(clientConnection.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig file: %v ", err)
		}

		cfg, err = clientcmd.BuildConfigFromFlags(master, clientConnection.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("unable to build config: %v", err)
		}
	} else {
		var err error
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to build in cluster config: %v", err)
		}
	}

	cfg.Burst = int(clientConnection.Burst)
	cfg.QPS = clientConnection.QPS

	return cfg, nil
}

func getMasterFromKubeconfig(filename string) (string, error) {
	config, err := clientcmd.LoadFromFile(filename)
	if err != nil {
		return "", err
	}

	context, ok := config.Contexts[config.CurrentContext]
	if !ok {
		return "", fmt.Errorf("get master address from kubeconfig, err: %w", errors.New("current context not found"))
	}

	val, ok := config.Clusters[context.Cluster]
	if !ok {
		return "", fmt.Errorf("get master address from kubeconfig, err: %w", errors.New("cluster information not found"))
	}
	return val.Server, nil
}
