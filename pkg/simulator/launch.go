package simulator

import (
	"fmt"
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
		return errors.Errorf("No infrastructure, please run simulator infra create:\n %#v", tfo)
	}

	scenarioPath := manifest.Find(id).Path

	po := MakePerturbOptions(*tfo, scenarioPath)
	fmt.Println("Converted usable terraform output into perturb options")
	fmt.Printf("%#v", po)
	c, err := tfo.ToSSHConfig()
	if err != nil {
		return err
	}

	cp, err := util.ExpandTilde("~/.ssh/config")
	if err != nil {
		return err
	}

	// BUG: (rem) doesnt work when SSH config doesnt exist in docker container
	written, err := util.EnsureFile(*cp, *c)
	if err != nil {
		return err
	}

	if !written {
		fmt.Printf("Please add the following lines to your ssh config\n---\n%s\n---\n", *c)
	}

	_, err = Perturb(&po)
	if err != nil {
		return errors.Wrapf(err, "Error running perturb with %#v", po)
	}

	return nil
}
