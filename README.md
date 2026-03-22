# Oracle Cloud Infrastructure DNS Provider for Caddy-DNS

This package provides an OCI DNS provider module for [Caddy](https://github.com/caddyserver/caddy) backed by [`github.com/Djelibeybi/libdns-oraclecloud`](https://github.com/Djelibeybi/libdns-oraclecloud).

## Caddy module name

```text
dns.providers.oraclecloud
```

## Building

Build Caddy with:

```sh
xcaddy build --with github.com/Djelibeybi/caddy-dns-oraclecloud
```

## Configuration

The wrapper exposes the underlying `libdns-oraclecloud` provider fields directly. Supported auth modes are:

- `auto` or empty
- `api_key`
- `config_file`
- `environment`
- `instance_principal`

`auto` is the default and follows the provider's own precedence order:

1. Explicit API key fields on the module
2. OCI config file credentials
3. `OCI_CLI_*` environment variables

For private zones accessed by name, Oracle Cloud requires both `scope PRIVATE` and `view_id <ocid>`.

### JSON example

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "oraclecloud",
        "auth": "config_file",
        "config_file": "{env.OCI_CLI_CONFIG_FILE}",
        "config_profile": "{env.OCI_CLI_PROFILE}",
        "region": "{env.OCI_CLI_REGION}"
      }
    }
  }
}
```

### Caddyfile examples

Using an OCI config file:

```caddyfile
tls {
	dns oraclecloud {
		auth config_file
		config_file {env.OCI_CLI_CONFIG_FILE}
		config_profile {env.OCI_CLI_PROFILE}
		region {env.OCI_CLI_REGION}
	}
}
```

Using direct API key fields:

```caddyfile
tls {
	dns oraclecloud {
		auth api_key
		tenancy_ocid {env.OCI_CLI_TENANCY}
		user_ocid {env.OCI_CLI_USER}
		fingerprint {env.OCI_CLI_FINGERPRINT}
		private_key_path {env.OCI_CLI_KEY_FILE}
		private_key_passphrase {env.OCI_CLI_PASSPHRASE}
		region {env.OCI_CLI_REGION}
	}
}
```

Using a private zone by name:

```caddyfile
tls {
	dns oraclecloud {
		auth config_file
		config_file {env.OCI_CLI_CONFIG_FILE}
		config_profile DEFAULT
		scope PRIVATE
		view_id ocid1.dnsview.oc1..exampleuniqueID
	}
}
```

## Available Caddyfile options

```text
oraclecloud [<auth_mode>] {
    auth <auth_mode>
    config_file <path>
    config_profile <profile>
    private_key <pem>
    private_key_path <path>
    private_key_passphrase <passphrase>
    tenancy_ocid <ocid>
    user_ocid <ocid>
    fingerprint <fingerprint>
    region <region>
    scope <GLOBAL|PRIVATE>
    view_id <ocid>
    compartment_id <ocid>
}
```

All string fields support Caddy placeholders such as `{env.OCI_CLI_REGION}`.

## Notes

- `compartment_id` is only needed if you use the underlying provider's `ListZones` capability; it is not required for standard ACME DNS challenge flows.
- For authentication details and provider behavior, see the [`libdns-oraclecloud` README](https://github.com/Djelibeybi/libdns-oraclecloud).

## Versioning

This repository uses Release Please with conventional commits to automate changelog entries, release PRs, Git tags, and GitHub Releases.

- `fix:` commits produce patch releases
- `feat:` commits produce minor releases
- `feat!:` or commits with a `BREAKING CHANGE:` footer produce major releases
