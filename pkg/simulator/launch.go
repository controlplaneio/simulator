package simulator

import (
	"strings"

	"github.com/kubernetes-simulator/simulator/pkg/scenario"
	"github.com/kubernetes-simulator/simulator/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Launch runs perturb.sh to setup a scenario with the supplied `id` assuming
// the infrastructure has been created.  Returns an error if the infrastructure
// is not ready or something goes wrong
func (s *Simulator) Launch() error {
	s.Logger.WithFields(logrus.Fields{
		"ScenariosDir": s.ScenariosDir,
	}).Debug("Loading scenario manifest")
	manifest, err := scenario.LoadManifest(s.ScenariosDir)
	if err != nil {
		return errors.Wrap(err, "Error loading scenario manifest file")
	}

	s.Logger.WithFields(logrus.Fields{
		"ScenarioID": s.ScenarioID,
	}).Debug("Checking manifest contains scenario")
	if !manifest.Contains(s.ScenarioID) {
		return errors.Errorf("Scenario not found: %s", s.ScenarioID)
	}

	s.Logger.Debug("Checking status of infrastructure")
	tfo, _ := s.Status()

	if !tfo.IsUsable() {
		return errors.Errorf("No infrastructure, please run simulator infra create")
	}
	s.Logger.Debug(tfo)

	s.Logger.WithFields(logrus.Fields{
		"ScenarioID": s.ScenarioID,
	}).Infof("Finding details of scenario")
	foundScenario := manifest.Find(s.ScenarioID)
	s.Logger.Debug(foundScenario)

	s.Logger.Debug(
		"Making options to pass to perturb from terraorm output and scnenario")
	po := MakePerturbOptions(*tfo, foundScenario.Path)
	s.Logger.Debug(po)

	s.Logger.Debug("Regenerating SSH config")
	cfg, err := tfo.ToSSHConfig()
	if err != nil {
		return errors.Wrap(err, "Error templating SSH config")
	}

	s.Logger.Info("Updating SSH config")
	err = s.SSHStateProvider.SaveSSHConfig(*cfg)
	if err != nil {
		return errors.Wrap(err, "Error writing SSH config")
	}

	bastion := tfo.BastionPublicIP.Value
	s.Logger.WithFields(logrus.Fields{
		"BastionIP": bastion,
	}).Info("Keyscanning bastion and updating known hosts")
	err = ssh.EnsureKnownHosts(bastion)
	if err != nil {
		return errors.Wrapf(err, "Error updating known hosts for bastion: %s",
			bastion)
	}

	s.Logger.WithFields(logrus.Fields{
		"Scenario": foundScenario.DisplayName,
	}).Info("Setting up the scenario on the cluster")
	_, err = Perturb(&po, s.Logger)
	if err != nil {
		if strings.Contains(err.Error(), "exit status 103") {
			s.Logger.Error("Scenario clash error from perturb.sh")
		} else {
			return errors.Wrapf(err, "Error running perturb with %#v", po)
		}
	}

	return nil
}
