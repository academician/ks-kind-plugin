package main

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/danielfoehrkn/kubeswitch/pkg/store/plugins"
	storetypes "github.com/danielfoehrkn/kubeswitch/pkg/store/types"

	kindcluster "sigs.k8s.io/kind/pkg/cluster"
)

// Store is the implementation of the store plugin
//
// Implement [kindlog.Logger]
type Store struct {
	Logger   hclog.Logger
	internal bool
	provider *kindcluster.Provider
}

func NewStore(logger hclog.Logger) *Store {
	s := &Store{
		Logger:   logger,
		internal: false,
		provider: kindcluster.NewProvider(
			kindcluster.ProviderWithLogger(
				&logWrapper{logger.Named("kind.cluster.Provider")})),
	}
	return s
}

// GetID returns the ID of the store
func (s *Store) GetID(ctx context.Context) (string, error) {
	return "kind", nil
}

// GetContextPrefix returns the context prefix
func (s *Store) GetContextPrefix(ctx context.Context, path string) (string, error) {
	return "kind", nil
}

// VerifyKubeconfigPaths verifies the kubeconfig paths
func (s *Store) VerifyKubeconfigPaths(ctx context.Context) error {
	return nil
}

// StartSearch starts the search
func (s *Store) StartSearch(ctx context.Context, channel chan storetypes.SearchResult) {
	s.Logger.Debug("Fetching provider.List()")
	list, err := s.provider.List()
	if err != nil {
		s.Logger.Error("Error received from provider.List()", "err", err)
		channel <- storetypes.SearchResult{
			Error: err,
		}
		return
	}
	s.Logger.Debug("Received list of clusters", "list", list)
	for _, name := range list {
		channel <- storetypes.SearchResult{
			KubeconfigPath: name,
			Error:          nil,
		}
	}
}

// GetKubeconfigForPath gets the kubeconfig for the path
func (s *Store) GetKubeconfigForPath(ctx context.Context, path string, _ map[string]string) ([]byte, error) {
	s.Logger.Debug("Fetching provider.KubeConfig", "path", path, "internal", s.internal)
	kubeconfig, err := s.provider.KubeConfig(path, s.internal)
	if err != nil {
		s.Logger.Error("invalid path", "path", path)
		return nil, err
	}
	s.Logger.Debug("Received kubeconfig", "byte_len", len([]byte(kubeconfig)))
	return []byte(kubeconfig), nil
}

func main() {
	logger := hclog.Default()
	logger.SetLevel(hclog.Info)
	store := NewStore(logger.Named("ks-kind-plugin.Store"))
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins: map[string]plugin.Plugin{
			"store": &plugins.StorePlugin{Impl: store},
		},
		Logger:     logger.Named("go-plugin.Serve"),
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
