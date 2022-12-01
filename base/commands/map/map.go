//go:build base || map

package _map

import (
	"context"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"

	"github.com/hazelcast/hazelcast-commandline-client/clc"
	"github.com/hazelcast/hazelcast-commandline-client/clc/paths"
	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
)

const (
	mapFlagName     = "name"
	mapFlagShowType = "show-type"
	mapPropertyName = "map"
)

type MapCommand struct {
	m *hazelcast.Map
}

func (mc *MapCommand) Init(cc plug.InitContext) error {
	cc.SetCommandGroup(clc.GroupDDSID)
	cc.AddStringFlag(mapFlagName, "n", defaultMapName, false, "Map name")
	cc.AddBoolFlag(mapFlagShowType, "", false, false, "add the type names to the output")
	if !cc.Interactive() {
		cc.AddStringFlag(clc.PropertySchemaDir, "", paths.Schemas(), false, "set the schema directory")
	}
	cc.SetTopLevel(true)
	cc.SetCommandUsage("map COMMAND [flags]")
	help := "Map operations"
	cc.SetCommandHelp(help, help)
	return nil
}

func (mc *MapCommand) Exec(context.Context, plug.ExecContext) error {
	return nil
}

func (mc *MapCommand) Augment(ec plug.ExecContext, props *plug.Properties) error {
	ctx := context.TODO()
	props.SetBlocking(mapPropertyName, func() (any, error) {
		if mc.m != nil {
			return mc.m, nil
		}
		mapName := ec.Props().GetString(mapFlagName)
		// empty map name is allowed
		ci, err := ec.ClientInternal(ctx)
		if err != nil {
			return nil, err
		}
		hint := fmt.Sprintf("Getting map %s", mapName)
		mv, stop, err := ec.ExecuteBlocking(ctx, hint, func(ctx context.Context) (any, error) {
			m, err := ci.Client().GetMap(ctx, mapName)
			if err != nil {
				return nil, err
			}
			return m, nil
		})
		if err != nil {
			return nil, err
		}
		stop()
		mc.m = mv.(*hazelcast.Map)
		return mc.m, nil
	})
	return nil
}

func init() {
	cmd := &MapCommand{}
	Must(plug.Registry.RegisterCommand("map", cmd))
	plug.Registry.RegisterAugmentor("20-map", cmd)
}