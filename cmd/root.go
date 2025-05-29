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

package cmd

import (
	"flag"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yshngg/kle/cmd/option"
	"k8s.io/klog/v2"
)

func NewKLECommand(out io.Writer) *cobra.Command {
	s := option.NewKLEServer()
	cmd := &cobra.Command{
		Use:   "kle",
		Short: "A Kubernetes Leader Election Demo",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = s.Apply(); err != nil {
				klog.Errorf("apply kle, err: %v", err)
				return err
			}

			if err = s.Run(cmd.Context()); err != nil {
				klog.Errorf("run kle, err: %v", err)
				return err
			}
			return nil
		},
	}
	s.AddFlags(cmd.Flags())
	return cmd
}

func Execute() {
	out := os.Stdout
	cmd := NewKLECommand(out)
	cmd.AddCommand(NewVersionCommand())

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	cmd.Flags().AddGoFlagSet(klogFlags)

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
