package ktntwitchgo

import (
	"strconv"
	"strings"
)

type User struct {
	BroadcasterType		string		`json:"broadcaster_type"`
	Description			string		`json:"description"`
	DisplayName			string		`json:"display_name"`
	Email				string		`json:"email"`
	ID					string		`json:"id"`
	Login				string		`json:"login"`
	OfflineImageURL		string		`json:"offline_image_url"`
	ProfileImageURL		string		`json:"profile_image_url"`
	Type				string		`json:"type"`
	ViewCount			string		`json:"view_count"`
}

type Channel struct {
	GameID				string		`json:"game_id"`
	ID					string		`json:"id"`
	DisplayName			string		`json:"display_name"`
	BroadcasterLanguage	string		`json:"broadcaster_language"`
	Title				string		`json:"title"`
	ThumbnailURL		string		`json:"thumbnail_url"`
	IsLive				bool		`json:"is_live"`
	StartedAt			string		`json:"string"`
	TagIDs				string		`json:"tag_ids"`
}

type ChannelInfo struct {
	BroadcasterID		string		`json:"broadcaster_id"`
	BroadcasterName		string		`json:"broadcaster_name"`
	GameName			string		`json:"game_name"`
	GameID				string		`json:"game_id"`
	BroadcasterLanguage string		`json:"broadcaster_language"`
	Title				string		`json:"title"`
}

type Game struct {
	BoxArtURL			string		`json:"box_art_url"`
	ID					string		`json:"id"`
	Name				string		`json:"name"`
}

type Follow struct {
	FollowedAt			string		`json:"followed_at"`
	FromID				string		`json:"from_id"`
	FromName			string		`json:"from_name"`
	ToID				string		`json:"to_id"`
	ToName				string		`json:"to_name"`
}

type ThumbnailUrlOptions struct {
	Width				int
	Height				int
}

type Stream struct {
	GameID				string		`json:"game_id"`
	GameName			string		`json:"game_name"`
	ID					string		`json:"id"`
	IsMature			bool		`json:"is_mature"`
	Language			string		`json:"language"`
	StartedAt			string		`json:"started_at"`
	TagIDs				string		`json:"tag_ids"`
	ThumbnailURL		string		`json:"thumbnail_url"`
	Title				string		`json:"title"`
	Type				string		`json:"type"`
	UserID				string		`json:"user_id"`
	UserLogin			string		`json:"user_login"`
	UserName			string		`json:"user_name"`
	ViewerCount			int			`json:"viewer_count"`
}
func (s *Stream) GetThumbnailUrl(options *ThumbnailUrlOptions) string {
	width := 1920
	height := 1080

	if options != nil {
		if options.Width > 0 {
			width = options.Width
		}

		if options.Height > 0 {
			height = options.Height
		}
	}

	url := s.ThumbnailURL
	url = strings.ReplaceAll(url, "{width}", strconv.Itoa(width))
	url = strings.ReplaceAll(url, "{height}", strconv.Itoa(height))

	return url
}

type StreamMarker struct {
	ID					string		`json:"id"`
	CreatedAt			string		`json:"created_at"`
	Description			string		`json:"description"`
	PositionSeconds		int			`json:"position_seconds"`
	URL					string		`json:"URL"`
	UserID				string		`json:"user_id"`
	UserName			string		`json:"user_name"`
	VideoID				string		`json:"video_id"`
}

type StreamKey struct {
	Key					string		`json:"stream_key"`
}

type VideoViewable string
const (
	VideoViewablePublic		VideoViewable = "public"
	VideoViewablePrivate	VideoViewable = "public"
)

type Video struct {
	CreatedAt			string			`json:"created_at"`
	Description			string			`json:"description"`
	Duration			string			`json:"duration"`
	ID					string			`json:"id"`
	Language			string			`json:"language"`
	PublishedAt			string			`json:"published_at"`
	ThumbnailURL		string			`json:"thumbnail_url"`
	Title				string			`json:"title"`
	Type				VideoType		`json:"type"`
	URL					string			`json:"url"`
	UserID				string			`json:"user_id"`
	ViewCount			int 			`json:"view_count"`
	Viewable			VideoViewable 	`json:"viewable"`
}

type Clip struct {
	BroadcasterID		string		`json:"broadcaster_id"`
	BroadcasterName		string		`json:"broadcaster_name"`
	CreatedAt			string		`json:"created_at"`
	CreatorID			string		`json:"creator_id"`
	CreatorName			string		`json:"creator_name"`
	EmbedURL			string		`json:"embed_url"`
	GameID				string		`json:"game_id"`
	ID					string		`json:"id"`
	Language			string		`json:"language"`
	ThumbnailURL		string		`json:"thumbnail_url"`
	Title				string		`json:"title"`
	URL					string		`json:"url"`
	VideoID				string		`json:"video_id"`
	ViewCount			int			`json:"view_count"`
}

type Tag struct {
	TagID						string				`json:"tag_id"`
	IsAuto						bool				`json:"is_auto"`
	LocalizationNames			map[string]string	`json:"localization_names"`
	LocalizationDescriptions 	map[string]string	`json:"localization_descriptions"`
}

