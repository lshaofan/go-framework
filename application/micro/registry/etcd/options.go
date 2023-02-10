package etcd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

type Options struct {
	EtcdAddresses []string
	EtcdAddrUrl   string
	DialTimeout   int
	closeChan     chan struct{}
	leases        *clientv3.LeaseGrantResponse
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	srv           *RpcServer
	srvTTL        int64
	cli           *clientv3.Client
	logger        *logrus.Logger
}

func NewOptions() *Options {
	return &Options{}
}

type Option func(*Options)

// Cli is the etcd client
func (o *Options) Cli() (*clientv3.Client, error) {
	var err error
	if o.EtcdAddrUrl == "" {
		o.EtcdAddrUrl = "http://" + o.EtcdAddresses[0]
	}

	o.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   o.EtcdAddresses,
		DialTimeout: time.Duration(o.srvTTL) * time.Millisecond,
	})

	//o.cli, err = clientv3.NewFromURL(o.EtcdAddrUrl)

	if err != nil {
		fmt.Println("etcd", err)
		return nil, err
	}
	return o.cli, nil
}

// Leases is the etcd leases id
func (o *Options) Leases() (*clientv3.LeaseGrantResponse, error) {
	var err error
	// 判断cli是否存在
	if o.cli == nil {
		_, err = o.Cli()
		if err != nil {
			return nil, err
		}
	}
	// 判断leases是否存在
	if o.leases == nil {
		// 创建租约
		o.leases, err = o.cli.Grant(o.cli.Ctx(), o.srvTTL)
		if err != nil {
			return nil, err
		}
	}

	return o.leases, nil
}

// GetGrpcServer GetSrv is the etcd grpc server
func (o *Options) GetGrpcServer(i ...grpc.UnaryServerInterceptor) *grpc.Server {
	return o.srv.GetGRPCSvr(i...)
}

// SetSrv is the etcd grpc server
func (o *Options) SetSrv(srv *RpcServer) {
	o.srv = srv
}
func WithLogger(logger *logrus.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// RunGrpcServer RunSrv is the etcd grpc server
func (o *Options) RunGrpcServer() error {
	return o.srv.RunGrpcServer()
}

func WithEtcdAddresses(etcdAddresses []string) Option {
	return func(o *Options) {
		o.EtcdAddresses = etcdAddresses
	}
}

// WithEtcdAddrUrl is the etcd address url
func WithEtcdAddrUrl(etcdAddrUrl string) Option {
	return func(o *Options) {
		o.EtcdAddrUrl = etcdAddrUrl
	}
}

// WithDialTimeout is the etcd dial timeout
func WithDialTimeout(dialTimeout int) Option {
	return func(o *Options) {
		o.DialTimeout = dialTimeout
	}
}

// WithCloseChan is the etcd close channel
func WithCloseChan(closeChan chan struct{}) Option {
	return func(o *Options) {
		o.closeChan = closeChan
	}
}

// WithKeepAliveChan is the etcd keep alive channel
func WithKeepAliveChan(keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse) Option {
	return func(o *Options) {
		o.keepAliveChan = keepAliveChan
	}
}

// WithSrv is the etcd server
func WithSrv(srv *RpcServer) Option {
	return func(o *Options) {
		o.srv = srv
	}
}

// WithSrvTTL is the etcd server ttl
func WithSrvTTL(srvTTL int64) Option {
	if srvTTL == 0 {
		srvTTL = 10
	}
	return func(o *Options) {
		o.srvTTL = srvTTL
	}
}

// WithRpcServer is the etcd rpc server
func WithRpcServer(srv *RpcServer) Option {
	return func(o *Options) {
		o.srv = srv
	}
}
