package ktntwitchgo

import (
	"os"
	"reflect"
	"testing"
)

func TestDecodeDataInstanceBytes(t *testing.T) {
	test := formTest(t, "decode data instance from bytes")

	// Test valid JSON
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	jsonData := []byte(`{"name":"test","value":123}`)
	result, err := DecodeDataInstanceBytes[TestStruct](jsonData)

	if err != nil {
		t.Errorf("Failed to decode: %v", err)
	}

	test.expect("test", result.Name)
	test.expect(123, result.Value)

	// Test invalid JSON
	invalidData := []byte(`{invalid json}`)
	_, err = DecodeDataInstanceBytes[TestStruct](invalidData)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}
}

func TestEncodeDataInstanceBytes(t *testing.T) {
	test := formTest(t, "encode data instance to bytes")

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	obj := &TestStruct{
		Name:  "test",
		Value: 123,
	}

	data, err := EncodeDataInstanceBytes(obj)
	if err != nil {
		t.Errorf("Failed to encode: %v", err)
	}

	expected := `{"name":"test","value":123}`
	test.expect(expected, string(data))
}

func TestDecodeDataInstanceString(t *testing.T) {
	test := formTest(t, "decode data instance from string")

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	jsonString := `{"name":"test","value":456}`
	result, err := DecodeDataInstanceString[TestStruct](jsonString)

	if err != nil {
		t.Errorf("Failed to decode: %v", err)
	}

	test.expect("test", result.Name)
	test.expect(456, result.Value)
}

func TestEncodeDataInstanceString(t *testing.T) {
	test := formTest(t, "encode data instance to string")

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	obj := &TestStruct{
		Name:  "test",
		Value: 789,
	}

	data, err := EncodeDataInstanceString(obj)
	if err != nil {
		t.Errorf("Failed to encode: %v", err)
	}

	expected := `{"name":"test","value":789}`
	test.expect(expected, data)
}

