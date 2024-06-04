package main

import (
	"selfchain/x/migration"
	migrationtypes "selfchain/x/migration/types"
	"selfchain/x/selfvesting"
	selfvestingtypes "selfchain/x/selfvesting/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/forbole/juno/v5/cmd"
	initcmd "github.com/forbole/juno/v5/cmd/init"
	parsetypes "github.com/forbole/juno/v5/cmd/parse/types"
	startcmd "github.com/forbole/juno/v5/cmd/start"
	"github.com/forbole/juno/v5/modules/messages"

	migratecmd "github.com/forbole/callisto/v4/cmd/migrate"
	parsecmd "github.com/forbole/callisto/v4/cmd/parse"

	"github.com/forbole/callisto/v4/types/config"

	"cosmossdk.io/simapp"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules"
)

func main() {
	initCfg := initcmd.NewConfig().
		WithConfigCreator(config.Creator)

	parseCfg := parsetypes.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("callisto").
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
		selfvestingtypes.ModuleName: selfvesting.AppModuleBasic{},
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
