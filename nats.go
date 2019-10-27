package nats

import (
	"github.com/im-kulikov/helium/module"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type (
	// Config alias
	Config = nats.Options

	// StreamerConfig for NSS client
	StreamerConfig struct {
		ClientID  string
		ClusterID string
		Options   []stan.Option
	}

	StreamerParams struct {
		dig.In

		Bus              *Client
		Viper            *viper.Viper
		OnConnectionLost stan.ConnectionLostHandler `optional:"true"`
	}

	// Client alias
	Client = nats.Conn

	// Error is constant error
	Error string
)

const (
	// ErrEmptyConfig when given empty options
	ErrEmptyConfig = Error("nats empty config")
	// ErrEmptyStreamerConfig when given empty options
	ErrEmptyStreamerConfig = Error("nats-streamer empty config")
	// ErrEmptyConnection when empty nats.Conn
	ErrEmptyConnection = Error("nats connection empty")
	// ErrClusterIDEmpty when empty clusterID
	ErrClusterIDEmpty = Error("nats.cluster_id cannot be empty")
	// ErrClientIDEmpty when empty clientID
	ErrClientIDEmpty = Error("nats.client_id cannot be empty")
)

// Error returns error message string
func (e Error) Error() string { return string(e) }

var (
	// Module is default Nats client
	Module = module.Module{
		{Constructor: NewDefaultConfig},
		{Constructor: NewConnection},
		{Constructor: NewDefaultStreamerConfig},
		{Constructor: NewStreamer},
	}
)

// NewDefaultConfig default settings for connection
func NewDefaultConfig(v *viper.Viper) (*Config, error) {
	if !v.IsSet("nats") {
		return nil, ErrEmptyConfig
	}

	var servers []string
	if v.IsSet("nats.servers") {
		servers = v.GetStringSlice("nats.servers")
	}

	return &Config{
		Url:              v.GetString("nats.url"),
		Servers:          servers,
		NoRandomize:      v.GetBool("nats.no_randomize"),
		Name:             v.GetString("nats.name"),
		Verbose:          v.GetBool("nats.verbose"),
		Pedantic:         v.GetBool("nats.pedantic"),
		Secure:           v.GetBool("nats.secure"),
		AllowReconnect:   v.GetBool("nats.allow_reconnect"),
		MaxReconnect:     v.GetInt("nats.max_reconnect"),
		ReconnectWait:    v.GetDuration("nats.reconnect_wait"),
		Timeout:          v.GetDuration("nats.timeout"),
		FlusherTimeout:   v.GetDuration("nats.flusher_timeout"),
		PingInterval:     v.GetDuration("nats.ping_interval"),
		MaxPingsOut:      v.GetInt("nats.max_pings_out"),
		ReconnectBufSize: v.GetInt("nats.reconnect_buf_size"),
		SubChanLen:       v.GetInt("nats.sub_chan_len"),
		User:             v.GetString("nats.user"),
		Password:         v.GetString("nats.password"),
		Token:            v.GetString("nats.token"),
	}, nil
}

// NewDefaultStreamerConfig default settings for streaming connection
func NewDefaultStreamerConfig(p StreamerParams) (*StreamerConfig, error) {
	if !p.Viper.IsSet("nats") {
		return nil, ErrEmptyConfig
	}

	var clusterID, clientID string
	if clusterID = p.Viper.GetString("nats.cluster_id"); clusterID == "" {
		return nil, ErrClusterIDEmpty
	}

	if clientID = p.Viper.GetString("nats.client_id"); clientID == "" {
		return nil, ErrClientIDEmpty
	}

	// set options:
	options := []stan.Option{stan.NatsConn(p.Bus)}

	// ConnectWait(t time.Duration)
	if v := p.Viper.GetDuration("nats.stan.connect_wait"); v > 0 {
		options = append(options, stan.ConnectWait(v))
	}

	// PubAckWait(t time.Duration)
	if v := p.Viper.GetDuration("nats.stan.pub_ack_wait"); v > 0 {
		options = append(options, stan.PubAckWait(v))
	}

	// MaxPubAcksInflight(max int)
	if v := p.Viper.GetInt("nats.stan.max_pub_acks_inflight"); v > 0 {
		options = append(options, stan.MaxPubAcksInflight(v))
	}

	// Pings(interval, maxOut int)
	pingMaxOut := p.Viper.GetInt("nats.stan.ping_max_out")
	pingInterval := p.Viper.GetInt("nats.stan.ping_interval")
	if pingMaxOut > 0 && pingInterval > 0 {
		options = append(options, stan.Pings(pingInterval, pingMaxOut))
	}

	// SetConnectionLostHandler(handler ConnectionLostHandler)
	if p.OnConnectionLost != nil {
		options = append(options, stan.SetConnectionLostHandler(p.OnConnectionLost))
	}

	return &StreamerConfig{
		ClientID:  clientID,
		ClusterID: clusterID,
		Options:   options,
	}, nil
}

// NewConnection of nats client
func NewConnection(opts *Config) (bus *Client, err error) {
	if opts == nil {
		return nil, ErrEmptyConfig
	}

	if bus, err = opts.Connect(); err != nil {
		return nil, err
	}

	return bus, nil
}

// NewStreamer is nats-streamer client
func NewStreamer(opts *StreamerConfig) (stan.Conn, error) {
	if opts == nil {
		return nil, ErrEmptyStreamerConfig
	}

	if opts.Options == nil || len(opts.Options) == 0 {
		return nil, ErrEmptyConnection
	}

	return stan.Connect(opts.ClusterID, opts.ClientID, opts.Options...)
}
