package micro

import (
	"context"
	"errors"
	"fmt"
	"github.com/lshaofan/go-framework/application/help/console"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	gresolver "google.golang.org/grpc/resolver"
	"net"
	"time"
)

type Option func(s *Server)

func WithEndpoints(endpoints []string) Option {
	return func(s *Server) {
		s.etcdEndPoints = endpoints
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.etcdDialTimeout = timeout
	}
}

func WithEtcdConfig(config *clientv3.Config) Option {
	return func(s *Server) {
		s.etcdConfig = config
	}
}

func WithGrpcListener(listener net.Listener) Option {
	return func(s *Server) {
		s.grpcListener = listener
	}
}

func WithGrpcServer(grpcSvr *grpc.Server) Option {
	return func(s *Server) {
		s.grpcSvr = grpcSvr
	}
}

func WithService(service *Service) Option {
	return func(s *Server) {
		s.service = service
	}
}

type ServiceOption func(s *Service)

// Service  服务
type Service struct {
	// 服务名称
	ServiceName string
	// 服务地址
	ServiceAddr string
}

func WithServiceName(serviceName string) ServiceOption {
	return func(s *Service) {
		s.ServiceName = serviceName
	}
}

func WithServiceAddr(serviceAddr string) ServiceOption {
	return func(s *Service) {
		s.ServiceAddr = serviceAddr
	}
}

func NewService(opts ...ServiceOption) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Server struct {
	grpcListener    net.Listener
	grpcSvr         *grpc.Server
	etcdEndPoints   []string
	etcdDialTimeout time.Duration
	etcdConfig      *clientv3.Config
	service         *Service
	etcdClient      *clientv3.Client
	etcdMonitor     endpoints.Manager
	etcdLease       *clientv3.LeaseGrantResponse
	etcdTtl         int64
	ctx             context.Context
	leaseID         clientv3.LeaseID
	etcdResolver    gresolver.Builder
}

func NewServer(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Init() error {
	var err error
	// 没有传入service
	if s.service == nil {
		return errors.New("service is nil")
	}
	// 没有传入grpc server listener抛出错误
	if s.grpcListener == nil {
		s.grpcListener, err = net.Listen("tcp", s.service.ServiceAddr)
		if err != nil {
			return err
		}
	}
	if s.etcdTtl == 0 {
		s.etcdTtl = 5
	}
	// 没有传入grpc server 初始化一个
	if s.grpcSvr == nil {
		s.grpcSvr = grpc.NewServer()
	}
	// 没有传入etcd config
	if s.etcdConfig == nil {
		s.etcdConfig = &clientv3.Config{
			Endpoints:   s.etcdEndPoints,
			DialTimeout: s.etcdDialTimeout,
		}
	} else {
		// 传入了etcd config 但是没有传入etcd endpoints
		if len(s.etcdConfig.Endpoints) == 0 {
			s.etcdConfig.Endpoints = s.etcdEndPoints
		}
		// 传入了etcd config 但是没有传入etcd dial timeout
		if s.etcdConfig.DialTimeout == 0 {
			s.etcdConfig.DialTimeout = s.etcdDialTimeout
		}
	}
	s.ctx = context.Background()
	s.etcdClient, err = clientv3.New(*s.etcdConfig)
	if err != nil {
		return err
	}
	s.etcdResolver, err = resolver.NewBuilder(s.etcdClient)
	if err != nil {
		return err
	}
	err = s.RegisterService()
	if err != nil {
		return err
	}

	return nil
}

// Start 启动grpc server
func (s *Server) Start() error {
	console.Success(fmt.Sprintf("grpc 服务启动 at %s", s.service.ServiceAddr))
	if err := s.grpcSvr.Serve(s.grpcListener); err != nil {
		return err
	}
	console.Error("grpc server start error")

	return nil
}

// Stop 停止grpc server
func (s *Server) Stop() {
	s.etcdClient.Revoke(s.ctx, s.leaseID)
	s.grpcSvr.Stop()
	s.etcdClient.Close()
	s.UnRegisterService()
	s.grpcListener.Close()
}

// GetGrpcServer 获取grpc server
func (s *Server) GetGrpcServer() *grpc.Server {
	return s.grpcSvr
}

// GetEtcdClient 获取etcd client
func (s *Server) GetEtcdClient() *clientv3.Client {
	return s.etcdClient
}

// GetService 获取service
func (s *Server) GetService() *Service {
	return s.service
}

// UnRegisterService 注销服务
func (s *Server) UnRegisterService() error {
	if s.etcdClient == nil || s.etcdMonitor == nil || s.etcdLease == nil {
		return nil
	}
	err := s.etcdMonitor.DeleteEndpoint(
		s.ctx,
		fmt.Sprintf("%s/%s", s.service.ServiceName, s.service.ServiceAddr),
	)
	if err != nil {
		return err
	}
	return nil
}

// ServiceDiscovery 服务发现
func (s *Server) ServiceDiscovery(name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("etcd:///%s", name),
		append(opts,

			grpc.WithResolvers(s.etcdResolver),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
				roundrobin.Name,
			)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return conn, err
}

// RegisterService 注册服务
func (s *Server) RegisterService() error {
	var err error

	s.etcdLease, err = s.etcdClient.Grant(s.ctx, s.etcdTtl)
	if err != nil {
		console.Error(fmt.Sprintf("etcd grant 失败:%s", err))
		return err
	}

	s.etcdMonitor, err = endpoints.NewManager(s.etcdClient, s.service.ServiceName)
	if err != nil {
		return err
	}
	err = s.etcdMonitor.AddEndpoint(
		s.ctx,
		fmt.Sprintf("%s/%s", s.service.ServiceName, s.service.ServiceAddr),
		endpoints.Endpoint{
			Addr: s.service.ServiceAddr,
		},
		clientv3.WithLease(s.etcdLease.ID),
	)

	if err != nil {
		return err
	}
	alive, err := s.etcdClient.KeepAlive(s.ctx, s.etcdLease.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case _, ok := <-alive:
				if !ok {
					return
				}
			}
		}

	}()
	return nil
}
