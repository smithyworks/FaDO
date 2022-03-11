package database

type CaddyConfig struct {
	Apps struct {
		HTTP struct {
			Servers map[string]LoadBalancerServerConfig `json:"servers,omitempty"`
		} `json:"http,omitempty"`
	} `json:"apps,omitempty"`
}

type LoadBalancerServerConfig struct {
	Listen []string `json:"listen,omitempty"`
	Routes []LoadBalancerRouteConfig `json:"routes,omitempty"`
	Terminal bool `json:"terminal,omitempty"`
}

type SelectionPolicyConfig struct {
	Policy string `json:"policy,omitempty"`
}

type LoadBalancingConfig struct {
	SelectionPolicy SelectionPolicyConfig `json:"selection_policy,omitempty"`
}

type UpstreamConfig struct {
	Dial string `json:"dial,omitempty"`
}

type HandleConfig struct {
	Handler string `json:"handler,omitempty"`
	LoadBalancing *LoadBalancingConfig `json:"load_balancing,omitempty"`
	Upstreams []UpstreamConfig `json:"upstreams,omitempty"`
	Routes []LoadBalancerRouteConfig `json:"routes,omitempty"`
}

type MatchConfig struct {
	Header map[string][]string `json:"header,omitempty"`
	Host []string `json:"host,omitempty"`
}

type LoadBalancerRouteConfig struct {
	Handle []HandleConfig `json:"handle,omitempty"`
	Match []MatchConfig `json:"match,omitempty"`
	Terminal *bool `json:"terminal,omitempty"`
}

// Override types

type LoadBalancerRouteSettings struct {
	BucketName string `json:"bucket_name"`
	Policy string `json:"policy,omitempty"`
	Upstreams []string `json:"upstreams,omitempty"`
}
