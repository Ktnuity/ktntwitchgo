package ktntwitchgo

type Scope string

const (
	// Analytics scopes
	ScopeAnalyticsReadExtensions		Scope = "analytics:read:extensions"
	ScopeAnalyticsReadGames				Scope = "analytics:read:games"

	// Bits scopes
	ScopeBitsRead						Scope = "bits:read"

	// Channel scopes
	ScopeChannelEditCommercial			Scope = "channel:edit:commercial"
	ScopeChannelReadHypeTrain			Scope = "channel:read:hype_train"
	ScopeChannelReadSubscriptions		Scope = "channel:read:subscriptions"
	ScopeChannelReadStreamKey			Scope = "channel:read:stream_key"
	ScopeChannelBot						Scope = "channel:bot"

	// Clips scopes
	ScopeClipsEdit						Scope = "clips:edit"

	// User scopes
	ScopeUserEdit						Scope = "user:edit"
	ScopeUserEditBroadcast				Scope = "user:edit:broadcast"
	ScopeUserEditFollows				Scope = "user:edit:follows"
	ScopeUserReadBroadcast				Scope = "user:read:broadcast"
	ScopeUserReadEmail					Scope = "user:read:email"
	ScopeUserReadChat					Scope = "user:read:chat"
	ScopeUserWriteChat					Scope = "user:write:chat"
	ScopeUserBot						Scope = "user:bot"

	// Moderation scopes
	ScopeModerationRead					Scope = "moderation:read"
)

func (s Scope) String() string {
	return string(s)
}

func (s Scope) IsValid() bool {
	switch s {
	case 	ScopeAnalyticsReadExtensions,
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
			ScopeModerationRead:
		return true
	default:
		return false
	}
}

func AllScopes() []Scope {
	return []Scope{
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
}

func ScopesByCategory() map[string][]Scope {
	return map[string][]Scope{
		"analytics": {
			ScopeAnalyticsReadExtensions,
			ScopeAnalyticsReadGames,
		},
		"bits": {
			ScopeBitsRead,
		},
		"channel": {
			ScopeChannelEditCommercial,
			ScopeChannelReadHypeTrain,
			ScopeChannelReadSubscriptions,
			ScopeChannelReadStreamKey,
			ScopeChannelBot,
		},
		"clips": {
			ScopeClipsEdit,
		},
		"user": {
			ScopeUserEdit,
			ScopeUserEditBroadcast,
			ScopeUserEditFollows,
			ScopeUserReadBroadcast,
			ScopeUserReadEmail,
			ScopeUserReadChat,
			ScopeUserWriteChat,
			ScopeUserBot,
		},
		"moderation": {
			ScopeModerationRead,
		},
	}
}

func HasScope(scopes []Scope, target Scope) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}
	return false
}

func ScopesToStrings(scopes []Scope) []string {
	result := make([]string, len(scopes))
	for i, scope := range scopes {
		result[i] = string(scope)
	}
	return result
}

func StringsToScopes(scopes []string) []Scope {
	result := make([]Scope, len(scopes))
	for i, scope := range scopes {
		result[i] = Scope(scope)
	}
	return result
}

func ValidateScopes(scopes []Scope) bool {
	for _, scope := range scopes {
		if !scope.IsValid() {
			return false
		}
	}
	return true
}

func FilterValidScopes(scopes []Scope) []Scope {
	result := make([]Scope, 0, len(scopes))
	for _, scope := range scopes {
		if scope.IsValid() {
			result = append(result, scope)
		}
	}
	return result
}

func FilterInvalidScopes(scopes []Scope) []Scope {
	result := make([]Scope, 0, len(scopes))
	for _, scope := range scopes {
		if !scope.IsValid() {
			result = append(result, scope)
		}
	}
	return result
}
