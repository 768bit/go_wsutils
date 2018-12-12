package go_wsutils

import (
	"errors"
	"time"
)

var FINAL_BYTE_MARKER = []byte{0xFF}
var SESH_ID_HEADER_SIZE_BYTES = 8
var SESH_ID_HEADER_INDEX_START = 0
var SESH_ID_HEADER_INDEX_END = SESH_ID_HEADER_INDEX_START + SESH_ID_HEADER_SIZE_BYTES - 1
var SEQ_ID_HEADER_SIZE_BYTES = 8
var SEQ_ID_HEADER_INDEX_START = SESH_ID_HEADER_INDEX_END + 1
var SEQ_ID_HEADER_INDEX_END = SEQ_ID_HEADER_INDEX_START + SEQ_ID_HEADER_SIZE_BYTES - 1
var SESH_KEY_HEADER_SIZE_BYTES = 16
var SESH_KEY_HEADER_INDEX_START = SEQ_ID_HEADER_INDEX_END + 1
var SESH_KEY_HEADER_INDEX_END = SESH_KEY_HEADER_INDEX_START + SESH_KEY_HEADER_SIZE_BYTES - 1
var CMD_ID_HEADER_SIZE_BYTES = 2
var CMD_ID_HEADER_INDEX_START = SESH_KEY_HEADER_INDEX_END + 1
var CMD_ID_HEADER_INDEX_END = CMD_ID_HEADER_INDEX_START + CMD_ID_HEADER_SIZE_BYTES - 1
var PAYLOAD_LENGTH_HEADER_SIZE_BYTES = 2
var PAYLOAD_LENGTH_HEADER_INDEX_START = CMD_ID_HEADER_INDEX_END + 1
var PAYLOAD_LENGTH_HEADER_INDEX_END = PAYLOAD_LENGTH_HEADER_INDEX_START + PAYLOAD_LENGTH_HEADER_SIZE_BYTES - 1
var HEADER_SIZE = (uint16)(SESH_ID_HEADER_SIZE_BYTES +
	SEQ_ID_HEADER_SIZE_BYTES +
	SESH_KEY_HEADER_SIZE_BYTES +
	CMD_ID_HEADER_SIZE_BYTES +
	PAYLOAD_LENGTH_HEADER_SIZE_BYTES)

var NEGOTIATE_TRANSFER = (uint16)(0x00)
var NEGOTIATE_TRANSFER_ACK = (uint16)(0x04)
var TRANSFER_BEGIN = (uint16)(0x08)
var TRANSFER_COMPLETE = (uint16)(0x10)
var TRANSFER_SEQ_ACK = (uint16)(0xAA)
var NEGOTIATE_TRANSFER_ERROR = (uint16)(0xFA)
var TRANSFER_ERROR = (uint16)(0xFB)
var SESSION_MISSING_ERROR = (uint16)(0xFC)
var SESSION_KEY_MISMATCH = (uint16)(0xFD)
var FATAL_ERROR = (uint16)(0xFF)

func NewWebSocketSessionStartRequestBody(requestID string, jwtTicketID string, userUUID string) *WebSocketRequestBody {

	return &WebSocketRequestBody{
		MessageType: RPCSessionStartMessage,
		ID:          requestID,
		Payload: map[string]interface{}{
			"jwtTicketID": jwtTicketID,
			"userUUID":    userUUID,
		},
	}

}

func NewWebSocketSessionStartResponseBody(requestID string, seshKey string) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: RPCSessionStartMessage,
		ID:          requestID,
		SeshKey:     seshKey,
		StatusCode:  RPCStatusOK,
	}

}

func NewWebSocketSessionStartErrorResponseBody(requestID string, status int, err error) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: RPCSessionStartErrorMessage,
		ID:          requestID,
		StatusCode:  status,
		Errors:      []error{err},
	}

}

