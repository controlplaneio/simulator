package runner

// TfVars struct representing the input variables for terraform to create the infrastructure
type TfVars struct {
	PublicKey  string
	AccessCIDR string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey string, accessCIDR string) TfVars {
	return TfVars{
		PublicKey:  publicKey,
		AccessCIDR: accessCIDR,
	}
}

var tmpl = `
access_key={{.PublicKey}}
access_cidr={{.AccessCIDR}}
`

func (tfv *TfVars) String() string {
	return "access_key = \"" + tfv.PublicKey + "\"\n" + "access_cidr = \"" + tfv.AccessCIDR + "\"\n"
}

// EnsureTfVarsFile writes an tfvars file if one hasnt already been made
func EnsureTfVarsFile(tfDir, publicKey, accessCIDR string) error {
	filename := tfDir + "/settings/bastion.tfVars"
	tfv := NewTfVars(publicKey, accessCIDR)

	return EnsureFile(filename, tfv.String())
}
