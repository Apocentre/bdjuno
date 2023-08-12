package main

import (
	"selfchain/x/migration"
	migrationtypes "selfchain/x/migration/types"
	selfvestingtypes "selfchain/x/selfvesting/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	shashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	"github.com/forbole/juno/v4/cmd"
	initcmd "github.com/forbole/juno/v4/cmd/init"
	parsetypes "github.com/forbole/juno/v4/cmd/parse/types"
	startcmd "github.com/forbole/juno/v4/cmd/start"
	"github.com/forbole/juno/v4/modules/messages"

	migratecmd "github.com/forbole/bdjuno/v4/cmd/migrate"
	parsecmd "github.com/forbole/bdjuno/v4/cmd/parse"

	"github.com/forbole/bdjuno/v4/types/config"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules"
)

func main() {
	initCfg := initcmd.NewConfig().
		WithConfigCreator(config.Creator)

	parseCfg := parsetypes.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("bdjuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := cmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		cmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(cfg.GetName(), cfg.GetParseConfig()),
		startcmd.NewStartCmd(cfg.GetParseConfig()),
	)

	executor := cmd.PrepareRootCmd(cfg.GetName(), rootCmd)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []module.BasicManager {
	modules := map[string]module.AppModuleBasic {
		authtypes.ModuleName: auth.AppModuleBasic{},
		stakingtypes.ModuleName: staking.AppModuleBasic{},
		distributiontypes.ModuleName: distribution.AppModuleBasic {},
		banktypes.ModuleName: bank.AppModuleBasic{},
		govtypes.ModuleName: gov.AppModuleBasic{},
		ibctypes.ModuleName: ibc.AppModuleBasic{},
		ibctransfertypes.ModuleName: ibctransfer.AppModuleBasic{},
		shashingtypes.ModuleName: slashing.AppModuleBasic{},
		migrationtypes.ModuleName: migration.AppModuleBasic{},
		selfvestingtypes.ModuleName: migration.AppModuleBasic{},
	}
	
	return []module.BasicManager {modules}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
	)
}
