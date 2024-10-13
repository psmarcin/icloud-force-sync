# iCloud Force Sync
iCloud Force Sync create random files in iCloud documents directory to force synchronization. 

## Instruction
1. `go build *.go`
2. Replace path in `dev.localhost.iCloudForceSync.plist` file
3. `cp ./dev.localhost.iCloudForceSync.plist $HOME/Library/LaunchAgents/dev.localhost.iCloudForceSync.plist` 
4. `launchctl load dev.localhost.iCloudForceSync.plist`


## Edit launchctl
```bash
micro dev.psmarcin.iCloudForceSync.plist
```
## Load launchctl
```bash
launchctl load dev.psmarcin.iCloudForceSync.plist
```

## Unload
```bash
launchctl unload dev.psmarcin.iCloudForceSync.plist
```
