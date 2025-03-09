# yarn-contained

yarn-contained is a drop-in replacement for Yarn. It runs your requested Yarn command inside a docker container, reducing the attack surface to any potential attackers. The container mounts your project as a volume, but does not expose e.g. your home directory where files with secrets may be harvested. Also, the program does not expose any environment variables, e.g. AWS keys, that you may have set in `.bashrc`, for harvesting by malicious npm packages.

The application works well enough for day-to-day use in August 2023 on Linux. It should also work on MacOS, and may work on Windows. I would be very grateful for bug reports for any issues on any platforms.

Here is a presentation I gave about it: [https://jamesrr39.github.io/talks/yarn-contained-talk/yarn-contained-talk.html](https://jamesrr39.github.io/talks/yarn-contained-talk/yarn-contained-talk.html#1)

## Install

```
go install github.com/jamesrr39/yarn-contained@latest
```

## Features

- Increased security due to isolation.
- `yarn` in container is run by as an ordinary user, not by root, therefore reducing the attack surface.
- Refuses to run any other command than `init` or `create` if the directory doesn't contain `package.json` - a relief to those who have accidentally run `yarn` in their home directory, for example!

## Limitations

- As environment variables, e.g. `NPM_TOKEN` are not carried into the container, dependencies from private NPM repositories cannot be fetched. This also applies to other private repositories or git remotes requiring authorization.

## Contributions

Pull requests welcome for small fixes, for features please open an issue to discuss the feature first.
