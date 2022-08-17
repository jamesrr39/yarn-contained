# yarn-contained

yarn-contained is a drop-in program for yarn. It runs your requested yarn command inside a docker container, reducing the attack surface to any potential attackers. The container mounts your project as a volume, but does not expose e.g. your home directory where files with secrets may be harvested. Also, the program does not expose any environment variables, e.g. AWS keys, that you may have set in `.bashrc`, for harvesting by malicious npm packages.

The application is in its early days as of August 2022. It works best on Linux, should work on MacOS, and may work on Windows. I would be very grateful for bug reports for any issues on any platforms.

Here is a presentation I gave about it: [https://jamesrr39.github.io/talks/yarn-contained-talk/yarn-contained-talk.html](https://jamesrr39.github.io/talks/yarn-contained-talk/yarn-contained-talk.html#1)

## Install

go install github.com/jamesrr39/yarn-contained@latest

## Limitations

- As environment variables, e.g. `NPM_TOKEN` are not carried into the container, dependencies from private NPM repositories cannot be fetched. This also applies to other private repositories or git remotes requiring authorization.
- Currently `yarn init` is not interactive, so to set properties such as name, version, main and license, you must run `yarn init`, and then edit the `package.json` file by hand.

## Contributions

Pull requests welcome for small fixes, for features please open an issue to discuss the feature first.
