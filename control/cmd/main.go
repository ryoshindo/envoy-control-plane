package main

import (
	"context"
	"flag"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"github.com/ryoshindo/envoy-control-plane/control/internal"
)

var (
	l      internal.Logger
	port   uint
	nodeID string
)

func init() {
	l = internal.Logger{}

	flag.BoolVar(&l.Debug, "debug", false, "Enable xDS server debug logging")
	flag.UintVar(&port, "port", 18000, "xDS management server port")
	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID") // FIXME: nodeID is dynamically obtained within the program, not as a runtime argument.
}

func main() {
	flag.Parse()

	cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)

	snapshot := internal.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}
	l.Debugf("will serve snapshot %+v", snapshot)

	if err := cache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	ctx := context.Background()
	cb := &test.Callbacks{Debug: l.Debug}
	srv := server.NewServer(ctx, cache, cb)
	internal.RunServer(srv, port)
}
