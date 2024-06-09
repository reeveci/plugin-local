# Reeve CI / CD - Local Plugin

This is a [Reeve](https://github.com/reeveci/reeve) plugin for providing pipeline environment variables from local storage.

## Configuration

This plugin supports configuration via the WebUI plugin.

Configuration is also supported via the CLI API:

```sh
reeve-cli --url <server-url> --usage
reeve-cli --url <server-url> local set-variable <name> <value>
reeve-cli --url <server-url> local set-secret <name> <value>
```

Encryption takes place on the server, so make sure to use a secure connection between reeve-cli and the server. That is, use TLS with a valid certificate and do not set the `--insecure` flag.

### Settings

Settings can be provided to the plugin through environment variables set to the reeve server.

Settings for this plugin should be prefixed by `REEVE_PLUGIN_LOCAL_`.

Settings may also be shared between plugins by prefixing them with `REEVE_SHARED_` instead.

- `ENABLED` - `true` enables this plugin
- `CONFIG_PATH` (required) - Path to where configuration should be stored on disk.
- `SECRET_KEY` (required) - Passphrase for encrypting secrets
- `PRIORITY` (default 1) - Priority of all variables returned by this plugin
