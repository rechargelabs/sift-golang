//
//

package sift

var (
	// NoContentStatusCodes -
	NoContentStatusCodes = map[int]string{204: "", 304: ""}

	// AvailableMethods List of all methods that Sift Science API accepts
	AvailableMethods = map[string]string{"GET": "get", "POST": "post", "DELETE": "delete"}

	// ErrorCodes A successful API request will respond with an HTTP 200. An invalid API
	// request will respond with an HTTP 400. The response body will be a JSON
	// object describing why the request failed.
	// These are JSON error response codes in case you need them
	ErrorCodes = map[int]string{
		-4:  "Service currently unavailable. Please try again later.",
		-3:  "Server-side timeout processing request. Please try again later.",
		-2:  "Unexpected server-side error",
		-1:  "Unexpected server-side error",
		0:   "Success",
		51:  "Invalid API key",
		52:  "Invalid characters in field name",
		53:  "Invalid characters in field value",
		55:  "Missing required field",
		56:  "Invalid JSON in request",
		57:  "Invalid HTTP body",
		60:  "Rate limited",
		104: "Invalid API version",
		105: "Not a valid reserved field",
	}
)
