# yarn-contained

yarn-contained is a drop-in program for yarn. It runs your requested yarn command inside a docker container, reducing the attack surface to any potential attackers. The container mounts your project as a volume, but does not expose e.g. your home directory where files with secrets may be harvested. Also, the program does not forward any environment variables, e.g. AWS keys, for harvesting by malicious npm packages.

The application is in its early days as of August 2022. It works best on Linux, should work on MacOS, and may work on Windows. I would be very grateful for bug reports for any issues on any platforms.
