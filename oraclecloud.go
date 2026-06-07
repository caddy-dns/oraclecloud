package oraclecloud

import (
	"fmt"

	libdnsoraclecloud "github.com/libdns/oraclecloud"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the libdns Oracle Cloud provider as a Caddy module.
type Provider struct{ *libdnsoraclecloud.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.oraclecloud",
		New: func() caddy.Module { return &Provider{Provider: new(libdnsoraclecloud.Provider)} },
	}
}

// Provision resolves placeholders before the provider is used.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.replacePlaceholders(caddy.NewReplacer())
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	oraclecloud [<auth_mode>] {
//	    auth <auth_mode>
//	    config_file <path>
//	    config_profile <profile>
//	    private_key <pem>
//	    private_key_path <path>
//	    private_key_passphrase <passphrase>
//	    tenancy_ocid <ocid>
//	    user_ocid <ocid>
//	    fingerprint <fingerprint>
//	    region <region>
//	    scope <GLOBAL|PRIVATE>
//	    view_id <ocid>
//	    compartment_id <ocid>
//	}
//
// If auth is omitted, the underlying libdns provider defaults to "auto".
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	provider := p.ensureProvider()

	for d.Next() {
		if d.NextArg() {
			if provider.Auth != "" {
				return d.Err("auth already set")
			}
			provider.Auth = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "auth":
				if err := parseSingleArg(d, &provider.Auth, "auth"); err != nil {
					return err
				}
			case "config_file":
				if err := parseSingleArg(d, &provider.ConfigFile, "config_file"); err != nil {
					return err
				}
			case "config_profile":
				if err := parseSingleArg(d, &provider.ConfigProfile, "config_profile"); err != nil {
					return err
				}
			case "private_key":
				if err := parseSingleArg(d, &provider.PrivateKey, "private_key"); err != nil {
					return err
				}
			case "private_key_path":
				if err := parseSingleArg(d, &provider.PrivateKeyPath, "private_key_path"); err != nil {
					return err
				}
			case "private_key_passphrase":
				if err := parseSingleArg(d, &provider.PrivateKeyPassphrase, "private_key_passphrase"); err != nil {
					return err
				}
			case "tenancy_ocid":
				if err := parseSingleArg(d, &provider.TenancyOCID, "tenancy_ocid"); err != nil {
					return err
				}
			case "user_ocid":
				if err := parseSingleArg(d, &provider.UserOCID, "user_ocid"); err != nil {
					return err
				}
			case "fingerprint":
				if err := parseSingleArg(d, &provider.Fingerprint, "fingerprint"); err != nil {
					return err
				}
			case "region":
				if err := parseSingleArg(d, &provider.Region, "region"); err != nil {
					return err
				}
			case "scope":
				if err := parseSingleArg(d, &provider.Scope, "scope"); err != nil {
					return err
				}
			case "view_id":
				if err := parseSingleArg(d, &provider.ViewID, "view_id"); err != nil {
					return err
				}
			case "compartment_id":
				if err := parseSingleArg(d, &provider.CompartmentID, "compartment_id"); err != nil {
					return err
				}
			default:
				return d.Errf("unrecognized subdirective %q", d.Val())
			}
		}
	}

	return nil
}

func (p *Provider) ensureProvider() *libdnsoraclecloud.Provider {
	if p.Provider == nil {
		p.Provider = new(libdnsoraclecloud.Provider)
	}
	return p.Provider
}

func (p *Provider) replacePlaceholders(repl *caddy.Replacer) {
	provider := p.ensureProvider()
	provider.Auth = repl.ReplaceAll(provider.Auth, "")
	provider.ConfigFile = repl.ReplaceAll(provider.ConfigFile, "")
	provider.ConfigProfile = repl.ReplaceAll(provider.ConfigProfile, "")
	provider.PrivateKey = repl.ReplaceAll(provider.PrivateKey, "")
	provider.PrivateKeyPath = repl.ReplaceAll(provider.PrivateKeyPath, "")
	provider.PrivateKeyPassphrase = repl.ReplaceAll(provider.PrivateKeyPassphrase, "")
	provider.TenancyOCID = repl.ReplaceAll(provider.TenancyOCID, "")
	provider.UserOCID = repl.ReplaceAll(provider.UserOCID, "")
	provider.Fingerprint = repl.ReplaceAll(provider.Fingerprint, "")
	provider.Region = repl.ReplaceAll(provider.Region, "")
	provider.Scope = repl.ReplaceAll(provider.Scope, "")
	provider.ViewID = repl.ReplaceAll(provider.ViewID, "")
	provider.CompartmentID = repl.ReplaceAll(provider.CompartmentID, "")
}

func parseSingleArg(d *caddyfile.Dispenser, target *string, field string) error {
	if *target != "" {
		return d.Errf("%s already set", field)
	}
	if !d.NextArg() {
		return d.ArgErr()
	}
	*target = d.Val()
	if d.NextArg() {
		return fmt.Errorf("%s %w", field, d.ArgErr())
	}
	return nil
}

// Interface guards.
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
