
# Project Vision - Why did I start this?

I started this project as I needed a way to provide hands-on Container and Kubernetes security training to engineers and could not find an appropriate solution in the marketplace (although KataCoda came close). I have had a lot of success in the past with building bespoke training environments as I have found that the engagement and retention obtained from these systems far exceeds that of slideware or reading books. By simulating a scenario, engineers are allowed to understand the strengths and weaknesses of a system and mitigate the latter in a real world simulation, using real world tools.

The ultimate vision of this project is to have a multi-player environment where we can pit blue team, red team and forensics team members against each other. The red team member would be presented with a security problem for investigation and to demonstrate that they can exploit it. Upon completion, the exercise would move to the forensics team member who would examine logs files and the system to gain an understanding of the issue and timeline of the incident. This would be followed by the blue team member who would be tasked with mitigating the issue. By creating scenarios incorporating the multiple security teams we create an environment where teams can increase collaboration.

## Intended Audience

Container security is a broad topic. The main focus of this training platform is to educate red and blue team members in the ways of exploiting and hardening a cluster.  Hence this is more beneficial to blue team members tasked with securing the Kubernetes platform itself. There are clearly various other team members tasked with securing an application running in Kubernetes; DevSecOps and application developers to name a few. Exercises may be useful to these communities however the main target audience is security
engineers tasked with securing the platform. Therefore the exercises aim to ensure standard security mitigations have been configured appropriately. These basic Kubernetes settings should already have been taken care of by platform engineering before an application developer uses the system.

## Project Phases

### Phase 1

The first phase of the project was to create an architecture that would allow a Kubernetes system to be provisioned within a secured environment, facilitating multiple security exercises to be loaded. The student would then be provided with a mechanism to connect to the Kubernetes system from multiple vantage points - shell in a container, access to a node or external cluster access only. The shell access into a container is to simulate a situation where a deployed application has been compromised.

### Phase 2

The second phase extends phase one capabilities by providing a mechanism to score the exercises and validate mitigations have been accomplished. Additionally a set of 25 exercises will be provided focussing mainly on major mitigations and protections that should be in place. We also aim to create a local training element where the system can be deployed to a laptop using KIND rather than a remote cloud provider. This would not be a full Kubernetes installation but would provide sufficient functionality as to allow the majority of exercises to function.  

### Phase 3

The third phase adds the multiple player element of the system. Allowing the training progress to be saved and handed over to another team member. This would enable the ultimate vision of having red team, blue team and forensics team members all working together on an exercise.

### Future Phases

I would like to add a further element of gamification to the system. Showing statistics on how well or quickly a student completed an exercise, adding scores to a leaderboard along with history to show a student’s progress.  This could potentially be extended, using the system to host hack-a-thon type wargames at conferences where multiple teams of blue and red teams compete against one another.

In previous training systems I have been involved with, we have had an extensive administration interface which allowed users and team leads to view progress and find any group wide weaknesses that would benefit from more intensive training to fill a gap.  

I have experimented in the past with adding mentorship capabilities into the system.  Allowing the student to review the “blessed” solution to an exercise along with their attempt and have a mentor review the changes to suggest additional tweaks. Whilst this obviously adds a scalability concern, depending on the use case it does allow for more tailored and engaging training.

## Acknowledgements

This project was conceived and implemented at JPMorgan Chase, built in collaboration with the control-plane.io team. Multiple team members have been involved in creating the project thus far, notably Raoul Millais and Jon Kent.

This project has benefited from many lessons learned in implementing a similar platform focussing on application security - Remediate the Flag. To get a better view of the target functionality, take a look at the RemediateTheFlag project at <https://github.com/sk4ddy/remediatetheflag>

There is still much to do on the project, please get in touch if you are interested in extending it or adding new exercises.

Jon Meadows (Product Owner)
