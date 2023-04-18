package evm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	icqkeeper "github.com/QuadrateOrg/core/x/interchainqueries/keeper"
	icqtypes "github.com/QuadrateOrg/core/x/interchainqueries/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ek evmkeeper.Keeper,
	fk feemarketkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[evmtypes.ModuleName] = mm.Modules[evmtypes.ModuleName].ConsensusVersion()
		fromVM[feemarkettypes.ModuleName] = mm.Modules[feemarkettypes.ModuleName].ConsensusVersion()

		var params EvmUpgradeParams
		err := json.Unmarshal([]byte(plan.Info), &params)
		if err != nil {
			panic(err)
		}

		ek.SetParams(ctx, params.Evm)
		fk.SetParams(ctx, params.FeeMarket)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func CreateUpgradeHandlerICQ(
	mm *module.Manager,
	configurator module.Configurator,
	icqKeeper icqkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		// todo: FIXME
		icqKeeper.SetParams(ctx, icqtypes.DefaultParams())

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}
