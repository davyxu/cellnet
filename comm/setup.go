package comm

import "github.com/davyxu/cellnet"

func ProcSysMsg(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		if userFunc == nil {
			return nil
		}

		switch ev := raw.(type) {

		case *cellnet.SessionConnectErrorEvent:
			userFunc(&cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionConnectError{}})
		case *cellnet.SessionClosedEvent:
			userFunc(&cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionClosed{}})
		case *cellnet.SessionAcceptedEvent:
			userFunc(&cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionAccepted{}})
		case *cellnet.SessionConnectedEvent:
			userFunc(&cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionConnected{}})
		}

		return userFunc(raw)
	}
}
