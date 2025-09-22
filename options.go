package ktntwitchgo

type TwitchApiConfig struct {
	ClientID			string			`json:"client_id"`
	ClientSecret		string			`json:"client_secret"`
	Scopes				[]Scope			`json:"scopes,omitempty"`
	AccessToken			*string			`json:"access_token,omitempty"`
	RefreshToken		*string			`json:"refresh_token,omitempty"`
	RedirectURI			*string			`json:"redirect_uri,omitempty"`
	ThrowRatelimitErrors *bool			`json:"throw_ratelimit_errors,omitempty"`
}

type GetFollowsOptions struct {
	After				*string			`json:"after,omitempty"`
	First				*int			`json:"first,omitempty"`
	FromID				*string 		`json:"from_id,omitempty"`
	ToID				*string 		`json:"to_id,omitempty"`
}

type BaseOptions struct {
	First				*int			`json:"first,omitempty"`
	After				*string			`json:"after,omitempty"`
	Before				*string			`json:"before,omitempty"`
}

type GetStreamsOptions struct {
	BaseOptions
	GameID				[]string		`json:"game_id,omitempty"`
	Language			[]string		`json:"language,omitempty"`
	Channels			[]string		`json:"channels,omitempty"`
	Channel				*string			`json:"channel,omitempty"`
}

type BaseClipsOptions struct {
	BaseOptions
	EndedAt				*string			`json:"ended_at,omitempty"`
	StartedAt			*string			`json:"started_at,omitempty"`
}

type ClipsBroadcasterIdOptions struct {
	BaseClipsOptions
	BroadcasterID		string			`json:"broadcaster_id"`
}

type ClipsGameIdOptions struct {
	BaseClipsOptions
	GameID				string			`json:"game_id"`
}

type ClipsIdOptions struct {
	BaseClipsOptions
	ID					[]string		`json:"id"`
}

type VideoPeriod string
const (
	VideoPeriodAll		VideoPeriod = "all"
	VideoPeriodDay		VideoPeriod = "day"
	VideoPeriodWeek		VideoPeriod = "week"
	VideoPeriodMonth	VideoPeriod = "month"
)

type VideoSort string
const (
	VideoSortTime		VideoSort = "time"
	VideoSortTrending	VideoSort = "trending"
	VideoSortViews		VideoSort = "views"
)

type VideoType string
const (
	VideoTypeAll		VideoType = "all"
	VideoTypeUpload		VideoType = "upload"
	VideoTypeArchive	VideoType = "archive"
	VideoTypeHighlight	VideoType = "highlight"
)

type GetVideosOptions struct {
	BaseOptions
	ID					[]string		`json:"id,omitempty"`
	UserID				*string			`json:"user_id,omitempty"`
	GameID				*string			`json:"game_id,omitempty"`
	Language			*string			`json:"language,omitempty"`
	Period 				*VideoPeriod	`json:"period,omitempty"`
	Sort				*VideoSort		`json:"sort,omitempty"`
	Type				*VideoType		`json:"type,omitempty"`
}

type GetAllStreamTagsOptions struct {
	After				*string			`json:"after,omitempty"`
	First				*int			`json:"first,omitempty"`
	TagID				[]string		`json:"tag_id,omitempty"`
}

type GetStreamTagsOptions struct {
	BroadcasterID string `json:"broadcaster_id"`
}

type BitsLeaderboardPeriod string
const (
	BitsLeaderboardPeriodDay	BitsLeaderboardPeriod = "day"
	BitsLeaderboardPeriodWeek	BitsLeaderboardPeriod = "week"
	BitsLeaderboardPeriodMonth	BitsLeaderboardPeriod = "month"
	BitsLeaderboardPeriodYear	BitsLeaderboardPeriod = "year"
	BitsLeaderboardPeriodAll	BitsLeaderboardPeriod = "all"
)

