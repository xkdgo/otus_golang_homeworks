package internalgrpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

func InterceptorWithLogger(logger Logger) grpc.UnaryServerInterceptor {
	fn := func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("middlware started")
		start := time.Now()

		clientRemoteAddr := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			clientRemoteAddr = p.Addr.String()
		}

		useragent := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ua, ok := md["user-agent"]; ok {
				useragent = strings.Join(ua, ",")
			}
		}
		method := info.FullMethod
		h, err := handler(ctx, req)
		// after executing rpc
		s, _ := status.FromError(err)
		logger.Infof("%s %s %s %s %s %s %s %s",
			clientRemoteAddr,
			start.Format("[02/Jan/2006:15:04:05 -0700]"),
			method,
			"",
			"prtoto3",
			s.Code().String(),
			time.Since(start),
			useragent,
		)
		return h, err
	}
	return fn
}
