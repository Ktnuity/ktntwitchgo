package ktntwitchgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EventHandler func(data any)

type Client struct {
	clientSecret	string
	clientID		string

	user			*User
	accessToken		*string
	refreshToken	*string
	scopes			[]Scope
	redirectURI		*string

	throwRateLimitErrors bool

	baseURL			string
	ingestBaseURL	string
	httpClient		*http.Client
	refreshAttempts	int
	ready			bool

	eventHandlers	map[string][]EventHandler

	Verbose			bool
}

func CreateTwitchApi(config TwitchApiConfig) (*Client, *sync.WaitGroup) {
	client := &Client{
		clientSecret:			config.ClientSecret,
		clientID:				config.ClientID,
		accessToken: 			config.AccessToken,
		refreshToken: 			config.RefreshToken,
		scopes:					config.Scopes,
		redirectURI: 			config.RedirectURI,
		throwRateLimitErrors:	config.ThrowRatelimitErrors != nil && *config.ThrowRatelimitErrors,
		baseURL:				"https://api.twitch.tv/helix",
		ingestBaseURL:			"https://ingest.twitch.tv",
		refreshAttempts: 		0,
		ready:					false,
		eventHandlers: 			make(map[string][]EventHandler),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Verbose: false,
	}

	wg := client.initialize()
	return client, wg
}

func (c *Client) log(fmt string, data...any) {
	if !c.Verbose { return }

	if !strings.HasSuffix(fmt, "\n") {
		fmt += "\n"
	}
	log.Printf(fmt, data...)
}

func (c *Client) AddEventHandler(event string, handler EventHandler) {
	c.eventHandlers[event] = append(c.eventHandlers[event], handler)
}

func (c *Client) RemoveEventHandler(event string) {
	delete(c.eventHandlers, event)
}

func (c *Client) emit(event string, data any) {
	if handlers, exists := c.eventHandlers[event]; exists {
		for _, handler := range handlers {
			handler(data)
		}
	}
}

func (c *Client) initialize() *sync.WaitGroup {
	if c.accessToken != nil {
		var wg sync.WaitGroup
		wg.Go(func() {
			user, err := c.GetCurrentUser()
			if err == nil && user != nil {
				c.user = user
			}
		})
		return &wg
	}

	return nil
}

func (c *Client) error(message string) error {
	return fmt.Errorf("%s", message)
}

func (c *Client) getAppAccessToken(ctx context.Context) (*string, error) {
	data := map[string]string{
		"client_id":		c.clientID,
		"client_secret":	c.clientSecret,
		"grant_type":		"client_credentials",
	}

	if len(c.scopes) > 0 {
		scopeStrings := ScopesToStrings(c.scopes)
		data["scope"] = strings.Join(scopeStrings, " ")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://id.twitch.tv/oauth2/token", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error getting app access token. Expected JSON but got: %s", string(body))
	}

	if accessToken, ok := result["access_token"].(string); ok {
		return &accessToken, nil
	}

	return nil, fmt.Errorf("no access_token in response")
}

