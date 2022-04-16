# SDCCProject
implementation of two mutual exclusion algorithms: Distributed Lamport and Ricart Agrawala.

### Application

To run the application follow the following steps:
* Open one terminal and digit docker-compose up --build;
* Open another terminal for register node with:  docker exec -it  app_register_node_1 /bin/sh, then digit: register;
* Open three terminal, one for each container and enter inside container with: docker exec -it app_peer_* /bin/sh (with 1,2 or 3 in *), then digit: node;
* Follow printed indication.
  <br />
  If you want change number of  peer change number of replicas in DockerCompose.yml and variable MAXPEERS in app/utility/setup


### Test

For run test, set variable *RunTest* to true, inside app/utility/setup.go:
<br /><br />
For testing Lamport:
* Uncomment line 18 in app/test.go;
* Open one terminal and digit docker-compose up --build;
* Open another terminal for register node with:  docker exec -it  app_register_node_1 /bin/sh, then digit: register;
* Open another three terminals, one for each container and enter inside container with:  docker exec -it app_peer_* /bin/sh (with 1,2 or 3 in *), then digit: node;
  <br /><br />

For testing Ricart-Agrawala:
* Uncomment line 19 in app/test.go;
* Follow same indication reported for testing Lamport.

