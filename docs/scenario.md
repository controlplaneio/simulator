# Scenarios

### The scenario runner - `perturb.sh`

* **Cloud provider** - is hardcoded to work with digital ocean and hobbykube
* **Node detection** - Expects exactly 3 nodes - 1 master and 2 slaves
  * `--auto-populate` uses `doctl` and a regexp to find the droplets

#### Phases

* `handle_arguments` - Argument parsing and vaildation
* `run_scenario` - Setup the scenario
* `run_test_loop` - Run Test loop

### Scenario setup

#### Validation `validate_instructions`

Before doing anything `perturb.sh` checks that it can only find bash scripts with well-known names.

Scenario scripts are found by looking for files with a `.sh` extension under the `$SCENARIO_DIR`

* `worker-any.sh` - runs on a randomly chosen slave
* `worker-1.sh` - runs on slave 1
* `worker-2.sh` - runs on slave 2
* `workers-every.sh` - runs on slaves 1 &  2
* `nodes-every.sh` - runs on master and slaves 1 and 2
* `master.sh` - runs on master
* `test.sh` - Ignored for setup

**Any other `.sh` scripts will cause an error**

#### Configure Kubernetes `run_kubectl_yaml`

1. Loop over subdirectories in `$SCENARIO_DIR`.  Subdirectories must be named after the command to be run by `kubectl`. **This is currently always only `apply`**
1. Concatenate all the `*.ya?ml` files together into a string
1. Run the concatenated string through ssh/kubectl on the master

### Run the validated shell scripts `run_scenario`

Runs each shell script on the appropriate host (see above)


### Cleanup

### Special scenarios

* `cleanup`
* `noop`

#### Scripts run always

* `no-cleanup.do` If this is present `scenario/cover-tracks.sh` will not run
* `reboot-all.do` If this is present `scenario/reboot.sh` will run on all nodes
* `reboot-workers.do` If this is present `scenario/reboot.sh` will run on master node
* `reboot-master.do` If this is present `scenario/reboot.sh` will run on worker nodes
* temporary script to copy base64 encoded`$SCENARIO_DIR/flag.txt` to `/root/flag.txt` if that file exists in `$SCENARIO_DIR`


#### Scripts only run for "normal scenarios

* temporary script to copy `$SCENARIO_DIR/challenge.txt` to `/opt/challenge.txt`