type SubscriptionTier			string
const (
	SubscriptionTier1			SubscriptionTier 	= "1000"
	SubscriptionTier2			SubscriptionTier	= "2000"
	SubscriptionTier3			SubscriptionTier	= "3000"
)

type Sub struct {
	BroadcasterID		string				`json:"broadcaster_id"`
	BroadcasterName		string				`json:"broadcaster_name"`
	IsGift				bool				`json:"is_gift"`
	Tier				SubscriptionTier	`json:"tier"`
	PlanName			string				`json:"plan_name"`
	UserID				string				`json:"user_id"`
	UserName			string				`json:"user_name"`
}

type BitsPosition struct {
	UserID				string		`json:"user_id"`
	UserName			string		`json:"user_name"`
	Rank				int			`json:"rank"`
	Score				int			`json:"score"`
}

type Ban struct {
	UserID				string		`json:"user_id"`
	UserName			string		`json:"user_name"`
	ExpiresAt			string		`json:"expires_at"`
}

type Moderator struct {
	UserID				string		`json:"user_id"`
	UserName			string		`json:"user_name"`
}

type CodeRedemptionStatus		string
const (
	StatusSuccessfullyRedeemed	CodeRedemptionStatus = "SUCCESSFULLY_REDEEMED"
	StatusAlreadyClaimed 		CodeRedemptionStatus = "ALREADY_CLAIMED"
	StatusExpired				CodeRedemptionStatus = "EXPIRED"
	StatusUserNotEligible		CodeRedemptionStatus = "USER_NOT_ELIGIBLE"
	StatusNotFound				CodeRedemptionStatus = "NOT_FOUND"
	StatusInactive				CodeRedemptionStatus = "INACTIVE"
	StatusUnused				CodeRedemptionStatus = "UNUSED"
	StatusIncorrectFormat		CodeRedemptionStatus = "INCORRECT_FORMAT"
	StatusInternalError			CodeRedemptionStatus = "INTERNAL_ERROR"
)

type CodeStatus struct {
	Code				string					`json:"code"`
	Status				CodeRedemptionStatus	`json:"status"`
}

type CreatedClip struct {
	ID					string		`json:"id"`
	EditURL				string		`json:"edit_url"`
}

type ProductData struct {
	Domain				*string		`json:"domain,omitempty"`
	Broadcast			*bool		`json:"broadcast,omitempty"`
	Expiration			*string		`json:"expiration,omitempty"`
	SKU					string		`json:"sku"`
	Cost				Cost		`json:"cost"`
	DisplayName			string		`json:"displayName"`
	InDevelopment		bool		`json:"inDevelopment"`
}

type Cost struct {
	Amount				int			`json:"amount"`
	Type				string		`json:"type"`
}

type ExtensionTransaction struct {
	ID					string		`json:"id"`
	Timestamp			string		`json:"timestamp"`
	BroadcasterID		string		`json:"broadcaster_id"`
	BroadcasterName		string		`json:"broadcaster_name"`
	UserID				string		`json:"user_id"`
	UserName			string		`json:"user_nme"`
	ProductType			string		`json:"product_type"`
	ProductData			ProductData	`json:"product_data"`
}

type Extension struct {
	CanActivate			bool		`json:"can_activate"`
	Type				[]string	`json:"type"`
	ID					string		`json:"id"`
	Name				string		`json:"name"`
	Version				string		`json:"version"`
}

type ActiveExtensionBase struct {
	Active				bool		`json:"active"`
	ID					string		`json:"id"`
	Version				string		`json:"version"`
	Name				string		`json:"name"`
}

type NotActive struct {
	Active				bool		`json:"active"`
}

type ExtensionComponent struct {
	ActiveExtensionBase
	X					int			`json:"x"`
	Y					int			`json:"y"`
}

type ActiveExtension struct {
	Panel				map[string]interface{}	`json:"panel"`
	Overlay				map[string]interface{}	`json:"overlay"`
	Component			map[string]interface{}	`json:"component"`
}

type CheermoteType string
const (
	CheermoteTypeGlobalFirstParty		CheermoteType = "global_first_party"
	CheermoteTypeGlobalThirdParty		CheermoteType = "global_third_party"
	CheermoteTypeChannelCustom			CheermoteType = "channel_custom"
	CheermoteTypeDisplayOnly			CheermoteType = "display_only"
	CheermoteTypeSponsored				CheermoteType = "sponsored"
)

type CheermoteTierID string
const (
	CheermoteTier1						CheermoteTierID = "1"
	CheermoteTier100					CheermoteTierID = "100"
	CheermoteTier500					CheermoteTierID = "500"
	CheermoteTier1000					CheermoteTierID = "1000"
	CheermoteTier5000					CheermoteTierID = "5000"
	CheermoteTier10k					CheermoteTierID = "10k"
	CheermoteTier100k					CheermoteTierID = "100k"
)

type CheermoteImages struct {
	Animated 			CheermoteImageSizes	`json:"animated"`
	Static				CheermoteImageSizes	`json:"static"`
}

