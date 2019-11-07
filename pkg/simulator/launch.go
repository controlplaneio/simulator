package simulator

import (
	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Launch runs perturb.sh to setup a scenario with the supplied `id` assuming
// the infrastructure has been created.  Returns an error if the infrastructure
// is not ready or something goes wrong
func Launch(logger *zap.SugaredLogger, tfDir, scenariosDir, bucketName, id, attackTag string) error {
	logger.Debugf("Loading scenario manifest from %s", scenariosDir)
	manifest, err := scenario.LoadManifest(scenariosDir)
	if err != nil {
		return errors.Wrap(err, "Error loading scenario manifest file")
	}

	logger.Debugf("Checking manifest contains %s", id)
	if !manifest.Contains(id) {
		return errors.Errorf("Scenario not found: %s", id)
	}

	logger.Debugf("Checking status of infrastructure")
	tfo, err := Status(logger, tfDir, bucketName, attackTag)
	if !tfo.IsUsable() {
		return errors.Errorf("No infrastructure, please run simulator infra create")
	}
	logger.Debug(tfo)

	logger.Infof("Finding details of scenario %s", id)
	s := manifest.Find(id)
	logger.Debug(s)

	logger.Debugf(
		"Making options to pass to perturb from terraorm output and scnenario")
	po := MakePerturbOptions(*tfo, s.Path)
	logger.Debug(po)

	logger.Debug("Regenerating SSH config")
	cfg, err := tfo.ToSSHConfig()
	if err != nil {
		return errors.Wrap(err, "Error templating SSH config")
	}

	logger.Info("Updating SSH config")
	err = ssh.EnsureSSHConfig(*cfg)
	if err != nil {
		return errors.Wrap(err, "Error writing SSH config")
	}

	bastion := tfo.BastionPublicIP.Value
	logger.Infof("Keyscanning %s and updating known hosts", bastion)
	err = ssh.EnsureKnownHosts(bastion)
	if err != nil {
		return errors.Wrapf(err, "Error updating known hosts for bastion: %s",
			bastion)
	}

	logger.Infof("Setting up the \"%s\" scenario on the cluster", s.DisplayName)
	_, err = Perturb(&po)
	if err != nil {
		return errors.Wrapf(err, "Error running perturb with %#v", po)
	}

	return nil
}
