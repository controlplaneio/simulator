package simulator

import (
	"strings"

	"github.com/controlplaneio/simulator-standalone/pkg/scenario"
	"github.com/controlplaneio/simulator-standalone/pkg/ssh"
	"github.com/pkg/errors"
)

// Launch runs perturb.sh to setup a scenario with the supplied `id` assuming
// the infrastructure has been created.  Returns an error if the infrastructure
// is not ready or something goes wrong
func (s *Simulator) Launch() error {
	s.Logger.Debugf("Loading scenario manifest from %s", s.ScenariosDir)
	manifest, err := scenario.LoadManifest(s.ScenariosDir)
	if err != nil {
		return errors.Wrap(err, "Error loading scenario manifest file")
	}

	s.Logger.Debugf("Checking manifest contains %s", s.ScenarioID)
	if !manifest.Contains(s.ScenarioID) {
		return errors.Errorf("Scenario not found: %s", s.ScenarioID)
	}

	s.Logger.Debugf("Checking status of infrastructure")

	simulator := NewSimulator(
		WithLogger(s.Logger),
		WithTfDir(s.TfDir),
		WithScenariosDir(s.ScenariosDir),
		WithAttackTag(s.AttackTag),
		WithBucketName(s.BucketName),
		WithTfVarsDir(s.TfVarsDir))

	tfo, err := simulator.Status()

	if !tfo.IsUsable() {
		return errors.Errorf("No infrastructure, please run simulator infra create")
	}
	s.Logger.Debug(tfo)

	s.Logger.Infof("Finding details of scenario %s", s.ScenarioID)
	sID := manifest.Find(s.ScenarioID)
	s.Logger.Debug(sID)

	s.Logger.Debugf(
		"Making options to pass to perturb from terraorm output and scnenario")
	po := MakePerturbOptions(*tfo, sID.Path)
	s.Logger.Debug(po)

	s.Logger.Debug("Regenerating SSH config")
	cfg, err := tfo.ToSSHConfig()
	if err != nil {
		return errors.Wrap(err, "Error templating SSH config")
	}

	s.Logger.Info("Updating SSH config")
	err = ssh.EnsureSSHConfig(*cfg)
	if err != nil {
		return errors.Wrap(err, "Error writing SSH config")
	}

	bastion := tfo.BastionPublicIP.Value
	s.Logger.Infof("Keyscanning %s and updating known hosts", bastion)
	err = ssh.EnsureKnownHosts(bastion)
	if err != nil {
		return errors.Wrapf(err, "Error updating known hosts for bastion: %s",
			bastion)
	}

	s.Logger.Infof("Setting up the \"%s\" scenario on the cluster", sID.DisplayName)
	_, err = Perturb(&po)
	if err != nil {
		if strings.Contains(err.Error(), "exit status 103") {
			s.Logger.Error("Scenario clash error from perturb.sh")
		} else {
			return errors.Wrapf(err, "Error running perturb with %#v", po)
		}
	}

	return nil
}
