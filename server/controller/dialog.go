package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/dialog"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/standup/notification"
	"github.com/standup-raven/standup-raven/server/util"
)

var submitStandupDialog = &Endpoint{
	Path:    "/dialog/standup",
	Method:  http.MethodPost,
	Execute: executeSubmitStandupDialog,
}

var submitConfigDialog = &Endpoint{
	Path:    "/dialog/config",
	Method:  http.MethodPost,
	Execute: executeSubmitConfigDialog,
}

func executeSubmitStandupDialog(w http.ResponseWriter, r *http.Request) error {
	var request model.SubmitDialogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Couldn't decode dialog submission request", err, nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	if request.Cancelled { //nolint:misspell // field name from Mattermost SDK
		w.WriteHeader(http.StatusOK)
		return nil
	}

	state, err := dialog.DecodeState(request.State)
	if err != nil {
		logger.Error("Couldn't decode dialog state", err, nil)
		http.Error(w, "Invalid dialog state", http.StatusBadRequest)
		return err
	}

	channelID := state.ChannelID
	userID := request.UserId

	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		respondWithDialogError(w, "Error fetching standup configuration.")
		return err
	}
	if standupConfig == nil {
		respondWithDialogError(w, "Standup is not configured for this channel.")
		return fmt.Errorf("standup not configured for channel: %s", channelID)
	}

	standupData := map[string]*[]string{}
	for i, section := range standupConfig.Sections {
		key := fmt.Sprintf("section_%d", i)
		raw, ok := request.Submission[key]
		if !ok || raw == nil {
			continue
		}
		rawStr, ok := raw.(string)
		if !ok || strings.TrimSpace(rawStr) == "" {
			continue
		}

		lines := strings.Split(rawStr, "\n")
		var tasks []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				tasks = append(tasks, trimmed)
			}
		}
		if len(tasks) > 0 {
			standupData[section] = &tasks
		}
	}

	userStandup := &standup.UserStandup{
		UserID:    userID,
		ChannelID: channelID,
		Standup:   standupData,
	}

	if err := userStandup.IsValid(); err != nil {
		respondWithDialogError(w, err.Error())
		return nil
	}

	if err := standup.SaveUserStandup(userStandup); err != nil {
		respondWithDialogError(w, "Failed to save standup.")
		return err
	}

	if standupConfig.PostingMode == config.PostingModeImmediate {
		if postErr := notification.PostIndividualStandup(userStandup); postErr != nil {
			logger.Error("Failed to post individual standup from dialog, but standup was saved", postErr, nil)
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func executeSubmitConfigDialog(w http.ResponseWriter, r *http.Request) error {
	var request model.SubmitDialogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Couldn't decode config dialog submission request", err, nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	if request.Cancelled { //nolint:misspell // field name from Mattermost SDK
		w.WriteHeader(http.StatusOK)
		return nil
	}

	state, err := dialog.DecodeState(request.State)
	if err != nil {
		logger.Error("Couldn't decode dialog state", err, nil)
		http.Error(w, "Invalid dialog state", http.StatusBadRequest)
		return err
	}

	channelID := state.ChannelID
	userID := request.UserId

	// Check guest permission
	userRoles, appErr := util.GetUserRoles(userID, channelID)
	if appErr != nil {
		respondWithDialogError(w, "Could not verify user permissions.")
		return appErr
	}
	userRolesMap := make(map[string]bool, len(userRoles))
	for _, role := range userRoles {
		userRolesMap[role] = true
	}
	if userRolesMap[model.SystemGuestRoleId] {
		respondWithDialogError(w, "Guest users are not allowed to configure standup.")
		return nil
	}

	// Check permission schema
	if config.GetConfig().PermissionSchemaEnabled {
		isAdmin := userRolesMap[model.SystemAdminRoleId] || userRolesMap[model.TeamAdminRoleId] || userRolesMap[model.ChannelAdminRoleId]
		if !isAdmin {
			respondWithDialogError(w, "You do not have permission to configure standup for this channel.")
			return nil
		}
	}

	// Parse submission fields
	enabled := true
	if raw, ok := request.Submission["enabled"]; ok && raw != nil {
		switch v := raw.(type) {
		case string:
			enabled = v == "true"
		case bool:
			enabled = v
		}
	}

	sectionsRaw := ""
	if raw, ok := request.Submission["sections"]; ok && raw != nil {
		if s, ok := raw.(string); ok {
			sectionsRaw = s
		}
	}
	var sections []string
	for _, line := range strings.Split(sectionsRaw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			sections = append(sections, trimmed)
		}
	}
	if len(sections) == 0 {
		respondWithDialogError(w, "At least one section is required.")
		return nil
	}

	openTimeStr := getSubmissionString(request.Submission, "window_open_time")
	closeTimeStr := getSubmissionString(request.Submission, "window_close_time")
	timezone := getSubmissionString(request.Submission, "timezone")
	postingMode := getSubmissionString(request.Submission, "posting_mode")
	reportFormat := getSubmissionString(request.Submission, "report_format")

	openTime, err := otime.Parse(openTimeStr)
	if err != nil {
		respondWithDialogFieldError(w, "window_open_time", "Invalid time format. Use HH:MM (e.g. 09:00).")
		return nil
	}

	closeTime, err := otime.Parse(closeTimeStr)
	if err != nil {
		respondWithDialogFieldError(w, "window_close_time", "Invalid time format. Use HH:MM (e.g. 17:00).")
		return nil
	}

	if _, err := time.LoadLocation(timezone); err != nil {
		respondWithDialogFieldError(w, "timezone", fmt.Sprintf("Invalid timezone: %q. Use IANA format (e.g. America/Sao_Paulo).", timezone))
		return nil
	}

	// Load existing config to preserve non-editable fields
	existing, _ := standup.GetStandupConfig(channelID)

	conf := &standup.Config{
		ChannelID:       channelID,
		Enabled:         enabled,
		Sections:        sections,
		WindowOpenTime:  openTime,
		WindowCloseTime: closeTime,
		Timezone:        timezone,
		PostingMode:     postingMode,
		ReportFormat:    reportFormat,
	}

	if existing != nil {
		conf.Members = existing.Members
		conf.RRuleString = existing.RRuleString
		conf.StartDate = existing.StartDate
		conf.WindowOpenReminderEnabled = existing.WindowOpenReminderEnabled
		conf.WindowCloseReminderEnabled = existing.WindowCloseReminderEnabled
		conf.ScheduleEnabled = existing.ScheduleEnabled
	} else {
		// Sensible defaults for new config
		conf.Members = []string{}
		conf.RRuleString = "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
		conf.StartDate = time.Now()
		conf.WindowOpenReminderEnabled = true
		conf.WindowCloseReminderEnabled = true
		conf.ScheduleEnabled = true
	}

	if err := conf.PreSave(); err != nil {
		respondWithDialogError(w, "Error processing configuration: "+err.Error())
		return nil
	}

	if err := conf.IsValid(); err != nil {
		respondWithDialogError(w, err.Error())
		return nil
	}

	if _, err := standup.SaveStandupConfig(conf); err != nil {
		respondWithDialogError(w, "Failed to save configuration.")
		return err
	}

	if err := standup.AddStandupChannel(channelID); err != nil {
		respondWithDialogError(w, "Failed to register standup channel.")
		return err
	}

	// Notify desktop clients about the config change
	event := "add_active_channel"
	if !conf.Enabled {
		event = "remove_active_channel"
	}
	config.Mattermost.PublishWebSocketEvent(
		event,
		map[string]interface{}{
			"channel_id": channelID,
		},
		&model.WebsocketBroadcast{
			UserId: userID,
		},
	)

	w.WriteHeader(http.StatusOK)
	return nil
}

func getSubmissionString(submission map[string]interface{}, key string) string {
	raw, ok := submission[key]
	if !ok || raw == nil {
		return ""
	}
	s, ok := raw.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

func respondWithDialogError(w http.ResponseWriter, message string) {
	resp := model.SubmitDialogResponse{
		Error: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func respondWithDialogFieldError(w http.ResponseWriter, field, message string) {
	resp := model.SubmitDialogResponse{
		Errors: map[string]string{
			field: message,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
