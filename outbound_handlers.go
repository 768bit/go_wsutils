package go_wsutils

import (
	"encoding/binary"
	"encoding/json"
	"github.com/768bit/websocket"
	"github.com/google/uuid"
)

func SendStreamCmdRequest(socketConn *websocket.Conn, id uint64, seq uint64, seshKey uuid.UUID, cmd uint16, payload []byte) {

	//send the payload after building the packet...

	var packetHeader []byte = make([]byte, HEADER_SIZE)

	binary.LittleEndian.PutUint64(payload[SESH_ID_HEADER_INDEX_START:SESH_ID_HEADER_INDEX_END], id)

	binary.LittleEndian.PutUint64(payload[SEQ_ID_HEADER_INDEX_START:SEQ_ID_HEADER_INDEX_END], id)

	copy(payload[SESH_KEY_HEADER_INDEX_START:SESH_KEY_HEADER_INDEX_END], socketConn.GetSeshKey())

	binary.LittleEndian.PutUint16(payload[CMD_ID_HEADER_INDEX_START:CMD_ID_HEADER_INDEX_END], cmd)

	binary.LittleEndian.PutUint16(payload[PAYLOAD_LENGTH_HEADER_INDEX_START:PAYLOAD_LENGTH_HEADER_INDEX_END], (uint16)(len(payload)))

	payload = append(packetHeader, payload...)

	socketConn.WriteMessage(websocket.BinaryMessage, payload)

	//payload written we n

}

//send JSON request will send a message to the server for processing - there are several types of JSON message this is the lowest level
func SendJSONRequest(requestID string, conn *websocket.Conn, payload interface{}, req *WSRequest) {

	encMsg, err := json.Marshal(payload)

	if err == nil {

		if sendErr := conn.WriteMessage(websocket.TextMessage, encMsg); err != nil {

			req.Errors = append(req.Errors, sendErr.Error())

			req.Done <- false

			req.Response <- NewWSRequestLocalErrorResponse(requestID, conn.GetSeshKey())

			conn.Close()

		}

	} else {

		req.Errors = append(req.Errors, "Error marshalling payload to JSON")
		req.Errors = append(req.Errors, err.Error())

		req.Done <- false

		req.Response <- NewWSRequestLocalErrorResponse(requestID, conn.GetSeshKey())

	}

}

func SendJSONMessage(conn *websocket.Conn, payload interface{}) error {

	encMsg, err := json.Marshal(payload)

	if err == nil {

		if sendErr := conn.WriteMessage(websocket.TextMessage, encMsg); err != nil {

			return sendErr

		} else {

			//conn.WriteMessage(websocket.BinaryMessage, []byte("hello"))

			return nil

		}

	} else {

		return err

	}

}
