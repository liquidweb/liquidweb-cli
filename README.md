# liquidweb-cli
Official command line interface for the LiquidWeb API
```
CLI interface for LiquidWeb.

Command line interface for interacting with LiquidWeb services via
LiquidWebs Public API.

If this is your first time running, you will need to setup at least
one auth context. An auth context contains authentication data for
accessing your LiquidWeb account. As such one auth context represents
one LiquidWeb account. You can have multiple auth contexts defined.

To setup your first auth context, you run 'auth init'. For further
information on auth contexts, be sure to checkout 'help auth' for a
list of capabilities.

As always, consult the various subcommands for specific features and
capabilities.

Usage:
  liquidweb-cli [command]

Available Commands:
  auth        authentication actions
  cloud       Interact with LiquidWebs Cloud platform.
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.liquidweb-cli.yaml)
  -h, --help            help for liquidweb-cli

Use "liquidweb-cli [command] --help" for more information about a command.
```
