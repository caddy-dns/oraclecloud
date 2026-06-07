package oraclecloud

import (
	"testing"

	libdnsoraclecloud "github.com/libdns/oraclecloud"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestUnmarshalCaddyfileFullConfig(t *testing.T) {
	t.Parallel()

	input := `
oraclecloud {
	auth config_file
	config_file ~/.oci/config
	config_profile DEFAULT
	private_key -----BEGIN
	private_key_path ~/.oci/oci_api_key.pem
	private_key_passphrase secret
	tenancy_ocid ocid1.tenancy.oc1..aaaa
	user_ocid ocid1.user.oc1..bbbb
	fingerprint 12:34:56
	region us-ashburn-1
	scope PRIVATE
	view_id ocid1.view.oc1..cccc
	compartment_id ocid1.compartment.oc1..dddd
}
`

	var provider Provider
	if err := provider.UnmarshalCaddyfile(caddyfile.NewTestDispenser(input)); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	got := provider.Provider
	if got == nil {
		t.Fatal("expected provider to be initialized")
	}

	assertEqual(t, got.Auth, "config_file", "auth")
	assertEqual(t, got.ConfigFile, "~/.oci/config", "config_file")
	assertEqual(t, got.ConfigProfile, "DEFAULT", "config_profile")
	assertEqual(t, got.PrivateKey, "-----BEGIN", "private_key")
	assertEqual(t, got.PrivateKeyPath, "~/.oci/oci_api_key.pem", "private_key_path")
	assertEqual(t, got.PrivateKeyPassphrase, "secret", "private_key_passphrase")
	assertEqual(t, got.TenancyOCID, "ocid1.tenancy.oc1..aaaa", "tenancy_ocid")
	assertEqual(t, got.UserOCID, "ocid1.user.oc1..bbbb", "user_ocid")
	assertEqual(t, got.Fingerprint, "12:34:56", "fingerprint")
	assertEqual(t, got.Region, "us-ashburn-1", "region")
	assertEqual(t, got.Scope, "PRIVATE", "scope")
	assertEqual(t, got.ViewID, "ocid1.view.oc1..cccc", "view_id")
	assertEqual(t, got.CompartmentID, "ocid1.compartment.oc1..dddd", "compartment_id")
}

func TestUnmarshalCaddyfileAuthShorthand(t *testing.T) {
	t.Parallel()

	input := `
oraclecloud config_file {
	config_file ~/.oci/config
}
`

	var provider Provider
	if err := provider.UnmarshalCaddyfile(caddyfile.NewTestDispenser(input)); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	assertEqual(t, provider.Provider.Auth, "config_file", "auth")
	assertEqual(t, provider.Provider.ConfigFile, "~/.oci/config", "config_file")
}

func TestUnmarshalCaddyfileRejectsDuplicateField(t *testing.T) {
	t.Parallel()

	input := `
oraclecloud {
	auth auto
	auth config_file
}
`

	var provider Provider
	if err := provider.UnmarshalCaddyfile(caddyfile.NewTestDispenser(input)); err == nil {
		t.Fatal("expected duplicate field error")
	}
}

func TestReplacePlaceholders(t *testing.T) {
	t.Setenv("OCI_CONFIG_FILE", "/Users/test/.oci/config")
	t.Setenv("OCI_REGION", "eu-frankfurt-1")

	provider := Provider{
		Provider: &libdnsoraclecloud.Provider{
			ConfigFile: "{env.OCI_CONFIG_FILE}",
			Region:     "{env.OCI_REGION}",
		},
	}

	provider.replacePlaceholders(caddy.NewReplacer())

	assertEqual(t, provider.Provider.ConfigFile, "/Users/test/.oci/config", "config_file")
	assertEqual(t, provider.Provider.Region, "eu-frankfurt-1", "region")
}

func assertEqual(t *testing.T, got, want, field string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %q, want %q", field, got, want)
	}
}
