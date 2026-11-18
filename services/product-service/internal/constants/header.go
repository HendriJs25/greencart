package constants

import "net/textproto"

var (
	CSRFToken = textproto.CanonicalMIMEHeaderKey("X-CSRF-Token")
)
