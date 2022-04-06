/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xyctruth/kubectl-replicas/pkg"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"os"
)

var (
	example = `
	# stash replicas of deployment
	$ kubectl replicas stash -n <namespace-name> 
	$ kubectl-replicas stash -n <namespace-name> 

	# recover replicas of deployment
	$ kubectl replicas recover -n <namespace-name> 
	$ kubectl-replicas recover -n <namespace-name> 
`
)

const (
	StashOp   = "stash"
	RecoverOp = "recover"
	Version   = "v0.0.1"
)

type ReplicasOptions struct {
	configFlags  *genericclioptions.ConfigFlags
	IOStreams    genericclioptions.IOStreams
	args         []string
	namespace    string
	op           string
	client       kubernetes.Interface
	printVersion bool
}

func NewReplicasOptions(streams genericclioptions.IOStreams) *ReplicasOptions {
	return &ReplicasOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func NewReplicasCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewReplicasOptions(streams)

	cmd := &cobra.Command{
		Use:          "replicas",
		Short:        "stash or recover replicas of deployment",
		Example:      example,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if o.printVersion {
				fmt.Println(Version)
				os.Exit(0)
			}
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				_ = c.Help()
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().BoolVar(&o.printVersion, "version", false, "prints version of plugin")
	o.configFlags.AddFlags(cmd.Flags())
	return cmd
}

func (o *ReplicasOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	o.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	o.namespace = getNamespace(o.configFlags)

	if len(o.args) > 0 {
		o.op = o.args[0]
	}
	return nil
}

func (o *ReplicasOptions) Validate() error {
	if len(o.args) != 1 {
		return fmt.Errorf("only one argument expected. got %d arguments", len(o.args))
	}

	if o.args[0] != StashOp && o.args[0] != RecoverOp {
		return fmt.Errorf("unknow arguments %s", o.args[0])
	}
	return nil
}

func (o *ReplicasOptions) Run() error {
	if o.op == StashOp {
		return pkg.Stash(o.client, o.namespace)
	}

	return pkg.Recover(o.client, o.namespace)
}

func getNamespace(flags *genericclioptions.ConfigFlags) string {
	namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()
	if err != nil || len(namespace) == 0 {
		namespace = "default"
	}
	return namespace
}
