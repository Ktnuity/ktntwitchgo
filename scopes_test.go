package ktntwitchgo

import "testing"

func TestScopeString(t *testing.T) {
	test := formTest(t, "convert scope to string")

	test.expect("analytics:read:extensions", ScopeAnalyticsReadExtensions.String())
	test.expect("bits:read", ScopeBitsRead.String())
	test.expect("user:bot", ScopeUserBot.String())
}

func TestScopeIsValid(t *testing.T) {
	test := formTest(t, "validate scope")

	// Test all valid scopes
	test.expect(true, ScopeAnalyticsReadExtensions.IsValid())
	test.expect(true, ScopeAnalyticsReadGames.IsValid())
	test.expect(true, ScopeBitsRead.IsValid())
	test.expect(true, ScopeChannelEditCommercial.IsValid())
	test.expect(true, ScopeChannelReadHypeTrain.IsValid())
	test.expect(true, ScopeChannelReadSubscriptions.IsValid())
	test.expect(true, ScopeChannelReadStreamKey.IsValid())
	test.expect(true, ScopeChannelBot.IsValid())
	test.expect(true, ScopeClipsEdit.IsValid())
	test.expect(true, ScopeUserEdit.IsValid())
	test.expect(true, ScopeUserEditBroadcast.IsValid())
	test.expect(true, ScopeUserEditFollows.IsValid())
	test.expect(true, ScopeUserReadBroadcast.IsValid())
	test.expect(true, ScopeUserReadEmail.IsValid())
	test.expect(true, ScopeUserReadChat.IsValid())
	test.expect(true, ScopeUserWriteChat.IsValid())
	test.expect(true, ScopeUserBot.IsValid())
	test.expect(true, ScopeModerationRead.IsValid())

	// Test invalid scopes
	test.expect(false, Scope("invalid:scope").IsValid())
	test.expect(false, Scope("").IsValid())
	test.expect(false, Scope("random").IsValid())
}

func TestAllScopes(t *testing.T) {
	scopes := AllScopes()

	if len(scopes) != 18 {
		t.Errorf("Expected 18 scopes, got %d", len(scopes))
	}

	// Verify all scopes are present
	expectedScopes := []Scope{
		ScopeAnalyticsReadExtensions,
		ScopeAnalyticsReadGames,
		ScopeBitsRead,
		ScopeChannelEditCommercial,
		ScopeChannelReadHypeTrain,
		ScopeChannelReadSubscriptions,
		ScopeChannelReadStreamKey,
		ScopeChannelBot,
		ScopeClipsEdit,
		ScopeUserEdit,
		ScopeUserEditBroadcast,
		ScopeUserEditFollows,
		ScopeUserReadBroadcast,
		ScopeUserReadEmail,
		ScopeUserReadChat,
		ScopeUserWriteChat,
		ScopeUserBot,
		ScopeModerationRead,
	}

	for i, expected := range expectedScopes {
		if scopes[i] != expected {
			t.Errorf("Expected scope at index %d to be %s, got %s", i, expected, scopes[i])
		}
	}
}

func TestScopesByCategory(t *testing.T) {
	categories := ScopesByCategory()

	// Test analytics category
	if len(categories["analytics"]) != 2 {
		t.Errorf("Expected 2 analytics scopes, got %d", len(categories["analytics"]))
	}

	// Test bits category
	if len(categories["bits"]) != 1 {
		t.Errorf("Expected 1 bits scope, got %d", len(categories["bits"]))
	}

	// Test channel category
	if len(categories["channel"]) != 5 {
		t.Errorf("Expected 5 channel scopes, got %d", len(categories["channel"]))
	}

	// Test clips category
	if len(categories["clips"]) != 1 {
		t.Errorf("Expected 1 clips scope, got %d", len(categories["clips"]))
	}

	// Test user category
	if len(categories["user"]) != 8 {
		t.Errorf("Expected 8 user scopes, got %d", len(categories["user"]))
	}

	// Test moderation category
	if len(categories["moderation"]) != 1 {
		t.Errorf("Expected 1 moderation scope, got %d", len(categories["moderation"]))
	}
}

func TestHasScope(t *testing.T) {
	test := formTest(t, "check if scope exists in slice")

	scopes := []Scope{
		ScopeUserReadChat,
		ScopeUserWriteChat,
		ScopeBitsRead,
	}

	test.expect(true, HasScope(scopes, ScopeUserReadChat))
	test.expect(true, HasScope(scopes, ScopeUserWriteChat))
	test.expect(true, HasScope(scopes, ScopeBitsRead))
	test.expect(false, HasScope(scopes, ScopeUserBot))
	test.expect(false, HasScope(scopes, ScopeChannelBot))

	// Test empty slice
	test.expect(false, HasScope([]Scope{}, ScopeUserReadChat))
}

