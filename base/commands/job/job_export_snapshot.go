package job

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hazelcast/hazelcast-commandline-client/clc"
	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
	"github.com/hazelcast/hazelcast-commandline-client/internal/proto/codec"
)

type ExportSnapshotCmd struct{}

func (cm ExportSnapshotCmd) Init(cc plug.InitContext) error {
	cc.SetCommandUsage("export-snapshot [job-ID/name]")
	help := "Exports a snapshot for a job"
	cc.SetCommandHelp(help, help)
	cc.SetPositionalArgCount(1, 1)
	cc.AddStringFlag(flagName, "", "", false, "specify the snapshot. By default an auto-genertaed snapshot name is used")
	cc.AddBoolFlag(flagCancel, "", false, false, "cancel the job after taking the snapshot")
	return nil
}

func (cm ExportSnapshotCmd) Exec(ctx context.Context, ec plug.ExecContext) error {
	ci, err := ec.ClientInternal(ctx)
	if err != nil {
		return err
	}
	var jm *jobNameToIDMap
	var jid int64
	var ok bool
	jobNameOrID := ec.Args()[0]
	name := ec.Props().GetString(flagName)
	cancel := ec.Props().GetBool(flagCancel)
	if name == "" {
		// create the default snapshot name
		jm, err = newJobNameToIDMap(ctx, ec, true)
		if err != nil {
			return err
		}
		jid, ok = jm.GetIDForName(jobNameOrID)
		if !ok {
			return errInvalidJobID
		}
		name, ok = jm.GetNameForID(jid)
		if !ok {
			name = "UNKNOWN"
		}
		name = autoGenerateSnapshotName(name)
	} else {
		jm, err = newJobNameToIDMap(ctx, ec, false)
		if err != nil {
			return err
		}
		jid, ok = jm.GetIDForName(jobNameOrID)
		if !ok {
			return errInvalidJobID
		}
	}
	if err != nil {
		return err
	}
	_, stop, err := ec.ExecuteBlocking(ctx, func(ctx context.Context, sp clc.Spinner) (any, error) {
		sp.SetText(fmt.Sprintf("Exporting snapshot: %s", name))
		req := codec.EncodeJetExportSnapshotRequest(jid, name, cancel)
		if _, err := ci.InvokeOnRandomTarget(ctx, req, nil); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	stop()
	return nil
}

func autoGenerateSnapshotName(jobName string) string {
	dt := time.Now().UTC().Format("06-01-02_150405")
	r := rand.Int31n(10_000)
	return fmt.Sprintf("%s-%s-%d", jobName, dt, r)
}

func init() {
	Must(plug.Registry.RegisterCommand("job:export-snapshot", &ExportSnapshotCmd{}))
}
