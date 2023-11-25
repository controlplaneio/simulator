FROM controlplane/simulator:dev

COPY --chown=ubuntu:ubuntu packer packer
COPY --chown=ubuntu:ubuntu terraform terraform
COPY --chown=ubuntu:ubuntu ansible ansible

RUN cd packer && packer init bastion.pkr.hcl && packer init k8s.pkr.hcl
RUN cd terraform/workspaces/simulator && terraform init -backend=false
