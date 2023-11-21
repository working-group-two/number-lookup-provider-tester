# Number Lookup Provider Tester

This application is designed to make it easier to develop a Number Lookup Provider, and will generate requests for
different phone numbers, and show the responses from the provider.

That is, it allows testing your integration with the Number Lookup API without having to make actual calls to the API.

## About the Number Lookup API

This is a test application for the WG2 Number Lookup API.

This API uses a bidirectional gRPC stream, where the server will send requests to the client for looking up information
about a phone number. The client will then respond with the information it has about the phone number.

This information will be used in the call flow to show the user information about the caller, and to determine if the
call should be blocked or not.

For more information about the API, see the following links:

- Docs: https://v0.docs.wgtwo.com/number-lookup/overview/
- Proto: https://github.com/working-group-two/wgtwoapis/blob/master/wgtwo/lookup/v0/number_lookup.proto

## About the test application

The test application does not inspect any authentication headers.

For testing authentication, you may either connect to our sandbox environment or the real production API.

It does not limit number of in-flight requests.

## Usage

```
Number Lookup Provider Tester

Usage:
  number-lookup [flags]

Flags:
  -h, --help                   help for number-lookup
      --phone-number strings   Phone number to use in requests
      --port int               Port to run the application on
      --print-progress         Print progress
      --print-requests         Print requests
      --print-responses        Print responses
      --rps uint32             Requests per second
```

### Example

```
go run main.go \
  --port 8080 \
  --rps 5 \
  --phone-number 4799990001 \
  --phone-number 4799990002,4799990003 \
  --print-progress \
  --print-requests \
  --print-responses
```
