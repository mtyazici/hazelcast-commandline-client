//go:build std || map

package _map

import (
	"context"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"

	"github.com/hazelcast/hazelcast-commandline-client/base"
	"github.com/hazelcast/hazelcast-commandline-client/base/commands"
	"github.com/hazelcast/hazelcast-commandline-client/clc"
	"github.com/hazelcast/hazelcast-commandline-client/clc/cmd"
	"github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/output"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
	"github.com/hazelcast/hazelcast-commandline-client/internal/proto/codec"
	"github.com/hazelcast/hazelcast-commandline-client/internal/proto/codec/control"
	"github.com/hazelcast/hazelcast-commandline-client/internal/serialization"
)

type MapGetEntryView struct{}

func (MapGetEntryView) Init(cc plug.InitContext) error {
	cc.SetCommandUsage("get-entry-view")
	long := `ss`
	short := "ss"
	cc.SetCommandHelp(long, short)
	commands.AddKeyTypeFlag(cc)
	cc.AddStringArg(commands.ArgKey, commands.ArgTitleKey)
	return nil
}

func (MapGetEntryView) Exec(ctx context.Context, ec plug.ExecContext) error {
	name := ec.Props().GetString(base.FlagName)
	keyStr := ec.GetStringArg(commands.ArgKey)
	var ci *hazelcast.ClientInternal
	entryv, stop, err := ec.ExecuteBlocking(ctx, func(ctx context.Context, sp clc.Spinner) (any, error) {
		var err error
		ci, err = cmd.ClientInternal(ctx, ec, sp)
		if err != nil {
			return nil, err
		}
		cmd.IncrementClusterMetric(ctx, ec, "total.map")
		if _, err = getMap(ctx, ec, sp); err != nil {
			return nil, err
		}
		keyData, err := commands.MakeKeyData(ec, ci, keyStr)
		if err != nil {
			return nil, err
		}
		sp.SetText(fmt.Sprintf("Getting map entry view from '%s'", name))
		req := codec.EncodeMapGetEntryViewRequest(name, keyData, 0)
		msg, err := ci.InvokeOnKey(ctx, req, keyData, nil)
		if err != nil {
			return nil, err
		}
		entry, _ := codec.DecodeMapGetEntryViewResponse(msg)
		return entry, nil
	})
	if err != nil {
		return err
	}
	stop()
	entry := entryv.(control.SimpleEntryView)
	showType := ec.Props().GetBool(base.FlagShowType)
	row := output.DecodePair(ci, hazelcast.NewPair(entry.Key, entry.Value), showType)
	if verbose := ec.Props().GetBool(clc.PropertyVerbose); verbose {
		row = append(row, output.Row{
			output.Column{
				Name:  "Cost",
				Type:  serialization.TypeInt64,
				Value: entry.Cost,
			},
			output.Column{
				Name:  "Creation Time",
				Type:  serialization.TypeInt64,
				Value: entry.CreationTime,
			},
			output.Column{
				Name:  "Expiration Time",
				Type:  serialization.TypeInt64,
				Value: entry.ExpirationTime,
			},
			output.Column{
				Name:  "Hits",
				Type:  serialization.TypeInt64,
				Value: entry.Hits,
			},
			output.Column{
				Name:  "Last Access Time",
				Type:  serialization.TypeInt64,
				Value: entry.LastAccessTime,
			},
			output.Column{
				Name:  "Last Stored Time",
				Type:  serialization.TypeInt64,
				Value: entry.LastStoredTime,
			},
			output.Column{
				Name:  "Last Update Time",
				Type:  serialization.TypeInt64,
				Value: entry.LastUpdateTime,
			},
			output.Column{
				Name:  "Version",
				Type:  serialization.TypeInt64,
				Value: entry.Version,
			},
			output.Column{
				Name:  "TTL",
				Type:  serialization.TypeInt64,
				Value: entry.TTL,
			},
			output.Column{
				Name:  "Max Idle",
				Type:  serialization.TypeInt64,
				Value: entry.MaxIdle,
			},
		}...)
	}
	return ec.AddOutputRows(ctx, row)
}

func init() {
	check.Must(plug.Registry.RegisterCommand("map:get-entry-view", &MapGetEntryView{}))
}
