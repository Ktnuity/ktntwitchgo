package ktntwitchgo

import (
	"testing"
)

func TestParseOptions(t *testing.T) {
	test := formTest(t, "parse options")

	input := &GetVideosOptions{
		UserID: asRef("12345"),
		GameID: asRef("twitch"),
	}

	test.expect("user_id=12345&game_id=twitch", parseOptions(input))

	input = &GetVideosOptions{
		ID: []string{"123", "456"},
	}

	test.expect("id=123&id=456", parseOptions(input))

	input = &GetVideosOptions{
		ID: []string{"123", "456"},
		GameID: asRef("twitch"),
	}

	test.expect("id=123&id=456&game_id=twitch", parseOptions(input))

	input = &GetVideosOptions{
		UserID: asRef("12345"),
	}
	input.First = asRef(1)

	test.expect("first=1&user_id=12345", parseOptions(input))
}

func TestMixedParam(t *testing.T) {
	test := formTest(t, "parse mixed param")

	test.expect("game_name=twitch", parseMixedParam("twitch", "game_name", "game_id"))

	test.expect("game_id=12345", parseMixedParam(12345, "game_name", "game_id"))

	test.expect("game_name=twitch&game_name=glitch", parseMixedParam([]string{"twitch","glitch"}, "game_name", "game_id"))

	test.expect("game_id=123&game_id=456", parseMixedParam([]int{123, 456}, "game_name", "game_id"))

	test.expect("game_id=123&game_id=456", parseMixedParam([]string{"123", "456"}, "game_name", "game_id"))

	test.expect("game_id=123&game_name=twitch", parseMixedParam([]string{"123", "twitch"}, "game_name", "game_id"))
}

func TestChooseKey(t *testing.T) {
	test := formTest(t, "choose key")

	test.expect("foo", chooseKey(true, "foo", "NEVER"))
	test.expect("bar", chooseKey(false, "NEVER", "bar"))
}

func TestIsNumber(t *testing.T) {
	test := formTest(t, "test number")

	test.expect(true, isNumber("123"))
	test.expect(true, isNumber("0000"))
	test.expect(true, isNumber("1425438400038830143"))
	test.expect(false, isNumber("twitch"))
	test.expect(false, isNumber("ABCDEF"))
	test.expect(false, isNumber("DeadBeef"))
	test.expect(false, isNumber(""))
	test.expect(false, isNumber("O"))
	test.expect(true, isNumber("0"))
	test.expect(false, isNumber(" "))
}

func TestCreateTwitchApi(t *testing.T) {
	config := TwitchApiConfig{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}

	client, wg := CreateTwitchApi(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	test := formTest(t, "create twitch api client")
	test.expect("test_client_id", client.clientID)
	test.expect("test_client_secret", client.clientSecret)
	test.expect("https://api.twitch.tv/helix", client.baseURL)
	test.expect("https://ingest.twitch.tv", client.ingestBaseURL)
	test.expect(false, client.Verbose)
	test.expect(false, client.throwRateLimitErrors)

	if wg != nil {
		t.Errorf("Expected nil wait group when access token is not provided")
	}

	// Test with access token
	accessToken := "test_access_token"
	refreshToken := "test_refresh_token"
	config2 := TwitchApiConfig{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		Scopes:       []Scope{ScopeUserReadChat, ScopeUserWriteChat},
	}

	client2, wg2 := CreateTwitchApi(config2)

	if client2 == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if wg2 == nil {
		t.Errorf("Expected wait group when access token is provided")
	}

	if len(client2.scopes) != 2 {
		t.Errorf("Expected 2 scopes, got %d", len(client2.scopes))
	}
}

func TestCreateTwitchApiWithRateLimitErrors(t *testing.T) {
	throwErrors := true
	config := TwitchApiConfig{
		ClientID:              "test_client_id",
		ClientSecret:          "test_client_secret",
		ThrowRatelimitErrors:  &throwErrors,
	}

	client, _ := CreateTwitchApi(config)

	test := formTest(t, "create client with rate limit errors")
	test.expect(true, client.throwRateLimitErrors)
}

func TestClientError(t *testing.T) {
	client := &Client{}
	err := client.error("test error message")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	test := formTest(t, "create client error")
	test.expect("test error message", err.Error())
}

func TestClientAddEventHandler(t *testing.T) {
	client := &Client{
		eventHandlers: make(map[string][]EventHandler),
	}

	callCount := 0
	handler := func(data any) {
		callCount++
	}

	client.AddEventHandler("test_event", handler)

	if len(client.eventHandlers["test_event"]) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(client.eventHandlers["test_event"]))
	}

	// Add another handler
	client.AddEventHandler("test_event", handler)

	if len(client.eventHandlers["test_event"]) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(client.eventHandlers["test_event"]))
	}

	// Trigger event
	client.emit("test_event", nil)

	if callCount != 2 {
		t.Errorf("Expected handlers to be called 2 times, got %d", callCount)
	}
}

