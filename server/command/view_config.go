package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

func commandViewConfig() *Config {
	return &Config{
		AutocompleteData: &model.AutocompleteData{
			Trigger:  "viewconfig",
			HelpText: "Display saved standup configuration for this channel.",
			RoleID:   model.SystemUserRoleId,
		},
		Validate: validateCommandViewConfig,
		Execute:  executeCommandViewConfig,
	}
}

func validateCommandViewConfig(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	return nil, nil
}

func executeCommandViewConfig(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	channelID := context.CommandArgs.ChannelId

	conf, err := standup.GetStandupConfig(channelID)
	if err != nil {
		return util.SendEphemeralText("Error fetching standup configuration for this channel.")
	}

	if conf == nil {
		return util.SendEphemeralText("Standup is not configured for this channel.")
	}

	enabledText := "No"
	if conf.Enabled {
		enabledText = "Yes"
	}

	scheduleEnabledText := "No"
	if conf.ScheduleEnabled {
		scheduleEnabledText = "Yes"
	}

	windowOpenReminderText := "No"
	if conf.WindowOpenReminderEnabled {
		windowOpenReminderText = "Yes"
	}

	windowCloseReminderText := "No"
	if conf.WindowCloseReminderEnabled {
		windowCloseReminderText = "Yes"
	}

	members := make([]string, len(conf.Members))
	for i, memberID := range conf.Members {
		user, appErr := config.Mattermost.GetUser(memberID)
		if appErr != nil {
			members[i] = memberID
		} else {
			members[i] = "@" + user.Username
		}
	}

	emptyTime := otime.OTime{}
	windowOpenTimeText := "Not set"
	if conf.WindowOpenTime != emptyTime {
		windowOpenTimeText = conf.WindowOpenTime.Format("15:04")
	}
	windowCloseTimeText := "Not set"
	if conf.WindowCloseTime != emptyTime {
		windowCloseTimeText = conf.WindowCloseTime.Format("15:04")
	}

	text := fmt.Sprintf(
		"### Standup Configuration\n\n"+
			"| Setting | Value |\n"+
			"|:--------|:------|\n"+
			"| **Enabled** | %s |\n"+
			"| **Timezone** | %s |\n"+
			"| **Window Open Time** | %s |\n"+
			"| **Window Close Time** | %s |\n"+
			"| **Report Format** | %s |\n"+
			"| **Schedule Enabled** | %s |\n"+
			"| **Window Open Reminder** | %s |\n"+
			"| **Window Close Reminder** | %s |\n"+
			"| **Sections** | %s |\n"+
			"| **Members** | %s |\n",
		enabledText,
		conf.Timezone,
		windowOpenTimeText,
		windowCloseTimeText,
		conf.ReportFormat,
		scheduleEnabledText,
		windowOpenReminderText,
		windowCloseReminderText,
		strings.Join(conf.Sections, ", "),
		strings.Join(members, ", "),
	)

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text,
	}, nil
}
