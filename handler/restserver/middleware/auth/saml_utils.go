package auth

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/crewjam/saml"
// 	"github.com/crewjam/saml/samlsp"
// )

// const (
// 	defaultSessionMaxAge  = time.Hour
// 	claimNameSessionIndex = "SessionIndex"
// )

// type JWTSessionCodec struct {
// 	Audience string
// 	Issuer   string
// 	MaxAge   time.Duration
// }

// func DefaultSessionCodec(opts samlsp.Options) JWTSessionCodec {
// 	// for backwards compatibility, support CookieMaxAge
// 	maxAge := defaultSessionMaxAge
// 	if opts.CookieMaxAge > 0 {
// 		maxAge = opts.CookieMaxAge
// 	}
// 	return JWTSessionCodec{
// 		Audience: opts.URL.String(),
// 		Issuer:   opts.URL.String(),
// 		MaxAge:   maxAge,
// 	}
// }

// func (c JWTSessionCodec) New(assertion *saml.Assertion) (samlsp.Session, error) {
// 	now := saml.TimeNow()
// 	claims := samlsp.JWTSessionClaims{}
// 	claims.SAMLSession = true
// 	claims.Audience = c.Audience
// 	claims.Issuer = c.Issuer
// 	claims.IssuedAt = now.Unix()
// 	claims.ExpiresAt = now.Add(c.MaxAge).Unix()
// 	claims.NotBefore = now.Unix()

// 	if sub := assertion.Subject; sub != nil {
// 		if nameID := sub.NameID; nameID != nil {
// 			claims.Subject = nameID.Value
// 		}
// 	}

// 	claims.Attributes = map[string][]string{}

// 	for _, attributeStatement := range assertion.AttributeStatements {
// 		for _, attr := range attributeStatement.Attributes {
// 			claimName := attr.FriendlyName
// 			if claimName == "" {
// 				claimName = attr.Name
// 			}
// 			for _, value := range attr.Values {
// 				claims.Attributes[claimName] = append(claims.Attributes[claimName], value.Value)
// 			}
// 		}
// 	}

// 	// add SessionIndex to claims Attributes
// 	for _, authnStatement := range assertion.AuthnStatements {
// 		claims.Attributes[claimNameSessionIndex] = append(claims.Attributes[claimNameSessionIndex],
// 			authnStatement.SessionIndex)
// 	}

// 	return claims, nil
// }

// func (c JWTSessionCodec) Encode(s samlsp.Session) (string, error) {
// 	buf, err := json.Marshal(s)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(buf), nil
// }

// func (c JWTSessionCodec) Decode(signed string) (samlsp.Session, error) {
// 	claims := samlsp.JWTSessionClaims{}
// 	if err := json.Unmarshal([]byte(signed), &claims); err != nil {
// 		return nil, err
// 	}
// 	if !claims.VerifyAudience(c.Audience, true) {
// 		return nil, fmt.Errorf("expected audience %q, got %q", c.Audience, claims.Audience)
// 	}
// 	if !claims.VerifyIssuer(c.Issuer, true) {
// 		return nil, fmt.Errorf("expected issuer %q, got %q", c.Issuer, claims.Issuer)
// 	}
// 	if claims.SAMLSession != true {
// 		return nil, errors.New("expected saml-session")
// 	}
// 	return claims, nil
// }
