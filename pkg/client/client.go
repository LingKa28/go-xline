package client

import (
	"context"
	"fmt"

	xlinepb "github.com/xline-kv/go-xline/api/gen/xline"
	"github.com/xline-kv/go-xline/pkg/curp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Xline client
type client struct {
	// Kv client
	KV
	// Lease client
	// svc.Lease
	// Watch client
	// svc.Watch
	// Lock client
	// svc.Lock
	// Auth client
	// svc.Auth
	// Maintenance client
	// svc.Maintenance
	// Cluster client
	// svc.Cluster
}

func Connect(allMembers []string, option *ClientOptions) (*client, error) {
	// conn, err := grpc.Dial(allMembers[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	panic(err)
	// }

	cfg := curp.NewDefaultClientConfig()
	if option != nil && option.CurpConfig != nil {
		cfg = option.CurpConfig
	}
	curpClient, e := curp.NewClientBuilder(cfg).DiscoveryFrom(allMembers).Build()
	if e != nil {
		panic(e)
	}

	// idGen := svc.NewLeaseId()

	var token *string = nil
	if option != nil && option.User != nil {
		conn, err := grpc.Dial(allMembers[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("request token fail. %v", err)
		}
		authClient := xlinepb.NewAuthClient(conn)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		res, err := authClient.Authenticate(ctx, &xlinepb.AuthenticateRequest{Name: option.User.Name, Password: option.User.Password})
		if err != nil {
			panic(err)
		}
		token = &res.Token
	}

	kv := newKV(curpClient, token)
	// lease := svc.NewLease(curpClient, conn, token, idGen)
	// watch := svc.NewWatch(conn)
	// lock := svc.NewLock(curpClient, lease, watch, token)
	// auth := svc.NewAuth(curpClient, token)
	// maintenance := NewMaintenance(conn)
	// cluster := NewCluster(conn)

	return &client{
		KV: kv,
		// Lease:       lease,
		// Watch:       watch,
		// Lock:        lock,
		// Auth:        auth,
		// Maintenance: maintenance,
		// Cluster:     cluster,
	}, nil
}

// Options for a client connection
type ClientOptions struct {
	// User is a pair values of name and password
	User *UserCredentials
	// Timeout settings for the curp client
	CurpConfig *curp.ClientConfig
}

// Options for a user
type UserCredentials struct {
	// Username
	Name string
	// Password
	Password string
}
