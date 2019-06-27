package runner

import (
	"bytes"
	"github.com/pkg/errors"
	"os/user"
	"text/template"
)

var sshConfigTmplSrc = `Host {{.Hostname}}
  IdentityFile {{.KeyFilePath}}
  ProxyCommand ssh {{.User}}@{{.BastionIP}} -W %h:%p
`

const hostPrefix = "simulator_"

type SshConfig struct {
	Hostname    string
	KeyFilePath string
	User        string
	BastionIP   string
}

func CreateSshConfig(tfo TerraformOutput) (*string, error) {
	var sshConfigTmpl, err = template.New("ssh-config").Parse(sshConfigTmplSrc)
	if err != nil {
		return nil, err

	}

	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get current user for generating sshconfig")
	}

	c := SshConfig{
		Hostname:    tfo.MasterNodesPrivateIP.Value[0],
		KeyFilePath: "~/.ssh/id_rsa.pub",
		User:        u.Username,
		BastionIP:   tfo.BastionPublicIP.Value,
	}

	var buf bytes.Buffer
	err = sshConfigTmpl.Execute(&buf, c)
	if err != nil {
		return nil, errors.Wrapf(err, "Error populating ssh config template with %+v", c)
	}

	var output = string(buf.Bytes())
	return &output, nil
}
