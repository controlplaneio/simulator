package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/pkg/errors"
)

// Launch runs perturb.sh to setup a scenario with the supplied `id` assuming the infrastructure has been created.
// Returns an error if the infrastructure is not ready or something goes wrong
func Launch(tfDir, scenariosDir, bucketName, id string) error {
	manifest, err := scenario.LoadManifest(scenariosDir)
	if err != nil {
		return err
	}

	if !manifest.Contains(id) {
		return errors.Errorf("scenario %s not found", id)
	}

	tfo, err := Status(tfDir, bucketName)
	if !tfo.IsUsable() {
		return errors.Errorf("No infrastructure, please run simulator infra create")
	}

	scenarioPath := manifest.Find(id).Path

	po := MakePerturbOptions(*tfo, scenarioPath)
	c, err := tfo.ToSSHConfig()
	if err != nil {
		return err
	}

	cp, err := util.ExpandTilde("~/.ssh/config")
	if err != nil {
		return err
	}

	err = util.OverwriteFile(*cp, *c)
	if err != nil {
		return err
	}

	bastion := tfo.BastionPublicIP.Value
	err = util.UpdateKnownHosts(bastion)
	if err != nil {
		return errors.Wrapf(err, "Error updating known hosts for bastion: %s", bastion)
	}

	_, err = Perturb(&po)
	if err != nil {
		return errors.Wrapf(err, "Error running perturb with %#v", po)
	}

	return nil
}