func TestClientRemoveEventHandler(t *testing.T) {
	client := &Client{
		eventHandlers: make(map[string][]EventHandler),
	}

	handler := func(data any) {}

	client.AddEventHandler("test_event", handler)
	client.AddEventHandler("other_event", handler)

	if len(client.eventHandlers) != 2 {
		t.Errorf("Expected 2 event types, got %d", len(client.eventHandlers))
	}

	client.RemoveEventHandler("test_event")

	if len(client.eventHandlers) != 1 {
		t.Errorf("Expected 1 event type after removal, got %d", len(client.eventHandlers))
	}

	if _, exists := client.eventHandlers["test_event"]; exists {
		t.Error("Expected test_event to be removed")
	}

	if _, exists := client.eventHandlers["other_event"]; !exists {
		t.Error("Expected other_event to still exist")
	}
}

func TestClientEmit(t *testing.T) {
	client := &Client{
		eventHandlers: make(map[string][]EventHandler),
	}

	receivedData := ""
	handler := func(data any) {
		if str, ok := data.(string); ok {
			receivedData = str
		}
	}

	client.AddEventHandler("test_event", handler)
	client.emit("test_event", "test_data")

	test := formTest(t, "emit event")
	test.expect("test_data", receivedData)

	// Test emitting event with no handlers
	client.emit("non_existent_event", "should not panic")
}

func TestClientHasScope(t *testing.T) {
	client := &Client{
		scopes: []Scope{ScopeUserReadChat, ScopeUserWriteChat, ScopeBitsRead},
	}

	test := formTest(t, "check if client has scope")
	test.expect(true, client.hasScope(ScopeUserReadChat))
	test.expect(true, client.hasScope(ScopeUserWriteChat))
	test.expect(true, client.hasScope(ScopeBitsRead))
	test.expect(false, client.hasScope(ScopeUserBot))
	test.expect(false, client.hasScope(ScopeChannelBot))
}

func TestClientGenerateAuthURL(t *testing.T) {
	redirectURI := "http://localhost:3000/callback"
	client := &Client{
		clientID:    "test_client_id",
		scopes:      []Scope{ScopeUserReadChat, ScopeUserWriteChat},
		redirectURI: &redirectURI,
	}

	url := client.GenerateAuthURL()

	if url == "" {
		t.Fatal("Expected non-empty URL")
	}

	// Check that URL contains required components
	if !containsHelper(url, "https://id.twitch.tv/oauth2/authorize") {
		t.Error("Expected URL to contain base authorization URL")
	}

	if !containsHelper(url, "client_id=test_client_id") {
		t.Error("Expected URL to contain client_id")
	}

	if !containsHelper(url, "response_type=code") {
		t.Error("Expected URL to contain response_type=code")
	}

	if !containsHelper(url, "redirect_uri=http") {
		t.Error("Expected URL to contain redirect_uri")
	}

	if !containsHelper(url, "scope=") {
		t.Error("Expected URL to contain scope parameter")
	}
}

func TestClientGenerateAuthURLNoScopes(t *testing.T) {
	client := &Client{
		clientID: "test_client_id",
	}

	url := client.GenerateAuthURL()

	if url == "" {
		t.Fatal("Expected non-empty URL")
	}

	if !containsHelper(url, "client_id=test_client_id") {
		t.Error("Expected URL to contain client_id")
	}
}

