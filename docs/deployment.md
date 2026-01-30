<img src="assets/images/banner.png" width="300px">

#

## Deployment

### Local deploy via mmctl

With Mattermost running locally (port 8065):

```bash
make dist      # build for all platforms
make deploy    # install and enable the plugin via mmctl --local
```

### Manual deploy

1. Build the packages:

```bash
make dist
```

2. In Mattermost, go to **System Console > Plugins > Management**
3. Upload the `.tar.gz` matching the server platform
4. Enable the plugin

### CI/CD deploy

The `release.yml` GitHub Actions workflow automatically creates a GitHub Release with artifacts when a `v*` tag is pushed:

```bash
git tag v3.4.0
git push origin v3.4.0
```

## Usage

1. Create a channel for your team standup or use an existing one.

2. Configure the standup:

        /standup config

    Opens a modal with the channel configuration.

3. Add members:

        /standup addmembers <usernames...>

    Usernames can be specified as @mentions.

4. Verify the saved config:

        /standup viewconfig

5. Fill your standup by clicking the Standup Raven icon in the channel header bar.

6. Access help anytime:

        /standup help