func (c *Client) refresh(ctx context.Context) error {
	valid, err := c.validate(ctx)
	if err != nil {
		return err
	}

	if valid {
		return nil
	}

	if c.refreshToken == nil {
		return c.error("refresh token is not set")
	}

	data := map[string]string{
		"client_id":		c.clientID,
		"client_secret":	c.clientSecret,
		"grant_type":		"refresh_token",
		"refresh_token":	url.QueryEscape(*c.refreshToken),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://id.twitch.tv/oauth2/token", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result AuthEvent
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.AccessToken != "" {
		c.accessToken = &result.AccessToken
	}
	if result.RefreshToken != "" {
		c.refreshToken = &result.RefreshToken
	}

	c.emit("refresh", result)

	if result.AccessToken == "" {
		c.refreshAttempts++
	}

	return nil
}

func (c *Client) validate(ctx context.Context) (bool, error) {
	if c.accessToken == nil {
		return false, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "OAuth " + *c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if message, ok := result["message"].(string); ok && message == "missing authorization token" {
		return false, c.error(message)
	}

	return resp.StatusCode == 200, nil
}

func (c *Client) get(ctx context.Context, endpoint string, apiType string) ([]byte, error) {
	c.log("GET: Endpoint(%s)", endpoint)
	if c.accessToken == nil {
		token, err := c.getAppAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		if token == nil {
			return nil, c.error("app access token could not be fetched. Please check your client_id and client_secret")
		}
		c.accessToken = token
	}

	var baseURL string
	switch apiType {
	case "helix":
		baseURL = c.baseURL
	case "ingest":
		baseURL = c.ingestBaseURL
	default:
		baseURL = c.baseURL
	}

	fullURL := baseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.clientID)
	req.Header.Set("Authorization", "Bearer " + *c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.handleRateLimit(resp.Header)

	if resp.StatusCode == 401 {
		if err := c.refresh(ctx); err != nil {
			return nil, err
		}
		return c.get(ctx, endpoint, apiType)
	}

	if resp.StatusCode == 429 {
		rateLimit := c.extractRateLimit(resp.Header)
		c.emit("ratelimit", rateLimit)

		if c.throwRateLimitErrors {
			return nil, &TwitchApiRateLimitError{RateLimit: rateLimit}
		}

		sleepTime := time.Duration(rateLimit.Reset) * time.Second - time.Duration(time.Now().Unix()) * time.Second
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}

		return c.get(ctx, endpoint, apiType)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) update(ctx context.Context, endpoint string, data any, method string) ([]byte, error) {
	if !strings.HasPrefix(endpoint, "/") {
		return nil, c.error("endpoint must start with a '/' (forward slash)")
	}

	c.log("UPDATE: Endpoint(%s)", endpoint)

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	fullURL := c.baseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), fullURL, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer " + *c.accessToken)
	req.Header.Set("Client-ID", c.clientID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.handleRateLimit(resp.Header)

	if resp.StatusCode == 401 {
		if err := c.refresh(ctx); err != nil {
			return nil, err
		}
		return c.update(ctx, endpoint, data, method)
	}

	if resp.StatusCode == 429 {
		rateLimit := c.extractRateLimit(resp.Header)
		c.emit("ratelimit", rateLimit)

		if c.throwRateLimitErrors {
			return nil, &TwitchApiRateLimitError{RateLimit: rateLimit}
		}

		sleepTime := time.Duration(rateLimit.Reset) * time.Second - time.Duration(time.Now().Unix()) * time.Second
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
		return c.update(ctx, endpoint, data, method)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) post(ctx context.Context, endpoint string, data any) ([]byte, error) {
	c.log("POST: Endpoint(%s)", endpoint)
	return c.update(ctx, endpoint, data, "post")
}

func (c *Client) put(ctx context.Context, endpoint string, data any) ([]byte, error) {
	c.log("PUT: Endpoint(%s)", endpoint)
	return c.update(ctx, endpoint, data, "put")
}

func (c *Client) patch(ctx context.Context, endpoint string, data any) ([]byte, error) {
	c.log("PATCH: Endpoint(%s)", endpoint)
	return c.update(ctx, endpoint, data, "patch")
}

func (c *Client) delete(ctx context.Context, endpoint string, data any) ([]byte, error) {
	c.log("DELETE: Endpoint(%s)", endpoint)
	return c.update(ctx, endpoint, data, "delete")
}

func (c *Client) handleRateLimit(headers http.Header) {
	rateLimit := c.extractRateLimit(headers)
	c.emit("ratelimitpoll", rateLimit)
}

func (c *Client) extractRateLimit(headers http.Header) TwitchApiRateLimit {
	limit, _ := strconv.Atoi(headers.Get("Ratelimit-Limit"))
	remaining, _ := strconv.Atoi(headers.Get("Ratelimit-Remaining"))
	reset, _ := strconv.Atoi(headers.Get("Ratelimit-Reset"))

	return TwitchApiRateLimit{
		Limit:		limit,
		Remaining:	remaining,
		Reset:		reset,
	}
}

func (c *Client) hasScope(scope Scope) bool {
	return slices.Contains(c.scopes, scope)
}

func parseMixedParam(values any, stringKey, numericKey string) string {
	var query []string

	addToQuery := func (value any) {
		switch v := value.(type) {
		case int:
			query = append(query, fmt.Sprintf("%s=%d", numericKey, v))
		case string:
			if !isNumber(v) {
				query = append(query, fmt.Sprintf("%s=%s", stringKey, v))
			} else {
				query = append(query, fmt.Sprintf("%s=%s", numericKey, v))
			}
		}
	}

	switch v := values.(type) {
	case []string:
		for _, val := range v {
			addToQuery(val)
		}
	case []int:
		for _, val := range v {
			addToQuery(val)
		}
	case int, string:
		addToQuery(v)
	}

	return strings.Join(query, "&")
}

func chooseKey(when bool, ifTrue, ifFalse string) string {
	if when {
		return ifTrue
	}

	return ifFalse
}

func getJSONFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return ""
	}

	name := strings.Split(tag, ",")[0]
	return name
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

func parseOptions[T any](options *T) string {
	if options == nil {
		return ""
	}

	v := reflect.ValueOf(options).Elem()
	t := v.Type()

	var parts []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Handle embedded structs
		if fieldType.Anonymous && field.Kind() == reflect.Struct {
			// Recursively parse embedded struct fields
			for j := 0; j < field.NumField(); j++ {
				embeddedField := field.Field(j)
				embeddedFieldType := field.Type().Field(j)

				if embeddedField.Kind() == reflect.Pointer {
					embeddedField = embeddedField.Elem()
				}

				if !embeddedField.IsValid() || isZeroValue(embeddedField) {
					continue
				}

				key := getJSONFieldName(embeddedFieldType)
				if key == "" || key == "-" {
					continue
				}

				if embeddedField.Kind() == reflect.Slice {
					for k := 0; k < embeddedField.Len(); k++ {
						parts = append(parts, fmt.Sprintf("%s=%v", key, embeddedField.Index(k)))
					}
				} else {
					parts = append(parts, fmt.Sprintf("%s=%v", key, embeddedField))
				}
			}
			continue
		}

		if field.Kind() == reflect.Pointer {
			field = field.Elem()
		}

		if !field.IsValid() || isZeroValue(field) {
			continue
		}

		key := getJSONFieldName(fieldType)
		if key == "" || key == "-" {
			continue
		}

		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				parts = append(parts, fmt.Sprintf("%s=%v", key, field.Index(j)))
			}
		} else {
			parts = append(parts, fmt.Sprintf("%s=%v", key, field))
		}
	}

	return strings.Join(parts, "&")
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (c *Client) GenerateAuthURL() string {
	base := "https://id.twitch.tv/oauth2/authorize"
	params := url.Values{}
	params.Add("client_id", c.clientID)
	params.Add("response_type", "code")

	if c.redirectURI != nil {
		params.Add("redirect_uri", *c.redirectURI)
	}

	if len(c.scopes) > 0 {
		scopeStrings := ScopesToStrings(c.scopes)
		params.Add("scope", strings.Join(scopeStrings, " "))
	}

	return base + "?" + params.Encode()
}

