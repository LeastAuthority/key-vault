package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/bloxapp/key-vault/utils/rpc"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Endpoints patterns
const (
	// ConfigPattern is the path pattern for config endpoint
	ConfigPattern = "config"
)

// Config contains the configuration for each mount
type Config struct {
	Network         core.Network `json:"network"`
	BeaconChainAddr string       `json:"beacon_chain_addr"`
	GenesisTime     *time.Time   `json:"genesis_time"`
}

func configPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: ConfigPattern,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWriteConfig,
				logical.UpdateOperation: b.pathWriteConfig,
				logical.ReadOperation:   b.pathReadConfig,
			},
			HelpSynopsis:    "Configure the Vault Ethereum plugin.",
			HelpDescription: "Configure the Vault Ethereum plugin.",
			Fields: map[string]*framework.FieldSchema{
				"network": {
					Type: framework.TypeString,
					Description: `Ethereum network - can be one of the following values:
					launchtest - Launch Test Network
					test 	   - Goerli Test Network`,
					AllowedValues: []interface{}{
						string(core.TestNetwork),
						string(core.LaunchTestNetwork),
					},
				},
				"beacon_chain_addr": {
					Type:        framework.TypeString,
					Description: `Beacon chain server RPC address`,
				},
			},
		},
	}
}

// pathWriteConfig is the write config path handler
func (b *backend) pathWriteConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	network := data.Get("network").(string)
	beaconChainAddr := data.Get("beacon_chain_addr").(string)

	// Load genesis time
	genesisTime, err := b.loadGenesisTime(ctx, beaconChainAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load genesis time")
	}

	configBundle := Config{
		Network:         core.NetworkFromString(network),
		BeaconChainAddr: beaconChainAddr,
		GenesisTime:     genesisTime,
	}

	// Create storage entry
	entry, err := logical.StorageEntryJSON("config", configBundle)
	if err != nil {
		return nil, err
	}

	// Store config
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"network":           configBundle.Network,
			"beacon_chain_addr": configBundle.BeaconChainAddr,
			"genesis_time":      configBundle.GenesisTime,
		},
	}, nil
}

// pathReadConfig is the read config path handler
func (b *backend) pathReadConfig(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	configBundle, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if configBundle == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"network":           configBundle.Network,
			"beacon_chain_addr": configBundle.BeaconChainAddr,
			"genesis_time":      configBundle.GenesisTime,
		},
	}, nil
}

// readConfig returns the configuration for this PluginBackend.
func (b *backend) readConfig(ctx context.Context, s logical.Storage) (*Config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("the plugin has not been configured yet")
	}

	var result Config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, errors.Wrap(err, "error reading configuration")
		}
	}

	return &result, nil
}

func (b *backend) configured(ctx context.Context, req *logical.Request) (*Config, error) {
	config, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (b *backend) loadGenesisTime(ctx context.Context, beaconChainAddr string) (*time.Time, error) {
	// Create RPC client
	client, err := rpc.Connect(beaconChainAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to beacon server %s", beaconChainAddr)
	}
	nodeClient := ethpb.NewNodeClient(client)
	defer client.Close()

	// Get genesis time
	genesisResp, err := nodeClient.GetGenesis(ctx, &types.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get genesis time")
	}

	// Parse genesis
	genesisTime, err := types.TimestampFromProto(genesisResp.GetGenesisTime())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse genesis time")
	}

	return &genesisTime, nil
}
