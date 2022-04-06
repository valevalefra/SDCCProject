# SDCCProject
implementation of two mutual exclusion algorithms: Distributed Lamport and Ricart Agrawala.

### Application

To run the application follow the following steps:
* Open one terminale and digit docker-compose up --build;
* Open three terminal, one for each container and enter inside container, digit  docker exec -it app_peer_* /bin/sh (with 1,2 or 3 in *);
* Follow printed indication.
<br />
If you want increase number of peer change number of replicas in DockerCompose.yml


### Test

For run test set variable *RunTest* to true, inside node.go.