type GetBitsLeaderboardOptions struct {
	Count				*int					`json:"count,omitempty"`
	Period				*BitsLeaderboardPeriod	`json:"period,omitempty"`
	StartedAt			*string					`json:"started_at,omitempty"`
	UserID				*string					`json:"user_id,omitempty"`
}

type GetSubsOptions struct {
	BroadcasterID		string			`json:"broadcaster_id"`
	UserID				[]string		`json:"user_id,omitempty"`
}

type GetChannelInfoOptions struct {
	BroadcasterID		[]string		`json:"broadcaster_id"`
}

type SearchOptions struct {
	Query				string			`json:"query"`
	First				*int			`json:"first,omitempty"`
	After				*string			`json:"after,omitempty"`
}

type SearchChannelsOptions struct {
	SearchOptions
	LiveOnly			*bool			`json:"live_only,omitempty"`
}

type SearchCategoriesOptions = SearchOptions
type GetBannedUsersOptions struct {
	BroadcasterID		string			`json:"broadcaster_id"`
	UserID				[]string		`json:"user_id,omitempty"`
	After				*string			`json:"after,omitempty"`
	Before				*string			`json:"before,omitempty"`
}

type GetExtensionTransactionsOptions struct {
	ExtensionID			string			`json:"extension_id"`
	ID					[]string		`json:"id,omitempty"`
	After				*string			`json:"after,omitempty"`
	First				*string			`json:"first,omitempty"`
}

type GetCheermotesOptions struct {
	BroadcasterID		*string			`json:"broadcaster_id,omitempty"`
}

type GetStreamKeyOptions struct {
	BroadcasterID		string			`json:"broadcaster_id"`
}

type GetStreamMarkerUserIdOptions struct {
	BaseOptions
	UserID				string			`json:"user_id"`
}

type GetStreamMarkerVideoIdOptions struct {
	BaseOptions
	VideoID				string			`json:"video_id"`
}

type CreateUserFollowsOptions struct {
	FromID				string			`json:"from_id"`
	ToID				string			`json:"to_id"`
	AllowNotifications	*bool			`json:"allow_notifications,omitempty"`
}

type DeleteUserFollowsOptions struct {
	FromID				string			`json:"from_id"`
	ToID				string			`json:"to_id"`
}

type GetUserActiveExtensionsOptions struct {
	UserID				*string			`json:"user_id,omitempty"`
}

type ModifyChannelInformationOptions struct {
	BroadcasterID		string 			`json:"broadcaster_id"`
	GameID				*string			`json:"game_id,omitempty"`
	BroadcasterLanguage *string			`json:"broadcaster_language,omitempty"`
	Title				*string			`json:"title,omitempty"`
}

type GetCodeStatusOptions struct {
	Code				[]string		`json:"code"`
	UserID				string			`json:"user_id"`
}

type ReplaceStreamTagsOptions struct {
	TagIDs				[]string		`json:"tag_ids,omitempty"`
	BroadcasterID		string			`json:"broadcaster_id"`
}

type UpdateUserOptions struct {
	Description			*string			`json:"description,omitempty"`
}

type CreateClipOptions struct {
	BroadcasterID		string			`json:"broadcaster_id"`
	HasDelay			*bool			`json:"has_delay,omitempty"`
}

type StartCommercialOptions struct {
	BroadcasterID		string				`json:"broadcaster_id"`
	Length				CommercialLength	`json:"length"`
}

type GetModeratorsOptions struct {
	BroadcasterID		string			`json:"broadcaster_id"`
	UserID				[]string		`json:"user_id,omitempty"`
	After				*string			`json:"after,omitempty"`
}

type GetChannelChatBadges struct {
	BroadcasterID		string			`json:"broadcaster_id"`
}

type SendChatMessageOptions struct {
	BroadcasterID		string			`json:"broadcaster_id"`
	SenderID			string			`json:"sender_id"`
	Message				string			`json:"message"`
	ReplyParentMessageID *string		`json:"reply_parent_message_id,omitempty"`
	ForSourceOnly		*bool			`json:"for_source_only,omitempty"`
}
