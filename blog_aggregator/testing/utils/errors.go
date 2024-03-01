package utils

import "log"

func RequestFormationError(e error, caller string) {
	log.Fatalf("Couldn't form request: %v (%s)", e.Error(), caller)
}

func EncodingError(e error, caller string) {
	log.Fatalf("Error encoding json: %v (%s)", e.Error(), caller)
}

func BadResponseError(e error, caller string) {
	log.Fatalf("Failed getting response: %s (%s)", e.Error(), caller)
}

func StatusCodeError(expected, actual int, caller string) {
	log.Fatalf("Expected response body '%v', got=%v (%s)", expected, actual, caller)
}

func BodyResponseError(e error, caller string) {
	log.Fatalf("Couldn't read response body: %s (%s)", e.Error(), caller)
}

func JSONUnmarshalError(e error, caller string) {
	log.Fatalf("Could not unmarshal response body: %s (%s)", e.Error(), caller)
}

func ResponsePayloadError(caller string) {
	log.Fatalf("Incorrect payload type received from %s", caller)
}

func WrongPayloadError(payload interface{}, caller string) {
	log.Fatalf("Unexpected payload. \ngot=%v (%s)", payload, caller)
}