func TestExtractRateLimit(t *testing.T) {
	client := &Client{}

	// Create mock headers
	headers := make(map[string][]string)
	headers["Ratelimit-Limit"] = []string{"800"}
	headers["Ratelimit-Remaining"] = []string{"799"}
	headers["Ratelimit-Reset"] = []string{"1234567890"}

	rateLimit := client.extractRateLimit(headers)

	test := formTest(t, "extract rate limit from headers")
	test.expect(800, rateLimit.Limit)
	test.expect(799, rateLimit.Remaining)
	test.expect(1234567890, rateLimit.Reset)
}

func TestExtractRateLimitInvalidHeaders(t *testing.T) {
	client := &Client{}

	// Create headers with invalid values
	headers := make(map[string][]string)
	headers["Ratelimit-Limit"] = []string{"invalid"}
	headers["Ratelimit-Remaining"] = []string{"also_invalid"}
	headers["Ratelimit-Reset"] = []string{"not_a_number"}

	rateLimit := client.extractRateLimit(headers)

	test := formTest(t, "extract rate limit with invalid headers")
	test.expect(0, rateLimit.Limit)
	test.expect(0, rateLimit.Remaining)
	test.expect(0, rateLimit.Reset)
}

func TestCommercialLengthIsValid(t *testing.T) {
	test := formTest(t, "validate commercial length")

	test.expect(true, CommercialLength30.IsValid())
	test.expect(true, CommercialLength60.IsValid())
	test.expect(true, CommercialLength90.IsValid())
	test.expect(true, CommercialLength120.IsValid())
	test.expect(true, CommercialLength150.IsValid())
	test.expect(true, CommercialLength180.IsValid())
	test.expect(false, CommercialLength(45).IsValid())
	test.expect(false, CommercialLength(0).IsValid())
}

func TestCommercialLengthString(t *testing.T) {
	test := formTest(t, "convert commercial length to string")

	test.expect("30", CommercialLength30.String())
	test.expect("60", CommercialLength60.String())
	test.expect("90", CommercialLength90.String())
	test.expect("120", CommercialLength120.String())
	test.expect("150", CommercialLength150.String())
	test.expect("180", CommercialLength180.String())
	test.expect("invalid", CommercialLength(45).String())
}

func TestCommercialLengthValidate(t *testing.T) {
	test := formTest(t, "validate and round commercial length")

	test.expect(CommercialLength30, CommercialLength(15).Validate())
	test.expect(CommercialLength30, CommercialLength(30).Validate())
	test.expect(CommercialLength60, CommercialLength(45).Validate())
	test.expect(CommercialLength60, CommercialLength(60).Validate())
	test.expect(CommercialLength90, CommercialLength(75).Validate())
	test.expect(CommercialLength90, CommercialLength(90).Validate())
	test.expect(CommercialLength120, CommercialLength(105).Validate())
	test.expect(CommercialLength120, CommercialLength(120).Validate())
	test.expect(CommercialLength150, CommercialLength(135).Validate())
	test.expect(CommercialLength150, CommercialLength(150).Validate())
	test.expect(CommercialLength180, CommercialLength(165).Validate())
	test.expect(CommercialLength180, CommercialLength(180).Validate())
	test.expect(CommercialLength180, CommercialLength(200).Validate())
}

func TestStreamGetThumbnailUrl(t *testing.T) {
	stream := &Stream{
		ThumbnailURL: "https://static-cdn.jtvnw.net/previews-ttv/live_user_test-{width}x{height}.jpg",
	}

	test := formTest(t, "get stream thumbnail URL")

	// Test with default options
	url := stream.GetThumbnailUrl(nil)
	test.expect("https://static-cdn.jtvnw.net/previews-ttv/live_user_test-1920x1080.jpg", url)

	// Test with custom dimensions
	url = stream.GetThumbnailUrl(&ThumbnailUrlOptions{Width: 640, Height: 360})
	test.expect("https://static-cdn.jtvnw.net/previews-ttv/live_user_test-640x360.jpg", url)

	// Test with only width
	url = stream.GetThumbnailUrl(&ThumbnailUrlOptions{Width: 1280})
	test.expect("https://static-cdn.jtvnw.net/previews-ttv/live_user_test-1280x1080.jpg", url)

	// Test with only height
	url = stream.GetThumbnailUrl(&ThumbnailUrlOptions{Height: 720})
	test.expect("https://static-cdn.jtvnw.net/previews-ttv/live_user_test-1920x720.jpg", url)
}

// Helper function already defined in util_test.go but we need it here too
func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
