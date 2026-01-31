package dialog

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
)

func BuildStandupDialog(triggerID, channelID, userID string) error {
	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		return fmt.Errorf("couldn't fetch standup config: %w", err)
	}
	if standupConfig == nil {
		return fmt.Errorf("standup not configured for this channel")
	}

	existing, err := standup.GetUserStandup(userID, channelID, otime.Now(standupConfig.Timezone))
	if err != nil {
		logger.Error("Couldn't fetch existing standup for dialog pre-population", err, nil)
	}

	var elements []model.DialogElement
	for i, section := range standupConfig.Sections {
		defaultValue := ""
		if existing != nil && existing.Standup[section] != nil {
			defaultValue = strings.Join(*existing.Standup[section], "\n")
		}

		elements = append(elements, model.DialogElement{
			DisplayName: section,
			Name:        fmt.Sprintf("section_%d", i),
			Type:        "textarea",
			Placeholder: fmt.Sprintf("Enter your %s items (one per line)", section),
			Optional:    true,
			Default:     defaultValue,
		})
	}

	siteURL := getSiteURL()

	dialog := model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       siteURL + "/plugins/" + config.PluginName + "/dialog/standup",
		Dialog: model.Dialog{
			CallbackId:       CallbackStandupSubmission,
			Title:            "Submit Standup",
			IntroductionText: "Enter your standup items for each section. One item per line.",
			SubmitLabel:      "Submit",
			NotifyOnCancel:   false,
			State:            EncodeState(channelID),
			Elements:         elements,
		},
	}

	if appErr := config.Mattermost.OpenInteractiveDialog(dialog); appErr != nil {
		return fmt.Errorf("couldn't open standup dialog: %s", appErr.Error())
	}

	return nil
}

func getSiteURL() string {
	mmConfig := config.Mattermost.GetConfig()
	if mmConfig != nil && mmConfig.ServiceSettings.SiteURL != nil {
		return *mmConfig.ServiceSettings.SiteURL
	}
	return ""
}
