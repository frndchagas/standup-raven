package dialog

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
)

func BuildConfigDialog(triggerID, channelID, userID string) error {
	existing, err := standup.GetStandupConfig(channelID)
	if err != nil {
		logger.Error("Couldn't fetch existing config for dialog pre-population", err, nil)
	}

	enabledDefault := "true"
	sectionsDefault := ""
	openTimeDefault := "09:00"
	closeTimeDefault := "17:00"
	timezoneDefault := config.GetConfig().TimeZone
	postingModeDefault := config.PostingModeScheduled
	reportFormatDefault := config.ReportFormatUserAggregated

	if existing != nil {
		if existing.Enabled {
			enabledDefault = "true"
		} else {
			enabledDefault = "false"
		}
		sectionsDefault = strings.Join(existing.Sections, "\n")
		if ts := existing.WindowOpenTime.GetTimeString(); ts != "" {
			openTimeDefault = ts
		}
		if ts := existing.WindowCloseTime.GetTimeString(); ts != "" {
			closeTimeDefault = ts
		}
		if existing.Timezone != "" {
			timezoneDefault = existing.Timezone
		}
		if existing.PostingMode != "" {
			postingModeDefault = existing.PostingMode
		}
		if existing.ReportFormat != "" {
			reportFormatDefault = existing.ReportFormat
		}
	}

	elements := []model.DialogElement{
		{
			DisplayName: "Enabled",
			Name:        "enabled",
			Type:        "bool",
			Default:     enabledDefault,
			HelpText:    "Enable or disable standup for this channel.",
		},
		{
			DisplayName: "Sections",
			Name:        "sections",
			Type:        "textarea",
			Default:     sectionsDefault,
			Placeholder: "One section per line (e.g. Yesterday, Today, Blockers)",
			HelpText:    "Standup sections, one per line.",
		},
		{
			DisplayName: "Window Open Time",
			Name:        "window_open_time",
			Type:        "text",
			SubType:     "text",
			Default:     openTimeDefault,
			Placeholder: "HH:MM (e.g. 09:00)",
			HelpText:    "Time when the standup window opens (24h format).",
		},
		{
			DisplayName: "Window Close Time",
			Name:        "window_close_time",
			Type:        "text",
			SubType:     "text",
			Default:     closeTimeDefault,
			Placeholder: "HH:MM (e.g. 17:00)",
			HelpText:    "Time when the standup window closes (24h format).",
		},
		{
			DisplayName: "Timezone",
			Name:        "timezone",
			Type:        "text",
			SubType:     "text",
			Default:     timezoneDefault,
			Placeholder: "IANA timezone (e.g. America/Sao_Paulo)",
			HelpText:    "IANA timezone identifier.",
		},
		{
			DisplayName: "Posting Mode",
			Name:        "posting_mode",
			Type:        "select",
			Default:     postingModeDefault,
			HelpText:    "When to post standups to the channel.",
			Options: []*model.PostActionOptions{
				{Text: "Scheduled", Value: config.PostingModeScheduled},
				{Text: "Immediate", Value: config.PostingModeImmediate},
			},
		},
		{
			DisplayName: "Report Format",
			Name:        "report_format",
			Type:        "select",
			Default:     reportFormatDefault,
			HelpText:    "How the standup report is organized.",
			Options: []*model.PostActionOptions{
				{Text: "User Aggregated", Value: config.ReportFormatUserAggregated},
				{Text: "Type Aggregated", Value: config.ReportFormatTypeAggregated},
			},
		},
	}

	siteURL := getSiteURL()

	dialog := model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       siteURL + "/plugins/" + config.PluginName + "/dialog/config",
		Dialog: model.Dialog{
			CallbackId:       CallbackConfigSubmission,
			Title:            "Standup Configuration",
			IntroductionText: "Configure standup settings for this channel. Use `/standup addmembers` and `/standup removemembers` to manage members.",
			SubmitLabel:      "Save",
			NotifyOnCancel:   false,
			State:            EncodeState(channelID),
			Elements:         elements,
		},
	}

	if appErr := config.Mattermost.OpenInteractiveDialog(dialog); appErr != nil {
		return fmt.Errorf("couldn't open config dialog: %s", appErr.Error())
	}

	return nil
}