func TestLoadLocalCache(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	originalUserFile := userFile
	userFile = tempDir + "/apiUser.json"
	defer func() { userFile = originalUserFile }()

	// Test with valid cache file
	testCache := LocalCache{
		AccessToken:  "test_access",
		RefreshToken: "test_refresh",
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}

	data, _ := EncodeDataInstanceBytes(&testCache)
	err := os.WriteFile(userFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cache, err := LoadLocalCache()
	if err != nil {
		t.Errorf("Failed to load cache: %v", err)
	}

	test := formTest(t, "load local cache")
	test.expect("test_access", cache.AccessToken)
	test.expect("test_refresh", cache.RefreshToken)
	test.expect("test_client_id", cache.ClientID)
	test.expect("test_client_secret", cache.ClientSecret)

	// Test with non-existent file
	userFile = tempDir + "/nonexistent.json"
	_, err = LoadLocalCache()
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestGetLocalAccessToken(t *testing.T) {
	tempDir := t.TempDir()
	originalUserFile := userFile
	userFile = tempDir + "/apiUser.json"
	defer func() { userFile = originalUserFile }()

	testCache := LocalCache{
		AccessToken: "my_access_token",
	}

	data, _ := EncodeDataInstanceBytes(&testCache)
	os.WriteFile(userFile, data, 0644)

	token, err := GetLocalAccessToken()
	if err != nil {
		t.Errorf("Failed to get access token: %v", err)
	}

	test := formTest(t, "get local access token")
	test.expect("my_access_token", token)
}

func TestGetLocalRefreshToken(t *testing.T) {
	tempDir := t.TempDir()
	originalUserFile := userFile
	userFile = tempDir + "/apiUser.json"
	defer func() { userFile = originalUserFile }()

	testCache := LocalCache{
		RefreshToken: "my_refresh_token",
	}

	data, _ := EncodeDataInstanceBytes(&testCache)
	os.WriteFile(userFile, data, 0644)

	token, err := GetLocalRefreshToken()
	if err != nil {
		t.Errorf("Failed to get refresh token: %v", err)
	}

	test := formTest(t, "get local refresh token")
	test.expect("my_refresh_token", token)
}

func TestGetLocalClientID(t *testing.T) {
	tempDir := t.TempDir()
	originalUserFile := userFile
	userFile = tempDir + "/apiUser.json"
	defer func() { userFile = originalUserFile }()

	testCache := LocalCache{
		ClientID: "my_client_id",
	}

	data, _ := EncodeDataInstanceBytes(&testCache)
	os.WriteFile(userFile, data, 0644)

	clientID, err := GetLocalClientID()
	if err != nil {
		t.Errorf("Failed to get client ID: %v", err)
	}

	test := formTest(t, "get local client ID")
	test.expect("my_client_id", clientID)
}

func TestGetLocalClientSecret(t *testing.T) {
	tempDir := t.TempDir()
	originalUserFile := userFile
	userFile = tempDir + "/apiUser.json"
	defer func() { userFile = originalUserFile }()

	testCache := LocalCache{
		ClientSecret: "my_client_secret",
	}

	data, _ := EncodeDataInstanceBytes(&testCache)
	os.WriteFile(userFile, data, 0644)

	secret, err := GetLocalClientSecret()
	if err != nil {
		t.Errorf("Failed to get client secret: %v", err)
	}

	test := formTest(t, "get local client secret")
	test.expect("my_client_secret", secret)
}

func TestParseOptionsUtil(t *testing.T) {
	test := formTest(t, "parse options from struct (util)")

	type TestOptions struct {
		Name    string   `json:"name"`
		Count   int      `json:"count"`
		Items   []string `json:"items"`
		Enabled bool     `json:"enabled"`
	}

	// Test basic struct
	opts := TestOptions{
		Name:    "test",
		Count:   5,
		Items:   []string{"a", "b", "c"},
		Enabled: true,
	}

	result := ParseOptions(opts)
	if result == "" {
		t.Errorf("Expected non-empty result")
	}

	// The result should contain all fields (order may vary)
	if !containsHelperUtil(result, "Name=test") {
		t.Errorf("Expected result to contain 'Name=test', got: %s", result)
	}
	if !containsHelperUtil(result, "Count=5") {
		t.Errorf("Expected result to contain 'Count=5', got: %s", result)
	}
	if !containsHelperUtil(result, "Enabled=true") {
		t.Errorf("Expected result to contain 'Enabled=true', got: %s", result)
	}

	// Test with pointer
	result = ParseOptions(&opts)
	if result == "" {
		t.Errorf("Expected non-empty result for pointer")
	}

	// Test with nil pointer
	var nilPtr *TestOptions
	result = ParseOptions(nilPtr)
	test.expect("", result)

	// Test with non-struct
	result = ParseOptions(123)
	test.expect("", result)
}

func TestParseOptionsSlices(t *testing.T) {
	type TestOptions struct {
		IDs []string `json:"ids"`
	}

	opts := TestOptions{
		IDs: []string{"123", "456", "789"},
	}

	result := ParseOptions(opts)

	// Should contain all IDs
	if !containsHelperUtil(result, "IDs=123") {
		t.Errorf("Expected result to contain 'IDs=123', got: %s", result)
	}
	if !containsHelperUtil(result, "IDs=456") {
		t.Errorf("Expected result to contain 'IDs=456', got: %s", result)
	}
	if !containsHelperUtil(result, "IDs=789") {
		t.Errorf("Expected result to contain 'IDs=789', got: %s", result)
	}
}

func TestParseOptionsEmpty(t *testing.T) {
	type TestOptions struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	// Empty struct
	opts := TestOptions{}
	result := ParseOptions(opts)

	// Since Count is 0, it should be included
	if !containsHelperUtil(result, "Count=0") {
		t.Errorf("Expected result to contain 'Count=0', got: %s", result)
	}
}

func TestParseOptionsNilFields(t *testing.T) {
	type TestOptions struct {
		Name  *string  `json:"name"`
		Items []string `json:"items"`
	}

	// Struct with nil pointer
	opts := TestOptions{
		Name:  nil,
		Items: nil,
	}

	result := ParseOptions(opts)

	// Nil fields should be skipped
	if containsHelperUtil(result, "name") {
		t.Errorf("Expected nil field 'name' to be skipped, got: %s", result)
	}
}

func TestParseOptionsReflectValue(t *testing.T) {
	type TestOptions struct {
		Value int `json:"value"`
	}

	opts := TestOptions{Value: 42}

	result := ParseOptions(opts)
	if !containsHelperUtil(result, "Value=42") {
		t.Errorf("Expected result to contain 'Value=42', got: %s", result)
	}
}

func TestGetJSONFieldName(t *testing.T) {
	test := formTest(t, "get JSON field name")

	type TestStruct struct {
		Name       string `json:"name"`
		UserID     string `json:"user_id,omitempty"`
		IgnoredTag string `json:"-"`
		NoTag      string
	}

	v := reflect.TypeOf(TestStruct{})

	nameField, _ := v.FieldByName("Name")
	test.expect("name", getJSONFieldName(nameField))

	userIDField, _ := v.FieldByName("UserID")
	test.expect("user_id", getJSONFieldName(userIDField))

	ignoredField, _ := v.FieldByName("IgnoredTag")
	test.expect("-", getJSONFieldName(ignoredField))

	noTagField, _ := v.FieldByName("NoTag")
	test.expect("", getJSONFieldName(noTagField))
}

func TestIsZeroValue(t *testing.T) {
	test := formTest(t, "check if value is zero")

	// Test slice
	var nilSlice []string
	test.expect(true, isZeroValue(reflect.ValueOf(nilSlice)))
	test.expect(true, isZeroValue(reflect.ValueOf([]string{})))
	test.expect(false, isZeroValue(reflect.ValueOf([]string{"a"})))

	// Test map
	var nilMap map[string]int
	test.expect(true, isZeroValue(reflect.ValueOf(nilMap)))
	test.expect(true, isZeroValue(reflect.ValueOf(map[string]int{})))
	test.expect(false, isZeroValue(reflect.ValueOf(map[string]int{"a": 1})))

	// Test pointer
	var nilPtr *string
	test.expect(true, isZeroValue(reflect.ValueOf(nilPtr)))
	str := "test"
	test.expect(false, isZeroValue(reflect.ValueOf(&str)))

	// Test zero values
	test.expect(true, isZeroValue(reflect.ValueOf(0)))
	test.expect(true, isZeroValue(reflect.ValueOf("")))
	test.expect(true, isZeroValue(reflect.ValueOf(false)))

	// Test non-zero values
	test.expect(false, isZeroValue(reflect.ValueOf(42)))
	test.expect(false, isZeroValue(reflect.ValueOf("test")))
	test.expect(false, isZeroValue(reflect.ValueOf(true)))
}

// Helper function to check if a string contains a substring
func containsHelperUtil(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
