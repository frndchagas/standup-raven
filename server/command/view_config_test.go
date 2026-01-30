package command

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/testutil"
)

func TestCommandViewConfig_Definition(t *testing.T) {
	cmd := commandViewConfig()
	assert.Equal(t, "viewconfig", cmd.AutocompleteData.Trigger)
	assert.Equal(t, model.SystemUserRoleId, cmd.AutocompleteData.RoleID)
	assert.NotNil(t, cmd.Validate)
	assert.NotNil(t, cmd.Execute)
}

func TestCommandViewConfig_Validate(t *testing.T) {
	response, err := validateCommandViewConfig(nil, Context{})
	assert.Nil(t, response)
	assert.Nil(t, err)
}

func TestCommandViewConfig_Execute_NoConfig(t *testing.T) {
	defer TearDown()

	testutil.Patch(standup.GetStandupConfig, func(channelID string) (*standup.Config, error) {
		return nil, nil
	})

	ctx := Context{
		CommandArgs: &model.CommandArgs{ChannelId: "test_channel"},
	}

	response, appErr := executeCommandViewConfig(nil, ctx)
	assert.Nil(t, appErr)
	assert.NotNil(t, response)
	assert.Contains(t, response.Text, "Standup is not configured for this channel.")
}

func TestCommandViewConfig_Execute_Error(t *testing.T) {
	defer TearDown()

	testutil.Patch(standup.GetStandupConfig, func(channelID string) (*standup.Config, error) {
		return nil, errors.New("kv store error")
	})

	ctx := Context{
		CommandArgs: &model.CommandArgs{ChannelId: "test_channel"},
	}

	response, appErr := executeCommandViewConfig(nil, ctx)
	assert.Nil(t, appErr)
	assert.NotNil(t, response)
	assert.Contains(t, response.Text, "Error fetching standup configuration")
}

func TestCommandViewConfig_Execute_Success(t *testing.T) {
	defer TearDown()

	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	mockAPI.On("GetUser", "user_1").Return(&model.User{Username: "alice"}, nil)
	mockAPI.On("GetUser", "user_2").Return(&model.User{Username: "bob"}, nil)

	location, _ := time.LoadLocation("America/Fortaleza")
	openTime := otime.OTime{Time: time.Date(2026, 1, 30, 9, 0, 0, 0, location)}
	closeTime := otime.OTime{Time: time.Date(2026, 1, 30, 12, 0, 0, 0, location)}

	testutil.Patch(standup.GetStandupConfig, func(channelID string) (*standup.Config, error) {
		return &standup.Config{
			ChannelID:                  channelID,
			Enabled:                    true,
			Timezone:                   "America/Fortaleza",
			WindowOpenTime:             openTime,
			WindowCloseTime:            closeTime,
			ReportFormat:               "user_aggregated",
			ScheduleEnabled:            true,
			WindowOpenReminderEnabled:  true,
			WindowCloseReminderEnabled: false,
			Sections:                   []string{"Done", "In Progress", "Blockers"},
			Members:                    []string{"user_1", "user_2"},
		}, nil
	})

	ctx := Context{
		CommandArgs: &model.CommandArgs{ChannelId: "test_channel"},
	}

	response, appErr := executeCommandViewConfig(nil, ctx)
	assert.Nil(t, appErr)
	assert.NotNil(t, response)
	assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)

	text := response.Text
	assert.Contains(t, text, "### Standup Configuration")
	assert.Contains(t, text, "| **Enabled** | Yes |")
	assert.Contains(t, text, "| **Timezone** | America/Fortaleza |")
	assert.Contains(t, text, "| **Window Open Time** | 09:00 |")
	assert.Contains(t, text, "| **Window Close Time** | 12:00 |")
	assert.Contains(t, text, "| **Report Format** | user_aggregated |")
	assert.Contains(t, text, "| **Schedule Enabled** | Yes |")
	assert.Contains(t, text, "| **Window Open Reminder** | Yes |")
	assert.Contains(t, text, "| **Window Close Reminder** | No |")
	assert.Contains(t, text, "Done, In Progress, Blockers")
	assert.Contains(t, text, "@alice, @bob")
}

func TestCommandViewConfig_Execute_ZeroWindowTimes(t *testing.T) {
	defer TearDown()

	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	testutil.Patch(standup.GetStandupConfig, func(channelID string) (*standup.Config, error) {
		return &standup.Config{
			ChannelID:    channelID,
			Enabled:      false,
			Timezone:     "UTC",
			ReportFormat: "user_aggregated",
			Sections:     []string{"Tasks"},
			Members:      []string{},
		}, nil
	})

	ctx := Context{
		CommandArgs: &model.CommandArgs{ChannelId: "test_channel"},
	}

	response, appErr := executeCommandViewConfig(nil, ctx)
	assert.Nil(t, appErr)
	assert.NotNil(t, response)

	text := response.Text
	assert.Contains(t, text, "| **Enabled** | No |")
	assert.Contains(t, text, "| **Window Open Time** | Not set |")
	assert.Contains(t, text, "| **Window Close Time** | Not set |")
}

func TestCommandViewConfig_Execute_UserNotFound(t *testing.T) {
	defer TearDown()

	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	mockAPI.On("GetUser", "user_exists").Return(&model.User{Username: "alice"}, nil)
	mockAPI.On("GetUser", "user_deleted").Return(nil, model.NewAppError("GetUser", "not_found", nil, "", 404))

	location, _ := time.LoadLocation("UTC")
	openTime := otime.OTime{Time: time.Date(2026, 1, 30, 9, 0, 0, 0, location)}
	closeTime := otime.OTime{Time: time.Date(2026, 1, 30, 12, 0, 0, 0, location)}

	testutil.Patch(standup.GetStandupConfig, func(channelID string) (*standup.Config, error) {
		return &standup.Config{
			ChannelID:       channelID,
			Enabled:         true,
			Timezone:        "UTC",
			WindowOpenTime:  openTime,
			WindowCloseTime: closeTime,
			ReportFormat:    "user_aggregated",
			Sections:        []string{"Tasks"},
			Members:         []string{"user_exists", "user_deleted"},
		}, nil
	})

	ctx := Context{
		CommandArgs: &model.CommandArgs{ChannelId: "test_channel"},
	}

	response, appErr := executeCommandViewConfig(nil, ctx)
	assert.Nil(t, appErr)
	assert.NotNil(t, response)

	text := response.Text
	// user_exists should be resolved to @alice
	assert.Contains(t, text, "@alice")
	// user_deleted should fallback to raw ID
	assert.Contains(t, text, "user_deleted")
	// Both should be present in the members field
	memberLine := ""
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "Members") {
			memberLine = line
			break
		}
	}
	assert.Contains(t, memberLine, "@alice")
	assert.Contains(t, memberLine, "user_deleted")
}

func TestCommandViewConfig_InCommandsMap(t *testing.T) {
	_, ok := commands["viewconfig"]
	assert.True(t, ok, "viewconfig should be registered in the commands map")
}
