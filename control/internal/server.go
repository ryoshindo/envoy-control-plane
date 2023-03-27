package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtime "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secret "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	grpcKeepAliveTime        = 30 * time.Second
	grpcKeepAliveTimeout     = 5 * time.Second
	grpcKeepAliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

type Server struct {
	xdsServer server.Server
}

func NewServer(ctx context.Context, cache cache.Cache, cb *test.Callbacks) *Server {
	srv := server.NewServer(ctx, cache, cb)
	return &Server{srv}
}

func (s *Server) registerServer(grpcServer *grpc.Server) {
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, s.xdsServer)
	endpoint.RegisterEndpointDiscoveryServiceServer(grpcServer, s.xdsServer)
	cluster.RegisterClusterDiscoveryServiceServer(grpcServer, s.xdsServer)
	route.RegisterRouteDiscoveryServiceServer(grpcServer, s.xdsServer)
	listener.RegisterListenerDiscoveryServiceServer(grpcServer, s.xdsServer)
	secret.RegisterSecretDiscoveryServiceServer(grpcServer, s.xdsServer)
	runtime.RegisterRuntimeDiscoveryServiceServer(grpcServer, s.xdsServer)
}

func (s *Server) Run(port uint) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepAliveTime,
			Timeout: grpcKeepAliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepAliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s.registerServer(grpcServer)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}

func registerServer(grpcServer *grpc.Server, server server.Server) {
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpoint.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	cluster.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	route.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listener.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secret.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtime.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
}

// RunServer starts an xDS server at the given port
func RunServer(srv server.Server, port uint) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepAliveTime,
			Timeout: grpcKeepAliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepAliveMinTime,
			PermitWithoutStream: true,
		}),
	)

	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	registerServer(grpcServer, srv)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}
