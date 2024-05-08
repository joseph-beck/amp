package status

const (
	// Informational Response
	Continue              = 100
	SwitchingProtocols    = 101
	Processing            = 102
	EarlyHints            = 103
	ResponseIsStale       = 110
	RevalidationFailed    = 111
	DisconnectedOperation = 112
	HeuristicExpiration   = 113
	// Success
	OK                          = 200
	Created                     = 201
	Accepted                    = 202
	NonAuthoritativeInformation = 203
	NoContent                   = 204
	ResetContent                = 205
	PartialContent              = 206
	MultiStatus                 = 207
	AlreadyReported             = 208
	ThisIsFine                  = 218 // unofficial
	ImUsed                      = 226
	// Redirection
	MultipleChoices   = 300
	MovedPermanently  = 301
	Found             = 302
	SeeOther          = 303
	NotModified       = 304
	UseProxy          = 305
	SwitchProxy       = 306
	TemporaryRedirect = 307
	PermanentRedirect = 308
	// Client Errors
	BadRequest                       = 400
	Unauthorized                     = 401
	PaymentRequired                  = 402
	Forbidden                        = 403
	NotFound                         = 404
	MethodNotAllowed                 = 405
	NotAcceptable                    = 406
	ProxyAuthenticationRequired      = 407
	RequestTimeout                   = 408
	Conflict                         = 409
	Gone                             = 410
	LengthRequired                   = 411
	PreconditionFailed               = 412
	PayloadTooLarge                  = 413
	URITooLong                       = 414
	UnsupportedMediaType             = 415
	RangeNotSatisfiable              = 416
	ExpectationFailed                = 417
	ImATeapot                        = 418
	PageExpired                      = 419 // Unofficial
	MethodFailure                    = 420 // Unofficial
	MisdirectedRequest               = 421
	UnprocessableContent             = 422
	Locked                           = 423
	FailedDependency                 = 424
	TooEarly                         = 425
	UpgradeRequired                  = 426
	PreconditionRequired             = 428
	TooManyRequests                  = 429
	RequestHeaderFieldsTooLarge      = 431
	NoResponse                       = 444
	BlockedByWindowsParentalControls = 450 // Unofficial
	UnavailableForLegalReasons       = 451
	InvalidToken                     = 498 // Unofficial
	TokenRequired                    = 499 // Unofficial
	// Server Errors
	InternalServerError           = 500
	NotImplemented                = 501
	BadGateway                    = 502
	ServiceUnavailable            = 503
	GatewayTimeout                = 504
	HTTPVersionNotSupported       = 505
	VariantAlsoNegotiates         = 506
	InsufficientStorage           = 507
	LoopDetected                  = 508
	BandwidthLimit                = 509 // Unofficial
	NotExtended                   = 510
	NetworkAuthenticationRequired = 511
	WebServerReturnedUnknownError = 520
	WebServerIsDown               = 521 // Cloudflare
	ConnectionTimedOut            = 522 // Cloudflare
	OriginIsUnreachable           = 523 // Cloudflare
	TimeoutOccurred               = 524 // Cloudflare
	RailgunError                  = 527 // Cloudflare
	SiteOverloaded                = 529 // Unofficial
	SiteFrozen                    = 530 // Unofficial
	NetworkReadTimeout            = 598 // Unofficial
	NetworkConnectTimeout         = 599 // Unofficial
)
