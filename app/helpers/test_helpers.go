package helpers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	quadrateapp "github.com/0xknstntn/quadrate/app"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "quadrate-app"
)

// DefaultConsensusParams defines the default Tendermint consensus params used
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type EmptyAppOptions struct{}

func (EmptyAppOptions) Get(o string) interface{} { return nil }

func Setup(t *testing.T, chainId string, isCheckTx bool, invCheckPeriod uint) *quadrateapp.QuadrateApp {
	t.Helper()

	app, genesisState := setup(!isCheckTx, invCheckPeriod)
	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				ChainId:         chainId,
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func setup(withGenesis bool, invCheckPeriod uint) (*quadrateapp.QuadrateApp, quadrateapp.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := quadrateapp.MakeEncodingConfig()
	app := quadrateapp.NewQuadrateApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		quadrateapp.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		EmptyAppOptions{},
	)
	if withGenesis {
		return app, quadrateapp.NewDefaultGenesisState()
	}

	return app, quadrateapp.GenesisState{}
}
