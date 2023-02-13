package etcd

import (
	"context"
	"fmt"
	"github.com/lshaofan/go-framework/application/help/console"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	namResolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	gresolver "google.golang.org/grpc/resolver"
)

type Registry struct {
	options Options
	ctx     context.Context
}

func NewRegistry(opts ...Option) *Registry {
	e := &Registry{
		options: Options{},
	}
	configure(e, opts...)
	return e
}
func configure(r *Registry, opts ...Option) error {
	for _, o := range opts {
		o(&r.options)
	}
	// 判断是否有传入grpc server
	if r.options.srv.Name == "" {
		panic("grpc server is nil")
	}
	return nil
}

// GetClient GetCli is the etcd client
func (e *Registry) GetClient() (*clientv3.Client, error) {
	return e.options.Cli()
}

// GetGrpcServer GetSrv is the etcd grpc server
func (e *Registry) GetGrpcServer(i ...grpc.UnaryServerInterceptor) *grpc.Server {
	return e.options.GetGrpcServer(i...)
}

// GetEtcdManager GetManager is the etcd manager
func (e *Registry) GetEtcdManager(name string) (em endpoints.Manager, err error) {

	client, err := e.GetClient()
	if err != nil {
		return
	}
	em, err = endpoints.NewManager(client, name)
	if err != nil {
		return
	}

	return
}

// GetLeasGrant  is the etcd lease
func (e *Registry) GetLeasGrant() (lease *clientv3.LeaseGrantResponse, err error) {
	return e.options.Leases()
}

// AddEndpoint add endpoint
func (e *Registry) AddEndpoint(servers ...*RpcServer) error {
	for _, srv := range servers {
		em, err := e.GetEtcdManager(srv.GetRpcServerName())
		if err != nil {
			return err
		}
		lease, err := e.GetLeasGrant()
		if err != nil {
			return err
		}
		err = em.AddEndpoint(
			context.TODO(),
			fmt.Sprintf("%s/%s", srv.GetRpcServerName(),
				srv.GetRpcServerAddr()),
			endpoints.Endpoint{Addr: srv.GetRpcServerAddr()},
			clientv3.WithLease(lease.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

// KeepAlive keep alive
func (e *Registry) KeepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	cli, err := e.GetClient()
	if err != nil {
		return nil, err
	}

	les, err := e.GetLeasGrant()
	if err != nil {
		return nil, err
	}
	return cli.KeepAlive(context.Background(), les.ID)

}

func (e *Registry) Init() error {
	c, err := e.KeepAlive()
	if err != nil {
		return err
	}
	err = e.AddEndpoint(e.options.srv)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			// 保持心跳 todo：先取走channel中的数据，再进行下一次心跳避免出现channel阻塞警告
			case <-c:

			}
		}
	}()
	return nil
}

// GetResolverBuilder get resolver builder
func (e *Registry) GetResolverBuilder(client *clientv3.Client) (gresolver.Builder, error) {
	return namResolver.NewBuilder(client)
}

// GetClientConn get client conn
func (e *Registry) GetClientConn(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	cli, err := e.GetClient()
	if err != nil {
		return nil, err
	}
	builder, err := e.GetResolverBuilder(cli)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("etcd:///%s", target),
		append(opts,
			grpc.WithResolvers(builder),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
				roundrobin.Name,
			)),
		)...,
	)
	return conn, err
}

// Run run
func (e *Registry) Run() error {
	err := e.options.RunGrpcServer()
	if err != nil {

		return err
	}
	return nil
}

func (e *Registry) Close() {
	err := e.etcdUnRegister()
	if err != nil {
		fmt.Println(err)
	}
	e.options.srv.Close()
}

func (e *Registry) etcdUnRegister() error {
	cli, err := e.GetClient()
	if err != nil {
		console.Error(fmt.Sprintf("etcdUnRegister err:%s", err))
		return err
	}
	if cli != nil {
		em, err := e.GetEtcdManager(e.options.srv.GetRpcServerName())
		if err != nil {
			console.Error(fmt.Sprintf("etcdUnRegister err:%s", err))
			return err
		}
		err = em.DeleteEndpoint(
			context.TODO(),
			fmt.Sprintf("%s/%s", e.options.srv.GetRpcServerName(),
				e.options.srv.GetRpcServerAddr(),
			))

		if err != nil {
			console.Error(fmt.Sprintf("etcdUnRegister err:%s", err))
			return err
		}
	}

	return nil
}
