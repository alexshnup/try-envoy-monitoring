// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package example

import (
	"fmt"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"google.golang.org/protobuf/types/known/durationpb"

	tcp_proxy "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
)

const (
	ClusterName = "redis_proxy_cluster"
	// RouteName    = "local_route"
	ListenerName = "listener_0"
	ListenerPort = 6666
	UpstreamPort = 6379
)

// func makeCluster(clusterName string) *cluster.Cluster {
// 	return &cluster.Cluster{
// 		Name:                 clusterName,
// 		ConnectTimeout:       durationpb.New(5 * time.Second),
// 		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
// 		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
// 		LoadAssignment:       makeEndpoint(clusterName),
// 		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
// 	}
// }

// func makeEndpoint(clusterName string) *endpoint.ClusterLoadAssignment {
// 	return &endpoint.ClusterLoadAssignment{
// 		ClusterName: clusterName,
// 		Endpoints: []*endpoint.LocalityLbEndpoints{{
// 			LbEndpoints: []*endpoint.LbEndpoint{{
// 				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
// 					Endpoint: &endpoint.Endpoint{
// 						Address: &core.Address{
// 							Address: &core.Address_SocketAddress{
// 								SocketAddress: &core.SocketAddress{
// 									Protocol: core.SocketAddress_TCP,
// 									// Address:  UpstreamHost,
// 									PortSpecifier: &core.SocketAddress_PortValue{
// 										PortValue: UpstreamPort,
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			}},
// 		}},
// 	}
// }

// func makeRoute(routeName string, clusterName string) *route.RouteConfiguration {
// 	return &route.RouteConfiguration{
// 		Name: routeName,
// 		VirtualHosts: []*route.VirtualHost{{
// 			Name:    "local_service",
// 			Domains: []string{"*"},
// 			Routes: []*route.Route{{
// 				Match: &route.RouteMatch{
// 					PathSpecifier: &route.RouteMatch_Prefix{
// 						Prefix: "/",
// 					},
// 				},
// 				Action: &route.Route_Route{
// 					Route: &route.RouteAction{
// 						ClusterSpecifier: &route.RouteAction_Cluster{
// 							Cluster: clusterName,
// 						},
// 						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
// 							HostRewriteLiteral: UpstreamHost,
// 						},
// 					},
// 				},
// 			}},
// 		}},
// 	}
// }

func MakeTCPListener(listenerName string) *listener.Listener {
	tcpProxyConfig := &tcp_proxy.TcpProxy{
		StatPrefix: "tcp_proxy",
		ClusterSpecifier: &tcp_proxy.TcpProxy_Cluster{
			Cluster: ClusterName,
		},
	}

	tcpProxyConfigPbst, err := ptypes.MarshalAny(tcpProxyConfig)
	if err != nil {
		// handle error
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.TCPProxy,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: tcpProxyConfigPbst,
				},
			}},
		}},
	}
}

// func makeHTTPListener(listenerName string, route string) *listener.Listener {
// 	routerConfig, _ := anypb.New(&router.Router{})
// 	// HTTP filter configuration
// 	manager := &hcm.HttpConnectionManager{
// 		CodecType:  hcm.HttpConnectionManager_AUTO,
// 		StatPrefix: "http",
// 		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
// 			Rds: &hcm.Rds{
// 				ConfigSource:    makeConfigSource(),
// 				RouteConfigName: route,
// 			},
// 		},
// 		HttpFilters: []*hcm.HttpFilter{{
// 			Name:       "http-router",
// 			ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
// 		}},
// 	}
// 	pbst, err := anypb.New(manager)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &listener.Listener{
// 		Name: listenerName,
// 		Address: &core.Address{
// 			Address: &core.Address_SocketAddress{
// 				SocketAddress: &core.SocketAddress{
// 					Protocol: core.SocketAddress_TCP,
// 					Address:  "0.0.0.0",
// 					PortSpecifier: &core.SocketAddress_PortValue{
// 						PortValue: ListenerPort,
// 					},
// 				},
// 			},
// 		},
// 		FilterChains: []*listener.FilterChain{{
// 			Filters: []*listener.Filter{{
// 				Name: "http-connection-manager",
// 				ConfigType: &listener.Filter_TypedConfig{
// 					TypedConfig: pbst,
// 				},
// 			}},
// 		}},
// 	}
// }

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

