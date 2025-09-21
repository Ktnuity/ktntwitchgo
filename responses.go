package ktntwitchgo

type Pagination struct {
	Cursor				string			`json:"cursor"`
}

type APIBaseResponse struct {
	Total				*int			`json:"total,omitempty"`
	Pagination			*Pagination		`json:"pagination,omitempty"`
}

type ObseleteResponse struct {
	Error				string			`json:"error"`
	Status				int				`json:"status"`
	Message				string			`json:"message"`
}

type APIGameResponse struct {
	APIBaseResponse
	Data				[]Game			`json:"data"`
}

type APIUserResponse struct {
	APIBaseResponse
	Data				[]User			`json:"data"`
}

type APIChannelResponse struct {
	APIBaseResponse
	Data				[]Channel		`json:"data"`
}

type APIChannelInfoResponse struct {
	Data				[]ChannelInfo	`json:"data"`
}

type APIFollowResponse = ObseleteResponse

type APIStreamResponse struct {
	APIBaseResponse
	Data				[]Stream		`json:"data"`
}

type APIStreamMarkerResponse struct {
	APIBaseResponse
	Data				[]StreamMarker	`json:"data"`
}

type APIStreamKeyResponse struct {
	Data				[]StreamKey		`json:"data"`
}

type APIVideoResponse struct {
	APIBaseResponse
	Data				[]Video			`json:"data"`
}

type APIClipsResponse struct {
	APIBaseResponse
	Data				[]Clip			`json:"data"`
}

type APITagResponse = ObseleteResponse

type APISubResponse struct {
	APIBaseResponse
	Data				[]Sub			`json:"data"`
}

type APIBanResponse struct {
	APIBaseResponse
	Data				[]Ban			`json:"data"`
}

type APIExtensionTransactionResponse struct {
	APIBaseResponse
	Data				[]ExtensionTransaction `json:"data"`
}

type APIExtensionResponse struct {
	Data				[]Extension		`json:"data"`
}

type APIActiveUserExtensionResponse struct {
	Data				[]ActiveExtension	`json:"data"`
}

type APICheermoteResponse struct {
	Data				[]Cheermote		`json:"data"`
}

type APIEmotesResponse struct {
	Data				[]Emote			`json:"data"`
	Template			string			`json:"template"`
}

type APICreateClipResponse struct {
	Data				[]CreatedClip	`json:"data"`
}

type APIModeratorResponse struct {
	APIBaseResponse
	Data				[]Moderator		`json:"data"`
}

type APICodeStatusResponse struct {
	Data				[]CodeStatus	`json:"data"`
}

type APIBitsLeaderboardResponse struct {
	Data				[]BitsPosition	`json:"data"`
	DateRage			DateRange		`json:"date_range"`
	Total				int				`json:"total"`
}

type APICommercialResponse struct {
	Data				[]Commercial	`json:"data"`
}

type APIBadgesResponse struct {
	Data				[]Badge			`json:"data"`
}

type APIIngestsResponse struct {
	Ingests				[]Ingest		`json:"ingests"`
}

type APIMessageResponse struct {
	Data				[]Message		`json:"data"`
}

type ResponseError struct {
	Error				string			`json:"error"`
	Status				int				`json:"status"`
	Message				string			`json:"message"`
}

func (r *ResponseError) IsError() bool {
	return r.Status >= 400
}

func NewObsoleteResponse() ObseleteResponse {
	return ObseleteResponse{
		Error:		"Gone",
		Status:		410,
		Message: 	"This API is not available.",
	}
}