func (c *Client) BanUser(ctx context.Context, channel, user, reason string) (*APIBanResponse, error) {
	if c.user == nil {
		return &APIBanResponse{Data: []Ban{}}, c.error("local user is null")
	}

	users, err := c.GetUsers(ctx, []string{channel, c.user.Login, user})
	if err != nil {
		return &APIBanResponse{Data: []Ban{}}, err
	}

	var channelUser, modUser, userUser *User
	for _, u := range users.Data {
		switch u.Login {
		case channel:
			channelUser = &u
		case c.user.Login:
			modUser = &u
		case user:
			userUser = &u
		}
	}

	if channelUser == nil || modUser == nil || userUser == nil {
		return &APIBanResponse{Data: []Ban{}}, c.error("failed to fetch required users")
	}

	query := fmt.Sprintf("?broadcaster_id=%s&moderator_id=%s", channelUser.ID, modUser.ID)
	endpoint := "/moderation/bans" + query

	banData := map[string]any{
		"user_id": userUser.ID,
	}

	if reason != "" {
		banData["reason"] = reason
	}

	requestData := map[string]any{"data": banData}
	data, err := c.post(ctx, endpoint, requestData)
	if err != nil {
		return &APIBanResponse{Data: []Ban{}}, err
	}

	var result APIBanResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return &APIBanResponse{Data: []Ban{}}, err
	}

	return &result, nil
}