func TestScopesToStrings(t *testing.T) {
	scopes := []Scope{
		ScopeUserReadChat,
		ScopeUserWriteChat,
		ScopeBitsRead,
	}

	strings := ScopesToStrings(scopes)

	if len(strings) != 3 {
		t.Errorf("Expected 3 strings, got %d", len(strings))
	}

	test := formTest(t, "convert scopes to strings")
	test.expect("user:read:chat", strings[0])
	test.expect("user:write:chat", strings[1])
	test.expect("bits:read", strings[2])

	// Test empty slice
	emptyStrings := ScopesToStrings([]Scope{})
	if len(emptyStrings) != 0 {
		t.Errorf("Expected 0 strings, got %d", len(emptyStrings))
	}
}

func TestStringsToScopes(t *testing.T) {
	strings := []string{
		"user:read:chat",
		"user:write:chat",
		"bits:read",
	}

	scopes := StringsToScopes(strings)

	if len(scopes) != 3 {
		t.Errorf("Expected 3 scopes, got %d", len(scopes))
	}

	test := formTest(t, "convert strings to scopes")
	test.expect(ScopeUserReadChat, scopes[0])
	test.expect(ScopeUserWriteChat, scopes[1])
	test.expect(ScopeBitsRead, scopes[2])

	// Test empty slice
	emptyScopes := StringsToScopes([]string{})
	if len(emptyScopes) != 0 {
		t.Errorf("Expected 0 scopes, got %d", len(emptyScopes))
	}
}

func TestValidateScopes(t *testing.T) {
	test := formTest(t, "validate scopes")

	// All valid scopes
	validScopes := []Scope{
		ScopeUserReadChat,
		ScopeUserWriteChat,
		ScopeBitsRead,
	}
	test.expect(true, ValidateScopes(validScopes))

	// Mixed valid and invalid
	mixedScopes := []Scope{
		ScopeUserReadChat,
		Scope("invalid:scope"),
		ScopeBitsRead,
	}
	test.expect(false, ValidateScopes(mixedScopes))

	// All invalid
	invalidScopes := []Scope{
		Scope("invalid:scope"),
		Scope("another:invalid"),
	}
	test.expect(false, ValidateScopes(invalidScopes))

	// Empty slice
	test.expect(true, ValidateScopes([]Scope{}))
}

func TestFilterValidScopes(t *testing.T) {
	mixedScopes := []Scope{
		ScopeUserReadChat,
		Scope("invalid:scope"),
		ScopeUserWriteChat,
		Scope("another:invalid"),
		ScopeBitsRead,
	}

	validScopes := FilterValidScopes(mixedScopes)

	if len(validScopes) != 3 {
		t.Errorf("Expected 3 valid scopes, got %d", len(validScopes))
	}

	test := formTest(t, "filter valid scopes")
	test.expect(ScopeUserReadChat, validScopes[0])
	test.expect(ScopeUserWriteChat, validScopes[1])
	test.expect(ScopeBitsRead, validScopes[2])

	// Test all invalid
	allInvalid := []Scope{
		Scope("invalid:scope"),
		Scope("another:invalid"),
	}
	filteredInvalid := FilterValidScopes(allInvalid)
	if len(filteredInvalid) != 0 {
		t.Errorf("Expected 0 valid scopes, got %d", len(filteredInvalid))
	}

	// Test empty slice
	emptyFiltered := FilterValidScopes([]Scope{})
	if len(emptyFiltered) != 0 {
		t.Errorf("Expected 0 scopes, got %d", len(emptyFiltered))
	}
}

func TestFilterInvalidScopes(t *testing.T) {
	mixedScopes := []Scope{
		ScopeUserReadChat,
		Scope("invalid:scope"),
		ScopeUserWriteChat,
		Scope("another:invalid"),
		ScopeBitsRead,
	}

	invalidScopes := FilterInvalidScopes(mixedScopes)

	if len(invalidScopes) != 2 {
		t.Errorf("Expected 2 invalid scopes, got %d", len(invalidScopes))
	}

	test := formTest(t, "filter invalid scopes")
	test.expect(Scope("invalid:scope"), invalidScopes[0])
	test.expect(Scope("another:invalid"), invalidScopes[1])

	// Test all valid
	allValid := []Scope{
		ScopeUserReadChat,
		ScopeUserWriteChat,
	}
	filteredValid := FilterInvalidScopes(allValid)
	if len(filteredValid) != 0 {
		t.Errorf("Expected 0 invalid scopes, got %d", len(filteredValid))
	}

	// Test empty slice
	emptyFiltered := FilterInvalidScopes([]Scope{})
	if len(emptyFiltered) != 0 {
		t.Errorf("Expected 0 scopes, got %d", len(emptyFiltered))
	}
}
