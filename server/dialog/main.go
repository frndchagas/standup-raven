package dialog

import "encoding/json"

const (
	CallbackStandupSubmission = "standup_dialog_submission"
	CallbackConfigSubmission  = "config_dialog_submission"
)

type DialogState struct {
	ChannelID string `json:"channel_id"`
}

func EncodeState(channelID string) string {
	data, _ := json.Marshal(DialogState{ChannelID: channelID})
	return string(data)
}

func DecodeState(state string) (DialogState, error) {
	var ds DialogState
	err := json.Unmarshal([]byte(state), &ds)
	return ds, err
}