type CheermoteImageSizes struct {
	Size1				string		`json:"1"`
	Size1_5				string		`json:"1.5"`
	Size2				string		`json:"2"`
	Size3				string		`json:"3"`
	Size4				string		`json:"4"`
}

type CheermoteTier struct {
	MinBits				int					`json:"min_bits"`
	ID					CheermoteTierID		`json:"id"`
	Color				string				`json:"color"`
	Images				CheermoteThemeImage	`json:"images"`
	CanCheer			bool				`json:"can_cheer"`
	ShowInBitsCard		bool				`json:"show_in_bits_card"`
}

type CheermoteThemeImage struct {
	Dark				CheermoteImages		`json:"dark"`
	Light				CheermoteImages		`json:"light"`
}

type Cheermote struct {
	Tiers				[]CheermoteTier		`json:"tiers"`
	Type				CheermoteType		`json:"type"`
	Order				int					`json:"order"`
	LastUpdated			string				`json:"last_updated"`
	IsCharitable		bool				`json:"is_charitable"`
}

type EmoteType			string
const (
	EmoteTypeBitsTier		EmoteType = "bitstier"
	EmoteTypeFollower		EmoteType = "follower"
	EmoteTypeSubscriptions	EmoteType = "subscriptions"
)

type EmoteFormat		string
const (
	EmoteFormatAnimated	EmoteFormat = "animated"
	EmoteFormatStatic	EmoteFormat = "static"
)

type EmoteScale			string
const (
	EmoteScale1_0		EmoteScale = "1.0"
	EmoteScale2_0		EmoteScale = "2.0"
	EmoteScale3_0		EmoteScale = "3.0"
)

type EmoteTheme			string
const (
	EmoteThemeLight		EmoteTheme = "light"
	EmoteThemeDark		EmoteTheme = "dark"
)

type EmoteImages struct {
	URL1x				string		`json:"url_1x"`
	URL2x				string		`json:"url_2x"`
	URL3x				string		`json:"url_3x"`
}

type Emote struct {
	ID					string 			`json:"id"`
	Name				string			`json:"name"`
	Images				EmoteImages		`json:"images"`
	Tier				string			`json:"tier"`
	EmoteType			EmoteType		`json:"emote_type"`
	EmoteSetID			string			`json:"emote_set_id"`
	Format				[]EmoteFormat	`json:"format"`
	Scale				[]EmoteScale	`json:"scale"`
	ThemeMode			[]EmoteTheme	`json:"theme_mode"`
}

type DateRange struct {
	StartedAt			string		`json:"started_at"`
	EndedAt				string		`json:"ended_at"`
}

type CommercialLength	int
const (
	CommercialLength30	CommercialLength = 30
	CommercialLength60	CommercialLength = 60
	CommercialLength90	CommercialLength = 90
	CommercialLength120	CommercialLength = 120
	CommercialLength150	CommercialLength = 150
	CommercialLength180	CommercialLength = 180
)

func (cl CommercialLength) IsValid() bool {
	return	cl == CommercialLength30 || cl == CommercialLength60 ||
			cl == CommercialLength90 || cl == CommercialLength120 ||
			cl == CommercialLength150 || cl == CommercialLength180
}

func (cl CommercialLength) String() string {
	switch cl {
	case CommercialLength30:
		return "30"
	case CommercialLength60:
		return "60"
	case CommercialLength90:
		return "90"
	case CommercialLength120:
		return "120"
	case CommercialLength150:
		return "150"
	case CommercialLength180:
		return "180"
	default:
		return "invalid"
	}
}

func (cl CommercialLength) Validate() CommercialLength {
	if cl <= CommercialLength30 {
		return CommercialLength30
	} else if cl <= CommercialLength60 {
		return CommercialLength60
	} else if cl <= CommercialLength90 {
		return CommercialLength90
	} else if cl <= CommercialLength120 {
		return CommercialLength120
	} else if cl <= CommercialLength150 {
		return CommercialLength150
	} else {
		return CommercialLength180
	}
}

type Commercial struct {
	Length				CommercialLength	`json:"length"`
	Message				string				`json:"message"`
	RetryAfter			int					`json:"retry_after"`
}

type BadgeVersion struct {
	ID					string				`json:"id"`
	ImageURL1x			string				`json:"image_url_1x"`
	ImageURL2x			string				`json:"image_url_2x"`
	ImageURL4x			string				`json:"image_url_4x"`
}

type Badge struct {
	SetID				string				`json:"set_id"`
	Versions			[]BadgeVersion		`json:"versions"`
}

type Ingest struct {
	ID					int					`json:"_id"`
	Availability		float64				`json:"availability"`
	Default				bool				`json:"default"`
	Name				string				`json:"name"`
	URLTemplate			string				`json:"url_template"`
	Priority			int					`json:"priority"`
}

type DropReason struct {
	Code				string				`json:"code"`
	Message				string				`json:"message"`
}

type Message struct {
	MessageID			string				`json:"message_id"`
	IsSent				bool				`json:"is_sent"`
	DropReason			*DropReason 		`json:"drop_reason,omitempty"`
}
