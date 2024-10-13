# iCloud Force Sync
iCloud Force Sync create random files in iCloud documents directory to force synchronization. 

## Instruction
1. Download latest version from [link](https://github.com/psmarcin/icloud-force-sync/releases)
2. Run it `./icloduf-force-sync`

It will automatically add task in background to run every time your use your computer. 


## Technical
### Load launchctl
```bash
launchctl load dev.localhost.iCloudForceSync.plist
```

## Unload
```bash
launchctl unload dev.localhost.iCloudForceSync.plist
```
