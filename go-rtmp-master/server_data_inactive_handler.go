//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package rtmp

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/yutopp/go-rtmp/internal"
	"github.com/yutopp/go-rtmp/message"
)

var _ stateHandler = (*serverDataInactiveHandler)(nil)

// serverDataInactiveHandler Handle data messages from a non operated client at server side.
//
//	transitions:
//	  | "publish" -> serverDataPublishHandler
//	  | "play"	   -> serverDataPlayHandler
//	  | _         -> self
type serverDataInactiveHandler struct {
	sh *streamHandler
}

func (h *serverDataInactiveHandler) onMessage(
	chunkStreamID int,
	timestamp uint32,
	msg message.Message,
) error {
	return internal.ErrPassThroughMsg
}

func (h *serverDataInactiveHandler) onData(
	chunkStreamID int,
	timestamp uint32,
	dataMsg *message.DataMessage,
	body interface{},
) error {
	return internal.ErrPassThroughMsg
}

func (h *serverDataInactiveHandler) onCommand(
	chunkStreamID int,
	timestamp uint32,
	cmdMsg *message.CommandMessage,
	body interface{},
) error {
	l := h.sh.Logger()

	switch cmd := body.(type) {
	case *message.NetStreamPublish:
		l.Infof("Publisher is comming: %#v", cmd)

		streamCtx := &StreamContext{
			StreamID: h.sh.stream.streamID,
		}
		if err := h.sh.stream.userHandler().OnPublish(streamCtx, timestamp, cmd); err != nil {
			// TODO: Support message.NetStreamOnStatusCodePublishBadName
			result := h.newOnStatus(message.NetStreamOnStatusCodePublishFailed, "Publish failed.")

			l.Infof("Reject a Publish request: Response = %#v, Err = %+v", result, err)
			if err1 := h.sh.stream.NotifyStatus(chunkStreamID, timestamp, result); err1 != nil {
				return errors.Wrapf(err, "Failed to reply response: Err = %+v", err1)
			}

			return err
		}

		result := h.newOnStatus(message.NetStreamOnStatusCodePublishStart, "Publish succeeded.")
		if err := h.sh.stream.NotifyStatus(chunkStreamID, timestamp, result); err != nil {
			return err
		}
		l.Infof("Publisher accepted")

		h.sh.ChangeState(streamStateServerPublish)

		go myApi()

		return nil

	case *message.NetStreamPlay:
		l.Infof("Player is comming: %#v", cmd)

		streamCtx := &StreamContext{
			StreamID: h.sh.stream.streamID,
		}
		if err := h.sh.stream.userHandler().OnPlay(streamCtx, timestamp, cmd); err != nil {
			result := h.newOnStatus(message.NetStreamOnStatusCodePlayFailed, "Play failed.")

			l.Infof("Reject a Play request: Response = %#v, Err = %+v", result, err)
			if err1 := h.sh.stream.NotifyStatus(chunkStreamID, timestamp, result); err1 != nil {
				return errors.Wrapf(err, "Failed to reply response: Err = %+v", err1)
			}

			return err
		}

		result := h.newOnStatus(message.NetStreamOnStatusCodePlayStart, "Play succeeded.")
		if err := h.sh.stream.NotifyStatus(chunkStreamID, timestamp, result); err != nil {
			return err
		}
		l.Infof("Player accepted")

		h.sh.ChangeState(streamStateServerPlay)

		return nil

	default:
		return internal.ErrPassThroughMsg
	}
}

func myApi() {
	apiURL1 := "http://localhost:7001/data?channel=" + "$tart"
	res2, err := http.Get(apiURL1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res2.Body.Close()
}

func (h *serverDataInactiveHandler) newOnStatus(
	code message.NetStreamOnStatusCode,
	description string,
) *message.NetStreamOnStatus {
	// https://helpx.adobe.com/adobe-media-server/ssaslr/netstream-class.html#netstream_onstatus
	level := message.NetStreamOnStatusLevelStatus
	switch code {
	case message.NetStreamOnStatusCodeConnectFailed:
		fallthrough
	case message.NetStreamOnStatusCodePlayFailed:
		fallthrough
	case message.NetStreamOnStatusCodePublishBadName, message.NetStreamOnStatusCodePublishFailed:
		level = message.NetStreamOnStatusLevelError
	}

	return &message.NetStreamOnStatus{
		InfoObject: message.NetStreamOnStatusInfoObject{
			Level:       level,
			Code:        code,
			Description: description,
		},
	}
}
