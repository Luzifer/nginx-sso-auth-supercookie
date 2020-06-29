module nginx-sso-auth-supercookie

go 1.14

// MUST use the SAME code, so may not download the code itself
// but MUST reference to the original source code already compiled
replace github.com/Luzifer/nginx-sso => ../nginx-sso

require (
	github.com/Luzifer/go_helpers/v2 v2.9.1 // Keep exactly the same version as nginx-sso
	github.com/Luzifer/nginx-sso v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.2.4 // Keep exactly the same version as nginx-sso
)
