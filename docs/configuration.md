<img src="assets/images/banner.png" width="300px">

#

## Plugin Configurations

These settings are configured in the Mattermost **System Console** under the Standup Raven plugin settings.

* **Time Zone** - The default timezone for your Mattermost instance. All datetimes are interpreted in this timezone unless overridden in a channel's standup configuration. Uses IANA timezone format (e.g. `America/Sao_Paulo`).

* **Enable Permission Schema** - Requires Mattermost Enterprise Edition. If enabled, only channel admins, team admins or system admins are allowed to configure standup for a channel or update it.

### Channel-Level Configuration

Each channel has its own standup configuration, managed via `/standup config`. See the [User Guide](user_guide.md) for details on all available channel settings including window times, sections, report format, posting mode, RRULE schedule, and reminders.