func (c *Client) ShoutoutUser(ctx context.Context, channel, user string) error {
	if c.user == nil {
		return c.error("local user is null")
	}

	users, err := c.GetUsers(ctx, []string{channel, c.user.Login, user})
	if err != nil {
		return err
	}

	var channelUser, modUser, userUser *User
	for _, u := range users.Data {
		switch u.Login {
		case channel:
			channelUser = &u
		case c.user.Login:
			modUser = &u
		case user:
			userUser = &u
		}
	}

	if channelUser == nil || modUser == nil || userUser == nil {
		return c.error("failed to fetch required users")
	}

	endpoint := "/chat/shoutouts"
	requestData := map[string]string{
		"from_broadcaster_id": channelUser.ID,
		"to_broadcaster_id": userUser.ID,
		"moderator_id": modUser.ID,
	}

	_, err = c.post(ctx, endpoint, requestData)
	return err
}

func (c *Client) GetUserAccess(ctx context.Context, code string) error {
	endpoint := "https://id.twitch.tv/oauth2/token"
	params := url.Values{}
	params.Add("client_id", c.clientID)
	params.Add("client_secret", c.clientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	if c.redirectURI != nil {
		params.Add("redirect_uri", *c.redirectURI)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint + "?" + params.Encode(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result AuthEvent
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.AccessToken != "" {
		c.accessToken = &result.AccessToken
	}

	if result.RefreshToken != "" {
		c.refreshToken = &result.RefreshToken
	}

	c.emit("user_auth", result)
	return nil
}

func simpleGetDecode[T any](c *Client, ctx context.Context, endpoint string, version string) (*T, error) {
	data, err := c.get(ctx, endpoint, version)
	if err != nil {
		return nil, err
	}

	c.log("GET Decode: Endpoint(%s), Data(%s)", endpoint, data)

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func simplePutDecode[T any](c *Client, ctx context.Context, endpoint string) (*T, error) {
	data, err := c.put(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	c.log("PUT Decode: Endpoint(%s), Data(%s)", endpoint, data)

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func simplePostDecode[T any](c *Client, ctx context.Context, endpoint string) (*T, error) {
	data, err := c.post(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	c.log("POST Decode: Endpoint(%s), Data(%s)", endpoint, data)

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetGames(ctx context.Context, games any) (*APIGameResponse, error) {
	query := "?" + parseMixedParam(games, "name", "id")
	endpoint := "/games" + query

	return simpleGetDecode[APIGameResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetTopGames(ctx context.Context, options *BaseOptions) (*APIGameResponse, error) {
	query := ""
	if options != nil {
		query = "?" + parseOptions(options)
	}
	endpoint := "/games/top" + query

	return simpleGetDecode[APIGameResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetUsers(ctx context.Context, ids any) (*APIUserResponse, error) {
	var query string

	switch v := ids.(type) {
	case []string, []int:
		query = "?" + parseMixedParam(v, "login", "id")
	case string:
		key := chooseKey(!isNumber(v), "login", "id")
		query = fmt.Sprintf("?%s=%s", key, v)
	case int:
		query = fmt.Sprintf("?id=%d", v)
	default:
		return nil, c.error("ids must be a string, int, or slice of strings/ints")
	}

	endpoint := "/users" + query

	return simpleGetDecode[APIUserResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetStreams(ctx context.Context, options *GetStreamsOptions) (*APIStreamResponse, error) {
	query := "?"
	endpoint := "/streams"

	if options == nil {
		return simpleGetDecode[APIStreamResponse](c, ctx, endpoint, "helix")
	}

	channel, channels := options.Channel, options.Channels

	if channel != nil {
		key := chooseKey(isNumber(*channel), "user_id", "user_login")
		query += fmt.Sprintf("%s=%s&", key, *channel)
	}

	if channels != nil {
		query += parseMixedParam(channels, "user_login", "user_id")
	}

	query += "&"
	query += parseOptions(options)

	return simpleGetDecode[APIStreamResponse](c, ctx, endpoint + query, "helix")
}

func (c *Client) GetGlobalBadges(ctx context.Context) (*APIBadgesResponse, error) {
	endpoint := "/chat/badges/global"
	return simpleGetDecode[APIBadgesResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetGlobalEmotes(ctx context.Context) (*APIEmotesResponse, error) {
	endpoint := "/chat/emotes/global"
	return simpleGetDecode[APIEmotesResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetVideos(ctx context.Context, options GetVideosOptions) (*APIVideoResponse, error) {
	query := "?" + parseOptions(&options)
	if strings.Contains(query, "/videos?user_id=0x") {
		c.log("FAIL")
	} else {
		c.log("SUCCESS")
	}
	endpoint := "/videos" + query
	return simpleGetDecode[APIVideoResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetClips(ctx context.Context, options any) (*APIClipsResponse, error) {
	query := "?" + parseOptions(&options)
	endpoint := "/clips" + query
	return simpleGetDecode[APIClipsResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetChannelInformation(ctx context.Context, options GetChannelInfoOptions) (*APIChannelInfoResponse, error) {
	query := "?" + parseOptions(&options)
	endpoint := "/channels" + query
	return simpleGetDecode[APIChannelInfoResponse](c, ctx, endpoint, "helix")
}

func (c *Client) SearchChannels(ctx context.Context, options SearchChannelsOptions) (*APIChannelResponse, error) {
	options.Query = url.QueryEscape(options.Query)
	query := "?" + parseOptions(&options)
	endpoint := "/search/channels" + query
	return simpleGetDecode[APIChannelResponse](c, ctx, endpoint, "helix")
}

func (c *Client) SearchCategories(ctx context.Context, options SearchCategoriesOptions) (*APIGameResponse, error) {
	options.Query = url.QueryEscape(options.Query)
	query := "?" + parseOptions(&options)
	endpoint := "/search/categories" + query
	return simpleGetDecode[APIGameResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetExtensionTransactions(ctx context.Context, options GetExtensionTransactionsOptions) (*APIExtensionTransactionResponse, error) {
	query := "?" + parseOptions(&options)
	endpoint := "/extensions/transactions" + query
	return simpleGetDecode[APIExtensionTransactionResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetCheermotes(ctx context.Context, options *GetCheermotesOptions) (*APICheermoteResponse, error) {
	query := ""
	if options != nil {
		query = "?" + parseOptions(options)
	}

	endpoint := "/bits/cheermotes" + query
	return simpleGetDecode[APICheermoteResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetChannelEmotes(ctx context.Context, broadcasterID string) (*APIEmotesResponse, error) {
	query := "?broadcaster_id=" + broadcasterID
	endpoint := "/chat/emotes" + query
	return simpleGetDecode[APIEmotesResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetChannelBadges(ctx context.Context, broadcasterID string) (*APIBadgesResponse, error) {
	query := "?broadcaster_id=" + broadcasterID
	endpoint := "/chat/badges" + query
	return simpleGetDecode[APIBadgesResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetBitsLeaderboard(ctx context.Context, options *GetBitsLeaderboardOptions) (*APIBitsLeaderboardResponse, error) {
	if !c.hasScope(ScopeBitsRead) {
		return nil, c.error("missing scope: bits:read")
	}

	query := ""
	if options != nil {
		query = "?" + parseOptions(options)
	}
	endpoint := "/bits/leaderboard" + query
	return simpleGetDecode[APIBitsLeaderboardResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetSubs(ctx context.Context, options GetSubsOptions) (*APISubResponse, error) {
	if !c.hasScope(ScopeChannelReadSubscriptions) {
		return nil, c.error("missing scope: channel:read:subscriptions")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/subscriptions" + query
	return simpleGetDecode[APISubResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetBannedUsers(ctx context.Context, options GetBannedUsersOptions) (*APIBanResponse, error) {
	if !c.hasScope(ScopeModerationRead) {
		return nil, c.error("missing scope: moderation:read")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/moderation/banned" + query
	return simpleGetDecode[APIBanResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetStreamMarkers(ctx context.Context, options any) (*APIStreamMarkerResponse, error) {
	if !c.hasScope(ScopeUserReadBroadcast) {
		return nil, c.error("missing scope: user:read:broadcast")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/streams/markers" + query
	return simpleGetDecode[APIStreamMarkerResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetUserExtensions(ctx context.Context) (*APIExtensionResponse, error) {
	if !c.hasScope(ScopeUserReadBroadcast) {
		return nil, c.error("missing scope: user:read:broadcast")
	}

	endpoint := "/users/extensions/list"
	return simpleGetDecode[APIExtensionResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetUserActiveExtensions(ctx context.Context, options *GetUserActiveExtensionsOptions) (*APIActiveUserExtensionResponse, error) {
	if !c.hasScope(ScopeUserReadBroadcast) && !c.hasScope(ScopeUserEditBroadcast) {
		return nil, c.error("missing scope: user:read:broadcast or user:edit:broadcast")
	}

	query := ""
	if options != nil {
		query = "?" + parseOptions(options)
	}
	endpoint := "/users/extensions" + query
	return simpleGetDecode[APIActiveUserExtensionResponse](c, ctx, endpoint, "helix")
}

func (c *Client) ModifyChannelInformation(ctx context.Context, options ModifyChannelInformationOptions) error {
	if !c.hasScope(ScopeUserEditBroadcast) {
		return c.error("missing scope: user:edit:broadcast")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/channels" + query

	_, err := c.patch(ctx, endpoint, nil)
	return err
}

func (c *Client) UpdateUser(ctx context.Context, options *UpdateUserOptions) (*APIUserResponse, error) {
	if !c.hasScope(ScopeUserEdit) {
		return nil, c.error("missing scope: user:edit")
	}

	query := ""
	if options != nil && options.Description != nil {
		query = "?" + parseOptions(options)
	}
	endpoint := "/users" + query
	return simplePutDecode[APIUserResponse](c, ctx, endpoint)
}

func (c *Client) CreateClip(ctx context.Context, options CreateClipOptions) (*APICreateClipResponse, error) {
	if !c.hasScope(ScopeClipsEdit) {
		return nil, c.error("missing scope: clips:edit")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/clips" + query
	return simplePostDecode[APICreateClipResponse](c, ctx, endpoint)
}

func (c *Client) GetModerators(ctx context.Context, options GetModeratorsOptions) (*APIModeratorResponse, error) {
	if !c.hasScope(ScopeModerationRead) {
		return nil, c.error("missing scope: moderation:read")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/moderation/moderators" + query
	return simpleGetDecode[APIModeratorResponse](c, ctx, endpoint, "helix")
}

func (c *Client) GetCodeStatus(ctx context.Context, options GetCodeStatusOptions) (*APICodeStatusResponse, error) {
	query := "?" + parseOptions(&options)
	endpoint := "/entitlements/codes" + query
	return simpleGetDecode[APICodeStatusResponse](c, ctx, endpoint, "helix")
}

func (c *Client) StartCommercial(ctx context.Context, options StartCommercialOptions) (*APICommercialResponse, error) {
	if !c.hasScope(ScopeChannelEditCommercial) {
		return nil, c.error("missing scope: channel:edit:commercial")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/channels/commercial" + query
	return simplePostDecode[APICommercialResponse](c, ctx, endpoint)
}

func (c *Client) GetCurrentUser() (*User, error) {
	ctx := context.Background()
	endpoint := "/users"

	data, err := c.get(ctx, endpoint, "helix")
	if err != nil {
		return nil, err
	}

	var result APIUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, c.error("failed to get current user")
	}

	return &result.Data[0], nil
}

func (c *Client) GetStreamKey(ctx context.Context, options GetStreamKeyOptions) (*string, error) {
	if !c.hasScope(ScopeChannelReadStreamKey) {
		return nil, c.error("missing scope: channel:read:stream_key")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/streams/key" + query

	data, err := c.get(ctx, endpoint, "helix")
	if err != nil {
		return nil, err
	}

	var result APIStreamKeyResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, c.error("no stream key found")
	}

	return &result.Data[0].Key, nil
}

func (c *Client) SendChatMessage(ctx context.Context, options SendChatMessageOptions) (*APIMessageResponse, error) {
	if !c.hasScope(ScopeUserBot) {
		return nil, c.error("missing scope: user:bot")
	}

	query := "?" + parseOptions(&options)
	endpoint := "/chat/messages" + query

	data, err := c.post(ctx, endpoint, options)
	if err != nil {
		return nil, err
	}

	var result APIMessageResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetIngestServers(ctx context.Context) (*APIIngestsResponse, error) {
	endpoint := "/ingests"

	data, err := c.get(ctx, endpoint, "ingest")
	if err != nil {
		return nil, err
	}

	var result APIIngestsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