func NewWebSocketSessionEndRequestBody(requestID string, seshKey string) *WebSocketRequestBody {

	return &WebSocketRequestBody{
		MessageType: RPCSessionEndMessage,
		ID:          requestID,
		SeshKey:     seshKey,
	}

}

func NewWebSocketSessionEndResponseBody(requestID string, seshKey string) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: RPCSessionEndMessage,
		ID:          requestID,
		StatusCode:  RPCStatusOK,
		SeshKey:     seshKey,
	}

}

func NewWebSocketSessionEndErrorResponseBody(requestID string, seshKey string, status int, err error) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: RPCSessionEndErrorMessage,
		ID:          requestID,
		StatusCode:  status,
		SeshKey:     seshKey,
		Errors:      []error{err},
	}

}

const (
	ServerHelloMessage          = 0x00
	RPCSessionStartMessage      = 0x01
	RPCSessionEndMessage        = 0x04
	RPCMessage                  = 0x20
	HTTPMessage                 = 0x40
	ByteSessionStartMessage     = 0xB0
	ByteSessionEndMessage       = 0xB4
	RPCSessionStartErrorMessage = 0xE0
	RPCSessionEndErrorMessage   = 0xE1
	BasicMessage                = 0xFF
)

const (
	RPCStatusOK               = 0x00C8 //200
	RPCStatusUnauthorised     = 0x0191 //401
	RPCStatusError            = 0x01F4 //500
	RPCStatusLocalError       = 0x0266 //550
	RPCStatusRequestCancelled = 0x029E //670
)

type WebSocketRequestBody struct {
	MessageType int                    `json:"messageType,omitEmpty"`
	Cmd         string                 `json:"cmd,omitEmpty"`
	Method      string                 `json:"method,omitEmpty"`
	Path        string                 `json:"path,omitEmpty"`
	ModuleURI   string                 `json:"moduleURI,omitEmpty"`
	ID          string                 `json:"id,omitEmpty"`
	SeshKey     string                 `json:"seshKey,omitEmpty"`
	Headers     map[string]string      `json:"headers,omitEmpty"`
	Payload     map[string]interface{} `json:"payload,omitEmpty"`
	Options     map[string]interface{} `json:"options,omitEmpty"`
	StatusCode  int                    `json:"statusCode,omitEmpty"`
}

type WebSocketResponseBody struct {
	MessageType int                    `json:"messageType,omitEmpty"`
	Cmd         string                 `json:"cmd,omitEmpty"`
	Method      string                 `json:"method,omitEmpty"`
	Path        string                 `json:"path,omitEmpty"`
	ModuleURI   string                 `json:"moduleURI,omitEmpty"`
	ID          string                 `json:"id,omitEmpty"`
	SeshKey     string                 `json:"seshKey,omitEmpty"`
	Headers     map[string]string      `json:"headers,omitEmpty"`
	Payload     map[string]interface{} `json:"payload,omitEmpty"`
	Options     map[string]interface{} `json:"options,omitEmpty"`
	StatusCode  int                    `json:"statusCode,omitEmpty"`
	Errors      []error                `json:"errors,omitEmpty"`
}

func NewBasicWebSocketResponseBody(statusCode int, requestID string, payload interface{}) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: BasicMessage,
		StatusCode:  statusCode,
		ID:          requestID,
		Payload:     map[string]interface{}{"response": payload},
	}

}

func NewBasicWebSocketHttpResponseBody(statusCode int, requestID string, method string, path string, payload interface{}) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: HTTPMessage,
		StatusCode:  statusCode,
		ID:          requestID,
		Method:      method,
		Path:        path,
		Payload:     map[string]interface{}{"http_response": payload},
	}

}

func NewWebSocketHttpResponseBody(statusCode int, seshKey string, requestID string, method string, path string, payload interface{}) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: HTTPMessage,
		StatusCode:  statusCode,
		SeshKey:     seshKey,
		ID:          requestID,
		Method:      method,
		Path:        path,
		Payload:     map[string]interface{}{"http_response": payload},
	}

}

