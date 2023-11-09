# Simulator Container Images

To build the Simulator Container Images run `make simulator-image` to build the two images.

| Name                          | Description                                                                    |
|-------------------------------|--------------------------------------------------------------------------------|
| controlplane/simulator:latest | The complete image, bundling the required tools, and all of the configuration. |
| controlplane/simulator:dev    | The development image, bundling the required tools.                            |

The following tools are bundled into both images.

* Ansible
* Packer
* Terraform
* The Simulator controlplane CLI

This allows users to execute the various commands without having to install the required tools locally and managing
compatible versions. The Simulator CLI will run the image and execute the specified command within the image.

The following directories will be bind mounted into the container at runtime.

| Name                    | Description                                                                              |
|-------------------------|------------------------------------------------------------------------------------------|
| $HOME/.aws              | The users AWS configuration directory for access AWS credentials.                        |
| $HOME/.simulator/admin  | The directory where Simulator will write the admin ssh bundle and ansible configuration. |
| $HOME/.simulator/player | The directory where Simulator will write the player ssh bundle.                          |

[//]: # (TODO: Use the same configuration directory from SIMULATOR_DIR for the configuration?)
