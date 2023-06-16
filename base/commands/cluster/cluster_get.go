//go:build base || cluster

package cluster

import (
	"context"
	"errors"

	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/output"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
	"github.com/hazelcast/hazelcast-commandline-client/internal/serialization"
)

type ClusterGetCmd struct{}

func (mc *ClusterGetCmd) Init(cc plug.InitContext) error {
	cc.SetPositionalArgCount(0, 0)
	help := "Get information about the cluster"
	cc.SetCommandHelp(help, help)
	cc.SetCommandUsage("get [flags]")
	return nil
}

func (mc *ClusterGetCmd) Exec(ctx context.Context, ec plug.ExecContext) error {
	ci, err := ec.ClientInternal(ctx)
	if err != nil {
		return err
	}
	uuid := ci.ClusterID()
	cl := ci.ClusterService().FailoverService().Current()
	if cl == nil {
		return errors.New("could not connect to the cluster")
	}
	rows := output.Row{
		output.Column{
			Name:  "Name",
			Value: cl.ClusterName,
			Type:  serialization.TypeString,
		},
		output.Column{
			Name:  "UUID",
			Value: uuid,
			Type:  serialization.TypeUUID,
		},
	}
	ec.AddOutputRows(ctx, rows)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Must(plug.Registry.RegisterCommand("cluster:get", &ClusterGetCmd{}))
}