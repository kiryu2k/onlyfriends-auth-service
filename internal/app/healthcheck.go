package app

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthcheck(srv *grpc.Server) {
	healthcheck := health.NewServer()
	healthcheck.SetServingStatus("onlyfriends-auth-service", healthpb.HealthCheckResponse_SERVING)
	healthgrpc.RegisterHealthServer(srv, healthcheck)
}
