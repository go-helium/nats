package nats

import (
	"testing"

	"github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func RunServer(ID string) *server.StanServer {
	s, err := server.RunServer(ID)
	if err != nil {
		panic(err)
	}
	return s
}

func TestNewDefaultConfig(t *testing.T) {
	t.Run("must fail on empty", func(t *testing.T) {
		v := viper.New()
		c, err := NewDefaultConfig(v)
		require.Nil(t, c)
		require.Error(t, err)
	})

	t.Run("servers should be nil", func(t *testing.T) {
		v := viper.New()
		v.SetDefault("nats.url", "something")

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Nil(t, c.Servers)
	})

	t.Run("servers should be slice of string", func(t *testing.T) {
		v := viper.New()
		v.SetDefault("nats.servers", "something")

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Len(t, c.Servers, 1)
		require.Equal(t, c.Servers[0], "something")
	})

	t.Run("should be ok", func(t *testing.T) {
		v := viper.New()
		url := "something"
		v.SetDefault("nats.url", url)

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Equal(t, c.Url, url)
	})

	t.Run("should fail for empty config", func(t *testing.T) {
		c, err := NewConnection(nil)
		require.Nil(t, c)
		require.EqualError(t, err, ErrEmptyConfig.Error())
	})

	t.Run("should fail for empty config on nats-stremer", func(t *testing.T) {
		c, err := NewStreamer(nil)
		require.Nil(t, c)
		require.EqualError(t, err, ErrEmptyStreamerConfig.Error())
	})

	t.Run("should fail client", func(t *testing.T) {
		v := viper.New()

		v.SetDefault("nats.url", nats.DefaultURL)

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Equal(t, c.Url, nats.DefaultURL)

		cli, err := NewConnection(c)
		require.Nil(t, cli)
		require.Error(t, err)
	})

	t.Run("should not fail with test server", func(t *testing.T) {
		v := viper.New()
		serve := RunServer(nats.DefaultURL)
		defer serve.Shutdown()

		v.SetDefault("nats.url", nats.DefaultURL)

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)

		cli, err := NewConnection(c)
		require.NoError(t, err)
		require.NotNil(t, cli)
	})

	t.Run("should fail with empty config", func(t *testing.T) {
		v := viper.New()
		cfg, err := NewDefaultStreamerConfig(StreamerParams{Viper: v})
		require.Nil(t, cfg)
		require.EqualError(t, err, ErrEmptyConfig.Error())
	})

	t.Run("should fail with empty clusterID", func(t *testing.T) {
		v := viper.New()
		v.SetDefault("nats.cluster_id", "")
		cfg, err := NewDefaultStreamerConfig(StreamerParams{Viper: v})
		require.Nil(t, cfg)
		require.EqualError(t, err, ErrClusterIDEmpty.Error())
	})

	t.Run("should fail with empty clientID", func(t *testing.T) {
		v := viper.New()
		v.SetDefault("nats.cluster_id", "myCluster")
		cfg, err := NewDefaultStreamerConfig(StreamerParams{Viper: v})
		require.Nil(t, cfg)
		require.EqualError(t, err, ErrClientIDEmpty.Error())
	})

	t.Run("should fail on connection empty", func(t *testing.T) {
		v := viper.New()
		v.SetDefault("nats.url", nats.DefaultURL)
		v.SetDefault("nats.client_id", "myClient")
		v.SetDefault("nats.cluster_id", "myCluster")

		cfg, err := NewDefaultStreamerConfig(StreamerParams{Viper: v})
		require.NoError(t, err)

		cfg.Options = nil

		serve, err := NewStreamer(cfg)
		require.Nil(t, serve)
		require.EqualError(t, err, ErrEmptyConnection.Error())
	})

	t.Run("should run streamer client", func(t *testing.T) {
		v := viper.New()
		v.SetDefault("nats.client_id", "myClient")
		v.SetDefault("nats.cluster_id", "myCluster")

		// stan options:
		v.SetDefault("nats.stan.connect_wait", stan.DefaultConnectWait)
		v.SetDefault("nats.stan.pub_ack_wait", stan.DefaultAckWait)
		v.SetDefault("nats.stan.ping_max_out", stan.DefaultPingMaxOut)
		v.SetDefault("nats.stan.ping_interval", stan.DefaultPingInterval)
		v.SetDefault("nats.stan.max_pub_acks_inflight", stan.DefaultMaxPubAcksInflight)

		// Run a NATS Streaming server
		s := RunServer("myCluster")
		defer s.Shutdown()

		con, err := nats.Connect(nats.DefaultURL)
		require.NoError(t, err)

		defer con.Close()

		cfg, err := NewDefaultStreamerConfig(StreamerParams{
			Viper:            v,
			Bus:              con,
			OnConnectionLost: func(stan.Conn, error) {},
		})
		require.NoError(t, err)

		st, err := NewStreamer(cfg)
		require.NoError(t, err)
		require.NoError(t, st.Close())
	})
}
