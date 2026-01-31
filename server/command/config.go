package command

import (
	"github.com/mattermost/mattermost/server/public/model"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/dialog"
	"github.com/standup-raven/standup-raven/server/util"
)

func commandConfig() *Config {
	return &Config{
		AutocompleteData: &model.AutocompleteData{
			Trigger:  "config",
			HelpText: "Open channel standup configuration dialog.",
			RoleID:   model.SystemUserRoleId,
		},
		ExtraHelpText: "",
		Validate:      validateCommandConfig,
		Execute:       executeCommandConfig,
	}
}

func validateCommandConfig(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	return nil, nil
}

func executeCommandConfig(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	if context.IsMobile {
		err := dialog.BuildConfigDialog(context.TriggerId, context.CommandArgs.ChannelId, context.CommandArgs.UserId)
		if err != nil {
			return util.SendEphemeralText("Could not open config dialog: " + err.Error())
		}
		return &model.CommandResponse{}, nil
	}

	config.Mattermost.PublishWebSocketEvent(
		"open_config_modal",
		map[string]interface{}{
			"channel_id": context.CommandArgs.ChannelId,
		},
		&model.WebsocketBroadcast{
			UserId: context.CommandArgs.UserId,
		},
	)

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "Configure your standup in the open modal!",
	}, nil
}