func NewWebSocketRPCResponseBody(statusCode int, seshKey string, requestID string, cmd string, payload interface{}) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		MessageType: RPCMessage,
		StatusCode:  statusCode,
		SeshKey:     seshKey,
		ID:          requestID,
		Cmd:         method,
		Payload:     payload,
	}

}

func NewWebSocketRPCErrorResponseBody(statusCode int, seshKey string, requestID string, cmd string, payload interface{}, err error) *WebSocketResponseBody {

	if statusCode <= RPCStatusOK {
		statusCode = RPCStatusError
	}

	return &WebSocketResponseBody{
		MessageType: RPCMessage,
		StatusCode:  statusCode,
		SeshKey:     seshKey,
		ID:          requestID,
		Cmd:         method,
		Payload:     payload,
		Errors:      []error{err},
	}

}

type WSRequest struct {
	requestID       string
	requestBody     *WebSocketRequestBody
	httpRequestBody *WebSocketRequestBody
	seshKey         string
	Cancelled       bool
	Timeout         int
	AckTimeout      int
	timeoutTimer    *time.Timer
	ackTimeoutTimer *time.Timer
	Done            chan bool
	Progress        chan *WSRequestProgress
	Response        chan *WebSocketResponseBody
	Errors          []error
}

func NewBasicWSRequest(requestID string, requestBody *WebSocketRequestBody) *WSRequest {

	return &WSRequest{
		requestID:   requestID,
		requestBody: requestBody,
		Cancelled:   false,
		Done:        make(chan bool),
		Progress:    make(chan *WSRequestProgress),
		Response:    make(chan *WebSocketResponseBody),
		Errors:      []error{},
	}

}

func NewWSRequest(requestID string, seshKey string, requestBody *WebSocketRequestBody) *WSRequest {

	return &WSRequest{
		requestID:   requestID,
		requestBody: requestBody,
		seshKey:     seshKey,
		Cancelled:   false,
		Done:        make(chan bool),
		Progress:    make(chan *WSRequestProgress),
		Response:    make(chan *WebSocketResponseBody),
		Errors:      []error{},
	}

}

func NewWSRequestWithTimeout(requestID string, seshKey string, requestBody *WebSocketRequestBody, timeout int) *WSRequest {

	return &WSRequest{
		requestID:    requestID,
		requestBody:  requestBody,
		seshKey:      seshKey,
		Cancelled:    false,
		Done:         make(chan bool),
		Progress:     make(chan *WSRequestProgress),
		Response:     make(chan *WebSocketResponseBody),
		Timeout:      timeout,
		timeoutTimer: time.NewTimer(time.Duration(timeout) * time.Second),
		Errors:       []error{},
	}

}

func NewWSRequestWithAckTimeout(requestID string, seshKey string, requestBody *WebSocketRequestBody, ackTimeout int) *WSRequest {

	return &WSRequest{
		requestID:       requestID,
		requestBody:     requestBody,
		seshKey:         seshKey,
		Cancelled:       false,
		Done:            make(chan bool),
		Progress:        make(chan *WSRequestProgress),
		Response:        make(chan *WebSocketResponseBody),
		AckTimeout:      ackTimeout,
		ackTimeoutTimer: time.NewTimer(time.Duration(ackTimeout) * time.Second),
		Errors:          []error{},
	}

}

func NewWSRequestWithAckTimeoutAndTimeout(requestID string, seshKey string, requestBody *WebSocketRequestBody, ackTimeout int, timeout int) *WSRequest {

	return &WSRequest{
		requestID:       requestID,
		requestBody:     requestBody,
		seshKey:         seshKey,
		Cancelled:       false,
		Done:            make(chan bool),
		Progress:        make(chan *WSRequestProgress),
		Response:        make(chan *WebSocketResponseBody),
		AckTimeout:      ackTimeout,
		ackTimeoutTimer: time.NewTimer(time.Duration(ackTimeout) * time.Second),
		Timeout:         timeout,
		timeoutTimer:    time.NewTimer(time.Duration(timeout) * time.Second),
		Errors:          []error{},
	}

}

