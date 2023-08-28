# yarn-contained

yarn-contained is a drop-in program for yarn. It runs your requested yarn command inside a docker container, reducing the attack surface to any potential attackers. The container mounts your project as a volume, but does not expose e.g. your home directory where files with secrets may be harvested. Also, the program does not expose any environment variables, e.g. AWS keys, that you may have set in `.bashrc`, for harvesting by malicious npm packages. Additionally, as purely a side feature, the application exits early if there is no `package.json` in the working direction (with exceptions for running the `yarn init` and `yarn create` commands).

The application works well enough for day-to-day use in August 2023 on Linux. It should also work on MacOS, and may work on Windows. I would be very grateful for bug reports for any issues on any platforms.

Here is a presentation I gave about it: [https://jamesrr39.github.io/talks/yarn-contained-talk/yarn-contained-talk.html](https://jamesrr39.github.io/talks/yarn-contained-talk/yarn-contained-talk.html#1)

## Install

go install github.com/jamesrr39/yarn-contained@latest

## Limitations

- As environment variables, e.g. `NPM_TOKEN` are not carried into the container, dependencies from private NPM repositories cannot be fetched. This also applies to other private repositories or git remotes requiring authorization.

## Contributions

Pull requests welcome for small fixes, for features please open an issue to discuss the feature first.
