package dto

// VanguardResponse is the common response envelope returned by Vanguard APIs.
type VanguardResponse struct {
	Data        any            `json:"data"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Error       *VanguardError `json:"error"`
}

// VanguardError describes an error returned by a Vanguard API.
type VanguardError struct {
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
}