func NewWSHttpRequest(requestID string, seshKey string, httpRequestBody *WebSocketRequestBody) *WSRequest {

	return &WSRequest{
		requestID:       requestID,
		httpRequestBody: httpRequestBody,
		seshKey:         seshKey,
		Cancelled:       false,
		Done:            make(chan bool),
		Progress:        make(chan *WSRequestProgress),
		Response:        make(chan *WebSocketResponseBody),
		Errors:          []error{},
	}

}

func NewWSHttpRequestWithTimeout(requestID string, seshKey string, httpRequestBody *WebSocketRequestBody, timeout int) *WSRequest {

	return &WSRequest{
		requestID:       requestID,
		httpRequestBody: httpRequestBody,
		seshKey:         seshKey,
		Cancelled:       false,
		Done:            make(chan bool),
		Progress:        make(chan *WSRequestProgress),
		Response:        make(chan *WebSocketResponseBody),
		Timeout:         timeout,
		timeoutTimer:    time.NewTimer(time.Duration(timeout) * time.Second),
		Errors:          []error{},
	}

}

func NewWSHttpRequestWithAckTimeout(requestID string, seshKey string, httpRequestBody *WebSocketRequestBody, ackTimeout int) *WSRequest {

	return &WSRequest{
		requestID:       requestID,
		httpRequestBody: httpRequestBody,
		seshKey:         seshKey,
		Cancelled:       false,
		Done:            make(chan bool),
		Progress:        make(chan *WSRequestProgress),
		Response:        make(chan *WebSocketResponseBody),
		AckTimeout:      ackTimeout,
		ackTimeoutTimer: time.NewTimer(time.Duration(ackTimeout) * time.Second),
		Errors:          []error{},
	}

}

func NewWSHttpRequestWithAckTimeoutAndTimeout(requestID string, seshKey string, httpRequestBody *WebSocketRequestBody, ackTimeout int, timeout int) *WSRequest {

	return &WSRequest{
		requestID:       requestID,
		httpRequestBody: httpRequestBody,
		seshKey:         seshKey,
		Cancelled:       false,
		Done:            make(chan bool),
		Progress:        make(chan *WSRequestProgress),
		Response:        make(chan *WebSocketResponseBody),
		AckTimeout:      ackTimeout,
		ackTimeoutTimer: time.NewTimer(time.Duration(ackTimeout) * time.Second),
		Timeout:         timeout,
		timeoutTimer:    time.NewTimer(time.Duration(timeout) * time.Second),
		Errors:          []error{},
	}

}

func (wsr *WSRequest) CancelRequest() {

	wsr.Cancelled = true

	//add cancelled error to stack... response will be a payload signifying it

	wsr.Errors = append(wsr.Errors, errors.New("Request was cancelled."))
	close(wsr.Progress)

	wsr.Done <- false
	wsr.Response <- NewWSRequestCancelledResponse(wsr.requestID, wsr.seshKey)
	close(wsr.Done)
	close(wsr.Response)

}

type WSRequestProgress struct {
	Percent    float32
	StatusCode int
	Error      error
}

func NewWSRequestLocalErrorResponse(requestID string, seshKey string) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		ID:         requestID,
		StatusCode: RPCStatusLocalError,
		SeshKey:    seshKey,
	}

}

func NewWSRequestCancelledResponse(requestID string, seshKey string) *WebSocketResponseBody {

	return &WebSocketResponseBody{
		ID:         requestID,
		StatusCode: RPCStatusRequestCancelled,
		SeshKey:    seshKey,
	}

}

type WSClientSession struct {
	SeshKey     string
	JWTTicketID string
}

func NewWSClientSession() *WSClientSession {

	return &WSClientSession{}

}
