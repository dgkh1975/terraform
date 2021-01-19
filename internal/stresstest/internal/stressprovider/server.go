package stressprovider

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform/internal/tfplugin5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// While we're doing normal stress-testing we just run the stressprovider
// in-process and access its API via normal function calls, but we also have
// an RPC implementation of it which is used by the "stresstest terraform"
// command so that the provider can be made available to a normal Terraform CLI
// process while debugging a test failure.

type Plugin struct {
	base *Provider
}

var _ plugin.GRPCPlugin = (*Plugin)(nil)
var _ plugin.Plugin = (*Plugin)(nil)

func (p *Provider) Plugin() *Plugin {
	return &Plugin{p}
}

func (p *Plugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	inst := p.base.NewInstance()
	tfplugin5.RegisterProviderServer(server, inst.Server())
	return nil
}

func (p *Plugin) GRPCClient(context.Context, *plugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	return nil, fmt.Errorf("this is only a server")
}

func (p *Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return nil, fmt.Errorf("go-plugin net/rpc is obsolete")
}

func (p *Plugin) Client(*plugin.MuxBroker, *rpc.Client) (interface{}, error) {
	return nil, fmt.Errorf("go-plugin net/rpc is obsolete")
}

type Server struct {
	provider *Provider
}

func (p *Provider) Server() *Server {
	return &Server{p}
}

var _ tfplugin5.ProviderServer = (*Server)(nil)

func (s *Server) GetSchema(ctx context.Context, req *tfplugin5.GetProviderSchema_Request) (*tfplugin5.GetProviderSchema_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) PrepareProviderConfig(ctx context.Context, req *tfplugin5.PrepareProviderConfig_Request) (*tfplugin5.PrepareProviderConfig_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) ValidateResourceTypeConfig(ctx context.Context, req *tfplugin5.ValidateResourceTypeConfig_Request) (*tfplugin5.ValidateResourceTypeConfig_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) ValidateDataSourceConfig(ctx context.Context, req *tfplugin5.ValidateDataSourceConfig_Request) (*tfplugin5.ValidateDataSourceConfig_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) UpgradeResourceState(ctx context.Context, req *tfplugin5.UpgradeResourceState_Request) (*tfplugin5.UpgradeResourceState_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) Configure(ctx context.Context, req *tfplugin5.Configure_Request) (*tfplugin5.Configure_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) ReadResource(ctx context.Context, req *tfplugin5.ReadResource_Request) (*tfplugin5.ReadResource_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) PlanResourceChange(ctx context.Context, req *tfplugin5.PlanResourceChange_Request) (*tfplugin5.PlanResourceChange_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) ApplyResourceChange(ctx context.Context, req *tfplugin5.ApplyResourceChange_Request) (*tfplugin5.ApplyResourceChange_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) ImportResourceState(ctx context.Context, req *tfplugin5.ImportResourceState_Request) (*tfplugin5.ImportResourceState_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) ReadDataSource(ctx context.Context, req *tfplugin5.ReadDataSource_Request) (*tfplugin5.ReadDataSource_Response, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) Stop(ctx context.Context, req *tfplugin5.Stop_Request) (*tfplugin5.Stop_Response, error) {
	// This provider isn't stoppable, because it's doing all of its work
	// locally in memory anyway.
	return nil, nil
}
