package go_wsutils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/768bit/websocket"
	"github.com/google/uuid"
	"strings"
)

func HandleByteStream(socketConn *websocket.Conn, data []byte) error {

	//check the length of the payload it must be a minimum of (8 + 8 + 16 + 2 + 2 + 1) = 37 bytes long...

	dataReadLen := uint16(len(data))

	if dataReadLen < 37 {

		//payload is less than 37 bytes... we will send back an error, the client may have called this incorrectly so we let them know but in effect this packet is completely discarded

		return errors.New("Unable to process byte stream message as it is too small")

	} else {

		//now we need to check what the payload length is marked as... get the 35th and 36th byte...

		var seshID uint64

		err := binary.Read(bytes.NewReader(data[SESH_ID_HEADER_INDEX_START:SESH_ID_HEADER_INDEX_END]), binary.LittleEndian, &seshID)

		var seqID uint64

		err = binary.Read(bytes.NewReader(data[SEQ_ID_HEADER_INDEX_START:SEQ_ID_HEADER_INDEX_END]), binary.LittleEndian, &seqID)

		seshKey, err := uuid.FromBytes(data[SESH_KEY_HEADER_INDEX_START:SESH_KEY_HEADER_INDEX_END])

		seshKeyStr := strings.Replace(seshKey.String(), "-", "", -1)

		if err != nil {

			//there was an error getting the session key for this request...

			return errors.New("Unable to get Session Key from message")

		} else if seshKeyStr != socketConn.GetSeshKey() {

			return errors.New("Session Key Mismatch")

		} else {

			var cmd uint16

			err = binary.Read(bytes.NewReader(data[CMD_ID_HEADER_INDEX_START:CMD_ID_HEADER_INDEX_END]), binary.LittleEndian, &cmd)

			var payloadSize uint16

			err = binary.Read(bytes.NewReader(data[PAYLOAD_LENGTH_HEADER_INDEX_START:PAYLOAD_LENGTH_HEADER_INDEX_END]), binary.LittleEndian, &payloadSize)

			if dataReadLen == 37 && payloadSize == 0 && data[36] == FINAL_BYTE_MARKER[0] && cmd == TRANSFER_COMPLETE {
				//IS THE TERMINATION MARKER FOR A PAYLOAD (IS A NULL PAYLOAD) - WILL CLEAN UP SUPPLIED SESSION
				return nil
			} else if dataReadLen == payloadSize+HEADER_SIZE {
				//PROCESS THE COMMAND...
				return ProcessStreamCommand(socketConn, seshID, seqID, seshKey, seshKeyStr, cmd, payloadSize, data[36:])
			} else {

				//INVALID PAYLOAD, LENGTHS MISMATCHED...

				return errors.New("Payload invalid")

			}

		}

	}

}

func ProcessStreamCommand(socketConn *websocket.Conn, id uint64, seq uint64, seshKey uuid.UUID, seshKeyStr string, cmd uint16, size uint16, payload []byte) error {

	//retrieve the session object from cache.. if it doesnt exist the system needs to be able to generate this...
	// in order for this to happen a session error needs to be responded with...

	//fullSeshID := seshKeyStr + ":" + fmt.Sprintf("%08")

	switch cmd {

	case NEGOTIATE_TRANSFER:
		//client is requesting to negotiate the transfer...
	case NEGOTIATE_TRANSFER_ACK:
		//client is acknowliging that the transfer request has succeeded..
	case TRANSFER_BEGIN:
		//open a transfer using the supplied key, this is the first packet in the sequence and will have a seqID of 0 too
	case TRANSFER_SEQ_ACK:
		//we received a packet in the sequence ok (when transmitting to client)
	case NEGOTIATE_TRANSFER_ERROR:
		//the client is reporting an error attempting to negotiate the transfer
	case TRANSFER_ERROR:
		//the client is reporting that it has encountered and error...

	}

	return errors.New("Not implemented")

}