func GenerateSnapshot() *cache.Snapshot {
	snap, _ := cache.NewSnapshot("version-0",
		map[resource.Type][]types.Resource{
			// resource.ClusterType: {makeCluster(ClusterName)},
			resource.ClusterType: {MakeRedisCluster(ClusterName, "redis-master")},
			// resource.RouteType:   {makeRoute(RouteName, ClusterName)},
			// resource.ListenerType: {makeHTTPListener(ListenerName, RouteName)},
			resource.ListenerType: {MakeTCPListener(ListenerName)},
		},
	)
	return snap
}

// func MakeRedisCluster(clusterName, masterNode string) *cluster.Cluster {
// 	return &cluster.Cluster{
// 		Name:                 clusterName,
// 		ConnectTimeout:       ptypes.DurationProto(5 * time.Second),
// 		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STRICT_DNS},
// 		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
// 		LoadAssignment:       MakeRedisEndpoint(clusterName, masterNode),
// 		HealthChecks: []*core.HealthCheck{
// 			{
// 				Timeout:            ptypes.DurationProto(5 * time.Second),
// 				Interval:           ptypes.DurationProto(10 * time.Second),
// 				UnhealthyThreshold: &wrappers.UInt32Value{Value: 1},
// 				HealthyThreshold:   &wrappers.UInt32Value{Value: 2},
// 				HealthChecker: &core.HealthCheck_TcpHealthCheck_{
// 					TcpHealthCheck: &core.HealthCheck_TcpHealthCheck{
// 						Send:    &core.HealthCheck_Payload{Payload: &core.HealthCheck_Payload_Text{Text: "PING\r\n"}},
// 						Receive: []*core.HealthCheck_Payload{{Payload: &core.HealthCheck_Payload_Text{Text: "+PONG"}}},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

func MakeRedisCluster(clusterName, masterNode string) *cluster.Cluster {
	// Hexadecimal representation of "PING\r\n"
	pingHex := fmt.Sprintf("%x", "PING\r\n") // convert to hex
	// Hexadecimal representation of "+PONG\r\n"
	pongHex := fmt.Sprintf("%x", "+PONG\r\n") // convert to hex

	return &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STRICT_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       MakeRedisEndpoint(clusterName, masterNode),
		HealthChecks: []*core.HealthCheck{
			{
				Timeout:            durationpb.New(5 * time.Second),
				Interval:           durationpb.New(10 * time.Second),
				UnhealthyThreshold: &wrappers.UInt32Value{Value: 1},
				HealthyThreshold:   &wrappers.UInt32Value{Value: 2},
				HealthChecker: &core.HealthCheck_TcpHealthCheck_{
					TcpHealthCheck: &core.HealthCheck_TcpHealthCheck{
						Send: &core.HealthCheck_Payload{
							Payload: &core.HealthCheck_Payload_Text{
								Text: pingHex,
							},
						},
						Receive: []*core.HealthCheck_Payload{
							{
								Payload: &core.HealthCheck_Payload_Text{
									Text: pongHex,
								},
							},
						},
					},
				},
			},
		},
	}
}

func MakeRedisEndpoint(clusterName, masterNode string) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{
			{
				LbEndpoints: []*endpoint.LbEndpoint{
					{
						HostIdentifier: &endpoint.LbEndpoint_Endpoint{
							Endpoint: &endpoint.Endpoint{
								Address: &core.Address{
									Address: &core.Address_SocketAddress{
										SocketAddress: &core.SocketAddress{
											Protocol: core.SocketAddress_TCP,
											Address:  masterNode,
											PortSpecifier: &core.SocketAddress_PortValue{
												PortValue: UpstreamPort, // 6379 Default Redis port
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
