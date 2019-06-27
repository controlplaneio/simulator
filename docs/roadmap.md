# Roadmap

## Milestone 1 - Open sourcing 4 scenarios

## Milestone 1 - Epic 1 - Separate standalone simulator and launch flow

### Immediate future - for launch demo

- ~~Remove digital ocean checks~~
- ~~Change terraform to create 2 worker nodes~~
- Automate ssh config for connecting to instances via bastion

### Backlog

- Split digital ocean automatic master / worker population out of perturb.sh into a helper
- Logging
  - Standardise logging messsages and levels
  - Make logging configurable
  - Write logs from various components to well known (configurable) files
- Attack container
  - nmap
  - kubectl
- Documentation
  - godoc comments for all libraries
- Testing
 - CI /CD
 - Code coverage
 - go vet
- Design for supporting multiple users in the same AWS Account(?)

## Milestone 1  - Epic 2
- Validate scenarios post running the scripts
  - drop "state" into nodes - to track whats been done on them(?)
- Evaluate user progress

## Future Milestones - Post initial open-source

- Support configurable worker node count
- Support configurable master node count - multimaster
- Reset cluster
- Multi cloud
