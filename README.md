## NexusModsUpdateChecker
[![License: MIT](https://img.shields.io/badge/License-MIT-informational.svg)](https://github.com/42milez/ProtocolStack/blob/main/LICENSE)

Watches new releases on [NexusMods.com](https://www.nexusmods.com) and automatically creates a pull request which includes the new file metadata.

### Usage
#### Build
```
make compile RELEASE=true
```

#### Update metadata
Execute `./bin/checker` with `update` sub-command, then mod metadata is stored into `mod.yml`. Also, a pull request which includes the updates is created automatically.
```shell
./bin/checker update
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
./bin/checker -c=false update
```

#### Download mods
```shell
./bin/checker download
```
