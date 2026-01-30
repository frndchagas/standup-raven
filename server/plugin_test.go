package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/standup-raven/standup-raven/server/testutil"
)

func TearDown() {
	testutil.UnpatchAll()
}

func setupMockAPI() *plugintest.API {
	api := &plugintest.API{}
	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	api.On("LogInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	api.On("KVSetWithOptions", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	api.On("GetServerVersion").Return("9.0.0")
	return api
}

func setupTestBundlePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	imgDir := filepath.Join(dir, "webapp", "static")
	err := os.MkdirAll(imgDir, 0750)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(imgDir, "logo.png"), []byte("fake-png"), 0600)
	assert.NoError(t, err)
	return dir
}

func TestSetUpBot(t *testing.T) {
	defer TearDown()
	bundlePath := setupTestBundlePath(t)
	api := setupMockAPI()
	api.On("EnsureBotUser", mock.Anything).Return("botUserID", nil)
	api.On("GetBundlePath").Return(bundlePath, nil)
	api.On("SetProfileImage", mock.Anything, mock.Anything).Return(nil)

	p := &Plugin{}
	p.SetAPI(api)
	botID, err := p.setUpBot()
	assert.Nil(t, err)
	assert.Equal(t, "botUserID", botID)
}

func TestSetUpBot_CreateBot(t *testing.T) {
	defer TearDown()
	bundlePath := setupTestBundlePath(t)
	api := setupMockAPI()
	api.On("EnsureBotUser", mock.Anything).Return("newBotID", nil)
	api.On("GetBundlePath").Return(bundlePath, nil)
	api.On("SetProfileImage", mock.Anything, mock.Anything).Return(nil)

	p := &Plugin{}
	p.SetAPI(api)
	botID, err := p.setUpBot()
	assert.Nil(t, err)
	assert.Equal(t, "newBotID", botID)
}

func TestSetUpBot_EnsureBot_Error(t *testing.T) {
	defer TearDown()
	api := setupMockAPI()
	api.On("EnsureBotUser", mock.Anything).Return("", errors.New("create failed"))

	p := &Plugin{}
	p.SetAPI(api)
	_, err := p.setUpBot()
	assert.NotNil(t, err)
}
