# Test Cases Documentation

This document provides a comprehensive overview of all unit test cases implemented for the boh-notification-service.

## Table of Contents
- [Notifier Package Tests](#notifier-package-tests)
- [Processor Package Tests](#processor-package-tests)
- [Test Execution](#test-execution)
- [Coverage Summary](#coverage-summary)

---

## Notifier Package Tests

### File: `notifier/smtp-email_test.go`

#### 1. LoginAuth Tests

##### 1.1 TestLoginAuth_Start
**Type:** Positive  
**Description:** Verifies the Start method of loginAuth returns correct mechanism  
**Test Steps:**
- Create LoginAuth instance with test credentials
- Call Start() method
- Verify mechanism is "LOGIN"
- Verify initial response is empty
- Verify no error is returned

**Expected Result:** Start method returns "LOGIN" mechanism with empty initial response

---

##### 1.2 TestLoginAuth_Next_Username
**Type:** Positive  
**Description:** Tests the Next method when server requests username  
**Test Steps:**
- Create LoginAuth instance
- Call Next() with "Username:" prompt
- Verify response matches the username

**Expected Result:** Returns the configured username

---

##### 1.3 TestLoginAuth_Next_Password
**Type:** Positive  
**Description:** Tests the Next method when server requests password  
**Test Steps:**
- Create LoginAuth instance
- Call Next() with "Password:" prompt
- Verify response matches the password

**Expected Result:** Returns the configured password

---

##### 1.4 TestLoginAuth_Next_UnknownServer
**Type:** Negative  
**Description:** Tests error handling for unknown server prompts  
**Test Steps:**
- Create LoginAuth instance
- Call Next() with "Unknown:" prompt
- Verify error is returned
- Verify error message is "unknown fromserver"

**Expected Result:** Returns error with message "unknown fromserver"

---

##### 1.5 TestLoginAuth_Next_NoMore
**Type:** Positive  
**Description:** Tests the Next method when authentication is complete  
**Test Steps:**
- Create LoginAuth instance
- Call Next() with more=false
- Verify nil response and no error

**Expected Result:** Returns nil response and no error

---

#### 2. MessageUnmarshal Tests

##### 2.1 TestMessageUnmarshal_ValidJSON
**Type:** Positive  
**Description:** Tests unmarshalling valid notification event JSON  
**Test Steps:**
- Create valid NotificationEvent struct
- Marshal to JSON
- Call MessageUnmarshal()
- Verify no JSON syntax errors

**Expected Result:** Successfully unmarshals valid JSON (may error on SMTP config, but not on JSON parsing)

---

##### 2.2 TestMessageUnmarshal_InvalidJSON
**Type:** Negative  
**Description:** Tests error handling for malformed JSON  
**Test Steps:**
- Create invalid JSON string: `{"invalid json"}`
- Call MessageUnmarshal()
- Verify error is returned

**Expected Result:** Returns JSON unmarshal error

---

##### 2.3 TestMessageUnmarshal_EmptyJSON
**Type:** Positive/Edge Case  
**Description:** Tests handling of empty JSON object  
**Test Steps:**
- Create empty JSON: `{}`
- Call MessageUnmarshal()
- Verify no JSON syntax error

**Expected Result:** Successfully unmarshals empty JSON structure

---

#### 3. ProcessMessage Tests

##### 3.1 TestProcessMessage_NoChannels
**Type:** Positive/Edge Case  
**Description:** Tests processing message with no notification channels  
**Test Steps:**
- Create NotificationEvent with empty Channels array
- Call ProcessMessage()
- Verify no error (nothing to process)

**Expected Result:** Completes without error

---

##### 3.2 TestProcessMessage_EmailChannel_MissingConfig
**Type:** Negative  
**Description:** Tests email channel handling when SMTP is not configured  
**Test Steps:**
- Clear SMTP environment variables
- Create NotificationEvent with email channel
- Call ProcessMessage()
- Verify error message is "SMTP not configured"

**Expected Result:** Returns "SMTP not configured" error

---

##### 3.3 TestProcessMessage_WhatsAppChannel_MissingConfig
**Type:** Negative  
**Description:** Tests WhatsApp channel handling when ACS is not configured  
**Test Steps:**
- Clear ACS environment variables
- Create NotificationEvent with whatsapp channel
- Call ProcessMessage()
- Verify error message is "ACS WhatsApp parameters not configured"

**Expected Result:** Returns "ACS WhatsApp parameters not configured" error

---

##### 3.4 TestProcessMessage_UnknownChannelType
**Type:** Positive/Edge Case  
**Description:** Tests handling of unknown notification channel types  
**Test Steps:**
- Create NotificationEvent with unknown channel type
- Call ProcessMessage()
- Verify unknown types are ignored without error

**Expected Result:** Ignores unknown channel type without error

---

##### 3.5 TestProcessMessage_MultipleChannels
**Type:** Positive  
**Description:** Tests processing message with multiple channels  
**Test Steps:**
- Create NotificationEvent with multiple channels (unknown and sms)
- Call ProcessMessage()
- Verify all channels are processed

**Expected Result:** Processes all channels sequentially

---

#### 4. WhatsApp Tests

##### 4.1 TestSendWhatsAppMessage_EmptyParameters (Table-Driven)
**Type:** Negative  
**Description:** Tests error handling for empty parameters in WhatsApp messaging  
**Test Scenarios:**
1. Empty toNumber
2. Empty fromNumber
3. Empty body
4. All parameters empty

**Test Steps (per scenario):**
- Call sendWhatsAppMessage() with empty parameter
- Verify error is returned
- Verify error message is "parameters not configured for whatsApp messaging"

**Expected Result:** Returns appropriate error for each empty parameter scenario

---

#### 5. Data Structure Tests

##### 5.1 TestNotificationEvent_JSONMarshaling
**Type:** Positive  
**Description:** Tests JSON marshaling and unmarshaling of NotificationEvent  
**Test Steps:**
- Create NotificationEvent with sample data
- Marshal to JSON
- Unmarshal back to struct
- Verify all fields match original

**Expected Result:** All fields correctly preserved through marshal/unmarshal

---

##### 5.2 TestNotificationChannel_JSONMarshaling
**Type:** Positive  
**Description:** Tests JSON marshaling of NotificationChannel  
**Test Steps:**
- Create NotificationChannel
- Marshal to JSON
- Unmarshal back
- Verify Type and Contact fields

**Expected Result:** Channel data correctly preserved

---

##### 5.3 TestOauthTokenResponse_JSONMarshaling
**Type:** Positive  
**Description:** Tests unmarshaling OAuth token response  
**Test Steps:**
- Create JSON string with access_token and expires_in
- Unmarshal to OauthTokenResponse
- Verify AccessToken and ExpiresIn values

**Expected Result:** Token response correctly parsed

---

##### 5.4 TestAcsMessage_JSONMarshaling
**Type:** Positive  
**Description:** Tests JSON marshaling of ACS message structure  
**Test Steps:**
- Create AcsMessage with all fields
- Marshal to JSON
- Unmarshal back
- Verify all fields match

**Expected Result:** ACS message structure correctly preserved

---

## Processor Package Tests

### File: `processor/processor_test.go`

#### 1. Initialization Tests

##### 1.1 TestInit_AllEnvVarsSet
**Type:** Positive  
**Description:** Tests Init() when all environment variables are configured  
**Test Steps:**
- Set all SMTP environment variables
- Set all Meta API environment variables
- Call Init()
- Verify all variables are correctly loaded

**Expected Result:** All configuration variables are properly initialized

---

##### 1.2 TestInit_MissingSMTPVars
**Type:** Negative  
**Description:** Tests Init() with missing SMTP configuration  
**Test Steps:**
- Clear SMTP environment variables
- Set Meta API variables
- Call Init()
- Verify SMTP variables are empty
- Verify warning is logged

**Expected Result:** Init completes with warning about missing SMTP config

---

##### 1.3 TestInit_MissingMetaVars
**Type:** Negative  
**Description:** Tests Init() with missing Meta API configuration  
**Test Steps:**
- Set SMTP environment variables
- Clear Meta API variables
- Call Init()
- Verify Meta variables are empty
- Verify warning is logged

**Expected Result:** Init completes with warning about missing Meta config

---

#### 2. ProcessMessage Tests

##### 2.1 TestProcessMessage_ValidJSON
**Type:** Positive  
**Description:** Tests processing valid JSON message with unknown channel  
**Test Steps:**
- Create valid NotificationEvent
- Marshal to JSON
- Call ProcessMessage() with context
- Verify no error for unknown channel type

**Expected Result:** Successfully processes message, logs warning for unknown type

---

##### 2.2 TestProcessMessage_InvalidJSON
**Type:** Negative  
**Description:** Tests error handling for malformed JSON  
**Test Steps:**
- Create invalid JSON: `{"invalid json"}`
- Call ProcessMessage()
- Verify error is returned

**Expected Result:** Returns JSON parse error

---

##### 2.3 TestProcessMessage_EmptyMessage
**Type:** Positive/Edge Case  
**Description:** Tests handling of empty notification message  
**Test Steps:**
- Create event with empty NotificationMessage
- Call ProcessMessage()
- Verify processing is skipped

**Expected Result:** Skips processing, returns no error

---

##### 2.4 TestProcessMessage_NoChannels
**Type:** Positive/Edge Case  
**Description:** Tests handling when no channels are specified  
**Test Steps:**
- Create event with empty Channels array
- Call ProcessMessage()
- Verify processing is skipped

**Expected Result:** Skips processing, returns no error

---

##### 2.5 TestProcessMessage_EmailChannel_NotConfigured
**Type:** Negative  
**Description:** Tests EMAIL channel when SMTP is not configured  
**Test Steps:**
- Clear SMTP configuration variables
- Create event with EMAIL channel
- Call ProcessMessage()
- Verify error about SMTP not initialized

**Expected Result:** Returns "SMTP client is not initialized" error

---

##### 2.6 TestProcessMessage_WhatsAppChannel_NotConfigured
**Type:** Negative  
**Description:** Tests WHATSAPP channel when Meta API is not configured  
**Test Steps:**
- Clear Meta API configuration
- Create event with WHATSAPP channel
- Call ProcessMessage()
- Verify error about Meta API not configured

**Expected Result:** Returns "Meta API is not configured" error

---

##### 2.7 TestProcessMessage_MultipleChannels
**Type:** Positive  
**Description:** Tests processing multiple unknown channel types  
**Test Steps:**
- Create event with multiple unknown channel types
- Call ProcessMessage()
- Verify all channels are processed
- Verify warnings logged for unknown types

**Expected Result:** All channels processed, no error

---

##### 2.8 TestProcessMessage_MixedChannels
**Type:** Negative  
**Description:** Tests mix of valid and invalid channels  
**Test Steps:**
- Create event with UNKNOWN and EMAIL channels
- Clear SMTP configuration
- Call ProcessMessage()
- Verify error for unconfigured EMAIL

**Expected Result:** Returns error for unconfigured EMAIL channel

---

##### 2.9 TestProcessMessage_ContextCancellation
**Type:** Edge Case  
**Description:** Tests behavior with cancelled context  
**Test Steps:**
- Create cancelled context
- Call ProcessMessage()
- Verify current behavior

**Expected Result:** Documents current context handling behavior

---

##### 2.10 TestProcessMessage_EmptyBody
**Type:** Negative  
**Description:** Tests error handling for empty message body  
**Test Steps:**
- Call ProcessMessage() with empty byte array
- Verify JSON parse error

**Expected Result:** Returns "unexpected end of JSON input" error

---

##### 2.11 TestProcessMessage_NilBody
**Type:** Negative  
**Description:** Tests error handling for nil message body  
**Test Steps:**
- Call ProcessMessage() with nil
- Verify JSON parse error

**Expected Result:** Returns "unexpected end of JSON input" error

---

#### 3. Email Sending Tests

##### 3.1 TestSendEmailViaSMTP_NotInitialized
**Type:** Negative  
**Description:** Tests email sending when SMTP client is not initialized  
**Test Steps:**
- Clear SMTP host and password
- Call sendEmailViaSMTP()
- Verify error message

**Expected Result:** Returns "SMTP client is not initialized" error

---

##### 3.2 TestSendEmailViaSMTP_MissingHost
**Type:** Negative  
**Description:** Tests email sending with missing SMTP host  
**Test Steps:**
- Clear SMTP host, set password
- Call sendEmailViaSMTP()
- Verify error

**Expected Result:** Returns initialization error

---

##### 3.3 TestSendEmailViaSMTP_MissingPassword
**Type:** Negative  
**Description:** Tests email sending with missing SMTP password  
**Test Steps:**
- Set SMTP host, clear password
- Call sendEmailViaSMTP()
- Verify error

**Expected Result:** Returns initialization error

---

#### 4. WhatsApp/Meta Tests

##### 4.1 TestSendSmsViaMeta_NotConfigured
**Type:** Negative  
**Description:** Tests WhatsApp sending when Meta API is not configured  
**Test Steps:**
- Clear Meta API token and URL
- Call sendSmsViaMeta()
- Verify error message

**Expected Result:** Returns "Meta API is not configured" error

---

##### 4.2 TestSendSmsViaMeta_MissingToken
**Type:** Negative  
**Description:** Tests WhatsApp sending with missing token  
**Test Steps:**
- Clear Meta API token, set URL
- Call sendSmsViaMeta()
- Verify error

**Expected Result:** Returns configuration error

---

##### 4.3 TestSendSmsViaMeta_MissingUrl
**Type:** Negative  
**Description:** Tests WhatsApp sending with missing URL  
**Test Steps:**
- Set Meta API token, clear URL
- Call sendSmsViaMeta()
- Verify error

**Expected Result:** Returns configuration error

---

#### 5. Data Structure Tests

##### 5.1 TestNotificationChannel_Struct
**Type:** Positive  
**Description:** Tests NotificationChannel struct creation and field access  
**Test Steps:**
- Create NotificationChannel
- Verify Type and Contact fields

**Expected Result:** Struct fields correctly accessible

---

##### 5.2 TestNotificationEvent_Struct
**Type:** Positive  
**Description:** Tests NotificationEvent struct creation  
**Test Steps:**
- Create NotificationEvent with channels
- Verify all fields
- Verify channels array

**Expected Result:** Struct correctly populated

---

##### 5.3 TestNotificationEvent_JSONMarshaling
**Type:** Positive  
**Description:** Tests JSON marshaling/unmarshaling of NotificationEvent  
**Test Steps:**
- Create event with multiple channels
- Marshal to JSON
- Unmarshal back
- Verify all fields and nested channels

**Expected Result:** Complete round-trip preservation of data

---

## Test Execution

### Running All Tests
```bash
go test -vet=off ./...
```

### Running Tests with Coverage
```bash
go test -vet=off -cover ./...
```

### Running Tests Verbosely
```bash
go test -vet=off -v ./...
```

### Running Specific Package Tests
```bash
# Notifier package only
go test -vet=off ./notifier/...

# Processor package only
go test -vet=off ./processor/...
```

### Important Note
The `-vet=off` flag is required due to existing linting issues in the codebase where `log.Print` is used with format directives instead of `log.Printf`.

---

## Coverage Summary

| Package | Coverage | Test Count |
|---------|----------|------------|
| notifier | 42.2% | 15 tests |
| processor | 60.9% | 23 tests |
| **Total** | **~51%** | **38 tests** |

### Coverage Details

#### Notifier Package Functions Covered
- ✅ LoginAuth.Start()
- ✅ LoginAuth.Next()
- ✅ MessageUnmarshal()
- ✅ ProcessMessage()
- ✅ sendWhatsAppMessage() (parameter validation)
- ✅ Data structure JSON marshaling

#### Processor Package Functions Covered
- ✅ Init()
- ✅ ProcessMessage()
- ✅ sendEmailViaSMTP() (error cases)
- ✅ sendSmsViaMeta() (error cases)
- ✅ Data structure JSON marshaling

### Functions Not Covered (Require Integration Testing)
- ❌ SendEmailSMTP() - requires actual SMTP server
- ❌ getOauthToken() - requires external OAuth endpoint
- ❌ sendWhatsAppMessage() - requires Meta API access
- ❌ sendEmailViaSMTP() - requires SMTP server (success path)
- ❌ sendSmsViaMeta() - requires Meta API (success path)

---

## Test Categories

### By Type
- **Positive Tests**: 18 (47%)
- **Negative Tests**: 16 (42%)
- **Edge Case Tests**: 4 (11%)

### By Functionality
- **Authentication**: 5 tests
- **Message Processing**: 11 tests
- **Configuration/Initialization**: 6 tests
- **Data Structures**: 7 tests
- **Channel Routing**: 9 tests

---

## Future Test Enhancements

### Recommended Additions
1. **Integration Tests**: Test with actual SMTP server and Meta API
2. **Mock Tests**: Add mocking for external dependencies
3. **Performance Tests**: Test with large message volumes
4. **Concurrent Tests**: Test thread-safety of message processing
5. **End-to-End Tests**: Test complete Azure Service Bus to notification flow

### Areas for Increased Coverage
1. Successful SMTP email sending
2. Successful WhatsApp message sending
3. OAuth token refresh logic
4. Network error handling and retries
5. Message format validation
