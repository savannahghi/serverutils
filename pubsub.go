package base

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

// pubsub constants
const (
	PubSubHandlerPath = "/pubsub"
	Aud               = "bewell.co.ke"

	authHeaderName    = "Authorization"
	googleIss         = "accounts.google.com"
	googleAccountsIss = "https://accounts.google.com"
	topicKey          = "topicID"
)

// PubSubMessage is a pub-sub message payload.
//
// See https://cloud.google.com/pubsub/docs/push for more context.
//
// The message that is POSTed looks like the example below:
//
// {
//     "message": {
//         "attributes": {
//             "key": "value"
//         },
//         "data": "SGVsbG8gQ2xvdWQgUHViL1N1YiEgSGVyZSBpcyBteSBtZXNzYWdlIQ==",
//         "messageId": "136969346945"
//     },
//    "subscription": "projects/myproject/subscriptions/mysubscription"
// }
type PubSubMessage struct {
	MessageID  string            `json:"messageId"`
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

// PubSubPayload is the payload of a Pub/Sub event.
type PubSubPayload struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func VerifyPubSubJWTAndDecodePayload(
	w http.ResponseWriter,
	r *http.Request,
) (*PubSubPayload, error) {
	authHeader := r.Header.Get(authHeaderName)
	if authHeader == "" || len(strings.Split(authHeader, " ")) != 2 {
		return nil, fmt.Errorf("missing Authorization Header")
	}

	token := strings.Split(authHeader, " ")[1]
	payload, err := idtoken.Validate(
		r.Context(), token, Aud)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if payload.Issuer != googleIss && payload.Issuer != googleAccountsIss {
		return nil, fmt.Errorf("%s is not a valid issuer", payload.Issuer)
	}

	// decode the message
	var m PubSubPayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read request body: %w", err)
	}

	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf(
			"can't unmarshal payload body into JSON struct: %w", err)
	}

	// return the decoded message if there is no error
	return &m, nil
}

// GetPubSubTopic retrieves a pubsub topic from a pubsub payload.
//
// It follows a convention where the topic is sent as an attribute under the
// `topicID` key.
func GetPubSubTopic(m *PubSubPayload) (string, error) {
	if m == nil {
		return "", fmt.Errorf("nil pub sub payload")
	}
	attrs := m.Message.Attributes
	topicID, prs := m.Message.Attributes[topicKey]
	if !prs {
		return "", fmt.Errorf(
			"no `%s` key in message attributes %#v", topicKey, attrs)
	}
	return topicID, nil
}
