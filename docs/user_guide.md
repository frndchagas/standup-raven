<img src="assets/images/banner.png" width="300px">

#

## User Guide

Once the plugin is installed in your Mattermost instance, enabling teams to use it is straightforward.

### 1. Creating channels for standup

Create a new channel, or use an existing one, for each team that wants to use Standup Raven for their standup.

### 2. Configuring channel standup

For each channel, any member can enter configurations for the channel standup. If you are on Mattermost Enterprise Edition and have *Permission Schema* enabled, only a channel admin, team admin or system admin can perform this operation.

Run the following slash command to open the configuration form:

    /standup config

On desktop, this opens a modal. On mobile, it opens an Interactive Dialog.

The following settings are available:

* **Status** - `Enabled` to enable standup for your channel or `Disabled` to disable it.

* **Window Open Time** - The time at which standup reminders will be sent in the channel.

* **Window Close Time** - The time at which an automated standup report will be sent in the channel. The report will include standups for all members who have filled their standups until this time. An additional reminder notification is sent in the channel at 80% completion of the window duration. This message tags those members who have not yet filled their standups.

* **Timezone** - Channel specific timezone to follow for standup notifications (IANA format, e.g. `America/Sao_Paulo`).

* **Sections** - Categories for standup items. For example, if your team fills their standup at the beginning of their work day, suggested sections would be `Yesterday`, `Today` and `Blockers`. At least one section is required.

* **Report Format** - Choose between:
  * `user_aggregated` - Tasks grouped by individual users
  * `type_aggregated` - Tasks grouped by section type

* **Posting Mode** - Choose between:
  * `scheduled` - Standups are batched and posted as a single report at window close time
  * `immediate` - Each standup is posted to the channel as soon as it is submitted

* **Schedule (RRULE)** - Recurring schedule for the standup (weekly or monthly). You can configure frequency, interval, and specific days.

* **Schedule Enabled** - When enabled, the standup schedule is displayed in the channel header.

* **Window Open Reminder** - Enable or disable the window open reminder notification.

* **Window Close Reminder** - Enable or disable the window close reminder notification.

### 3. Adding standup members

The following slash command adds members to the channel's standup:

    /standup addmembers @user1 @user2 @user3

You can specify multiple members separated by a space. Members who are not present in the channel will be automatically added to the channel as well.

### 4. Removing standup members

To remove members from the channel's standup:

    /standup removemembers @user1 @user2

Members are removed from the standup but NOT removed from the channel.

### 5. Filling your standup

Once all the configuration is complete, you can fill your standup in two ways:

* **Desktop**: Click on the Standup Raven button in the channel header to open a modal, or run `/standup`.
* **Mobile**: Run `/standup` to open an Interactive Dialog.

Once saved, you can open the form again to view or update your filled standup.

### 6. Viewing current configuration

To see the saved configuration for the current channel:

    /standup viewconfig

### 7. Generating manual reports

To generate standup reports for specific dates:

    /standup report <public|private> DD-MM-YYYY [date2] [date3]...

* `public` - The generated report is visible to everyone in the channel.
* `private` - The generated report is visible only to you.
* Dates must be in `DD-MM-YYYY` format.

### 8. Getting help

Run the following for a summary of all available commands:

    /standup help
