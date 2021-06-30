## NexusModsWatcher
[![License: MIT](https://img.shields.io/badge/License-MIT-informational.svg)](https://github.com/42milez/ProtocolStack/blob/main/LICENSE)

Watch new releases on [NexusMods.com](https://www.nexusmods.com) and automatically create a pull request which includes the new files.

### Usage
#### Build
```
make compile RELEASE=true
```

#### Update mod metadata
Execute `./bin/watcher` with `update` sub-command, then mod metadata is stored into `mod.yml`. Also, a pull request which includes the updates is created automatically.
```shell
./bin/watcher update
```
If you do not yet add mod id to mod.yml, add it prior to updating.
```yaml
cyberpunk2077:
  - id: MOD_ID         // required
  - name: MOD_NAME     // optional
  - author: MOD_AUTHOR // optional
```
Also, if you would not like to create pull request automatically, you can use `-c` flag.
```shell
./bin/watcher -c=false update
```

#### Download mod
```shell
./bin/watcher download
```

### Development
#### References
- GitHub
  - [GitHub REST API](https://docs.github.com/en/rest)
  - [Permissions for the GITHUB_TOKEN](https://docs.github.com/en/actions/reference/authentication-in-a-workflow#permissions-for-the-github_token)
- Go
  - [Error handling and Go](https://blog.golang.org/error-handling-and-go)
  - [go-github (v35)](https://pkg.go.dev/github.com/google/go-github/v35/github)
  - [go-yaml (v3)](https://pkg.go.dev/gopkg.in/yaml.v3?utm_source=godoc)
- Nexus Mods
  - [API Acceptable Use Policy](https://help.nexusmods.com/article/114-api-acceptable-use-policy)
  - [Nexus Mods Public API](https://app.swaggerhub.com/apis-docs/NexusMods/nexus-mods_public_api_params_in_form_data)
- Stack Overflow
  - [File permission with six bytes in git. What does it mean?](https://unix.stackexchange.com/questions/450480/file-permission-with-six-bytes-in-git-what-does-it-mean)
  - [What could happen if I don't close response.Body?](https://stackoverflow.com/questions/33238518/what-could-happen-if-i-dont-close-response-body)
