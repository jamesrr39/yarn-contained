# yarn-contained

yarn-contained is a drop-in program for yarn. It runs your requested yarn command inside a docker container, reducing the attack surface to any potential attackers. The container mounts your project as a volume, but does not expose e.g. your home directory where files with secrets may be harvested. Also, the program does not forward any environment variables, e.g. AWS keys, for harvesting by malicious npm packages.

Support
