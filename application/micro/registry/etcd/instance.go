package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lshaofan/go-framework/application/help/console"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
	"strings"
)

type RpcServer struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int    `json:"weight"`
	Listen  net.Listener
	srv     *grpc.Server
}

type ServerOption func(*RpcServer)

func NewRpcServer(name, addr string, opts ...ServerOption) *RpcServer {
	s := &RpcServer{
		Name: name,
		Addr: addr,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// GetRpcServerName 从服务地址中获取服务名
func (s *RpcServer) GetRpcServerName() string {
	return s.Name
}

// GetRpcServerVersion 从服务地址中获取服务版本
func (s *RpcServer) GetRpcServerVersion() string {
	return s.Version
}

// GetRpcServerAddr 从服务地址中获取服务地址
func (s *RpcServer) GetRpcServerAddr() string {
	return s.Addr
}

// GetGRPCSvr 获取Grpc服务
func (s *RpcServer) GetGRPCSvr(i ...grpc.UnaryServerInterceptor) *grpc.Server {
	if len(i) == 0 {
		s.srv = grpc.NewServer()
		return s.srv
	}
	if len(i) > 1 {
		panic(errors.New("too many interceptors"))
	}
	s.srv = grpc.NewServer(grpc.UnaryInterceptor(i[0]))
	return s.srv
}

// SetGrpcServer 设置Grpc服务
func (s *RpcServer) SetGrpcServer(srv *grpc.Server) {
	s.srv = srv
}

// InitGrpcResolver InitResolver 初始化服务发现
func (s *RpcServer) InitGrpcResolver() {
	if s.Name == "" {
		panic(errors.New("服务名不能为空"))
	}
	if s.Version == "" {
		panic(errors.New("服务版本不能为空"))
	}
	if s.Addr == "" {
		panic(errors.New("服务地址不能为空"))
	}

}

// RunGrpcServer 运行Grpc服务
func (s *RpcServer) RunGrpcServer() error {
	console.Success("grpc服务启动中...")
	var err error
	s.Listen, err = net.Listen("tcp", s.GetRpcServerAddr())
	if err != nil {
		return err
	}
	if err := s.srv.Serve(s.Listen); err != nil {
		console.Error(fmt.Sprintf("grpc 服务启动失败:%s", err))
		return err
	}
	return nil
}

// Close 关闭服务
func (s *RpcServer) Close() {
	err := s.Listen.Close()
	if err != nil {
		console.Error(fmt.Sprintf("grpc 服务关闭失败:%s", err))
		return
	}
	defer s.srv.GracefulStop()
}

// BuildKey 构建服务的key
func BuildKey(server RpcServer) string {
	if server.Version == "" {
		return fmt.Sprintf("/%s/", server.Name)
	}
	return fmt.Sprintf("/%s/%s/", server.Name, server.Version)
}

// BuildRegisterKey 构建注册服务的key
func BuildRegisterKey(server RpcServer) string {
	return fmt.Sprintf("%s%s", BuildKey(server), server.Addr)
}

// ParseValue 将服务信息序列化为json
func ParseValue(value []byte) (RpcServer, error) {
	var server RpcServer
	if err := json.Unmarshal(value, &server); err != nil {
		return server, err
	}
	return server, nil

}

// SplitKey 将服务的key分割为服务名和服务地址
func SplitKey(key string) (RpcServer, error) {
	var server RpcServer
	str := strings.Split(key, "/")
	if len(str) == 0 {
		return server, errors.New(fmt.Sprintf("key %s is invalid", key))
	}
	server.Addr = str[len(str)-1]
	return server, nil

}

// Exist 判断服务是否存在
func Exist(l []resolver.Address, addr resolver.Address) bool {
	for _, v := range l {
		if v.Addr == addr.Addr {
			return true
		}
	}
	return false

}
