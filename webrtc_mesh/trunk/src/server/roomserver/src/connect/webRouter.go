package connect

import (
	"roomserver/src/logic"
	"roomserver/src/model"
	"roomserver/src/model/iceBalance"
	"roomserver/src/util"
	"crypto/tls"
	"errors"
	Xlog "github.com/cihub/seelog"
	"github.com/json-iterator/go"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"time"

	"net/http"
	"strconv"
)


// Collider 碰撞机
type WebRouter struct {
	collier *logic.Collider		// 房间table 字典
	WsCount int
}

func NewCollider(rs string) *WebRouter {
	return &WebRouter{
		collier: logic.NewCollider(),
		WsCount: 0,
	}
}

// Run starts the collider server and blocks the thread until the program exits.
func (router *WebRouter) Run(p int, useTls bool) {
	// 注意：这里 url/ws
	http.Handle("/ws", websocket.Handler(router.wsHandler))	// 主要
	//http.HandleFunc("/status", router.httpStatusHandler)
	http.HandleFunc("/", router.httpHandler)
	http.HandleFunc("/voiptoken", router.httpHandleVoIPToken)
	http.HandleFunc("/voiplogon", router.httpHandleVoIPLogon)
	http.HandleFunc("/websocketlist", router.httpHandleWebsocketList)
	http.HandleFunc("/icebandwidthload", router.httpHandleBandwidthIceLoad) // 主要

	var e error

	pstr := ":" + strconv.Itoa(p)
	//Xlog.Error("here only for notify the admin that roomserver have start, it is not error.")
	//Xlog.Info("pstr ", pstr)
	if useTls {
		config := &tls.Config {
			// Only allow ciphers that support forward secrecy for iOS9 compatibility:
			// https://developer.apple.com/library/prerelease/ios/technotes/App-Transport-Security-Technote/
			CipherSuites: []uint16 {
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			},
			PreferServerCipherSuites: true,
		}
		server := &http.Server{ Addr: pstr, Handler: nil, TLSConfig: config }

		e = server.ListenAndServeTLS("/cert/cert.pem", "/cert/key.pem")
	} else {
		e = http.ListenAndServe(pstr, nil)
	}

	if e != nil {
		Xlog.Critical("Run: " + e.Error())
	}
}

// httpStatusHandler is a HTTP handler that handles GET requests to get the
// status of collider.
func (c *WebRouter) httpStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")

	Xlog.Info("httpStatusHandler")
	//rp := router.dash.getReport(router.roomTable)
	//enc := json.NewEncoder(w)
	//if err := enrouter.Encode(rp); err != nil {
	//	err = errors.New("Failed to encode to JSON: err=" + err.Error())
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	router.dash.onHttpErr(err)
	//}
}

// httpHandler is a HTTP handler that handles GET/POST/DELETE requests.
// POST request to path "/$ROOMID/$CLIENTID" is used to send a message to the other client of the room.
// $CLIENTID is the source client ID.
// The request must have a form value "msg", which is the message to send.
// DELETE request to path "/$ROOMID/$CLIENTID" is used to delete all records of a client, including the queued message from the client.
// "OK" is returned if the request is valid.
func (c *WebRouter) httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE, GET")

	Xlog.Info("httpHandler, url ", r.URL.Path)
	http.ServeFile(w, r, "index.html")


	io.WriteString(w, "OK\n")
}

func (c *WebRouter) httpHandleVoIPToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE, GET")
	Xlog.Info("httpHandleVoIPToken into, method " , r.Method)
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			c.httpError("Failed to read request body: "+ err.Error(), w)
			Xlog.Error("Failed to read request body: "+ err.Error())
			return
		}
		m := string(body)
		if m == "" {
			c.httpError("Empty request body", w)
			Xlog.Error("Empty request body")
			return
		}
		Xlog.Info("httpHandleVoIPToken recive ", m)
		var tokenReq  util.TokenRequest
		err = jsoniter.Unmarshal(body, &tokenReq)
		if err  != nil {
			c.httpError("body parse faild ", w)
			Xlog.Error("body parse faild ")
			return
		}
		now := time.Now().Unix()		//返回秒
		validTimeInSeconds := int64(3600*24)
		expiredTsInSeconds := now + validTimeInSeconds
		token, err := util.GenerateSignalingToken(tokenReq.UID, tokenReq.AppID,
			util.DefaultAppCertificate, expiredTsInSeconds)
		if err == nil {
			w.WriteHeader(200)
			io.WriteString(w, token)
		} else {
			c.httpError("GenerateSignalingToken failed", w)
			Xlog.Error("GenerateSignalingToken faild ")
		}
		return
	default:
		return
	}

	c.httpError("httpHandleVoIPToken failed ", w)
	Xlog.Error("httpHandleVoIPToken faild ")
}
func (c *WebRouter) httpHandleBandwidthIceLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE, GET")
	//Xlog.Info("httpHandleBandwidthIceLoad into, method " , r.Method)
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(201)
			c.httpError("Failed to read request body: "+ err.Error(), w)
			Xlog.Error("Failed to read request body: "+ err.Error())
			return
		}
		m := string(body)
		if m == "" {
			w.WriteHeader(201)
			c.httpError("Empty request body", w)
			Xlog.Error("Empty request body")
			return
		}

		jsonData := jsoniter.Get(body, "cmd")
		cmd := jsonData.ToString()
		if cmd == iceBalance.CMD_ICE_REPORT_RX_TX_RATE {
			Xlog.Debugf("httpHandleBandwidthIceLoad recive %s", m)
		} else {
			Xlog.Infof("httpHandleBandwidthIceLoad recive %s", m)
		}
		ret := iceBalance.ICE_NO_FOUND
		switch cmd {
		case iceBalance.CMD_ICE_REGISTER:
			ret = c.collier.HandleIceRegister(body, r.RemoteAddr)
			break
		case iceBalance.CMD_ICE_DEREGISTER:
			ret = c.collier.HandleIceDeregister(body)
			break
		case iceBalance.CMD_ICE_REPORT_RX_TX_RATE:
			ret = c.collier.HandleIceReportIceRxTxRate(body)
			break
		}
		w.WriteHeader(200)
		io.WriteString(w, iceBalance.ICE_ERROR_MSG[ret])
		return
	default:
		return
	}
	w.WriteHeader(201)
	c.httpError("httpHandleVoIPToken failed ", w)
	Xlog.Error("httpHandleVoIPToken faild ")
}

func (c *WebRouter) httpHandleVoIPLogon(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE, GET")

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			c.httpError("Failed to read request body: "+err.Error(), w)
			return
		}
		m := string(body)
		if m == "" {
			c.httpError("Empty request body", w)
			return
		}
		Xlog.Info("httpHandleVoIPLogon recive ", m)
	case "DELETE":

	default:
		return
	}
	w.WriteHeader(200)
	io.WriteString(w, "token csfdsfdsfdsfdsfds\n")
}

func (c *WebRouter) httpHandleWebsocketList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE, GET")
	Xlog.Info("httpHandleWebsocketList into, method " , r.Method)
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			c.httpError("Failed to read request body: "+err.Error(), w)
			return
		}
		m := string(body)
		if m == "" {
			c.httpError("Empty request body", w)
			return
		}
		Xlog.Info("httpHandleWebsocketList recive ", m)
	case "DELETE":

	case "GET":		// GET方法数据在URL，那怎么回复数据呢？
		Xlog.Info("url :", r.URL)
	default:
		return
	}
	w.WriteHeader(200)
	io.WriteString(w, "['127.0.0.1:9000/ws', '127.0.0.1:9000/ws']")
}

// wsHandler is a WebSocket server that handles requests from the WebSocket client in the form of:
// 1. { 'cmd': 'register', 'roomid': $ROOM, 'clientid': $CLIENT' },
// which binds the WebSocket client to a client ID and room ID.
// A client should send this message only once right after the connection is open.
// or
// 2. { 'cmd': 'send', 'msg': $MSG }, which sends the message to the other client of the room.
// It should be sent to the server only after 'regiser' has been sent.
// The message may be cached by the server if the other client has not joined.
//
// Unexpected messages will cause the WebSocket connection to be closed.
// 每个客户端都维持一个ws连接
/*
1,string 转为[]byte
var str string = "test"
var data []byte = []byte(str)

2,byte转为string
var data [10]byte 
byte[0] = 'T'
byte[1] = 'E'
var str string = string(data[:])

//先检测命令，然后使用结构体序列化结构体
*/


func (router *WebRouter) wsHandler(ws *websocket.Conn) {
	router.WsCount++
	Xlog.Infof("wsHandler into WsCount = %d, remoteIp:%s, %s",
		router.WsCount, ws.Request().Header.Get("Remote_addr"), ws.RemoteAddr().String())

	var cmd string
	var message []byte
//loop:
	for {
		err := ws.SetReadDeadline(time.Now().Add(time.Duration(WS_READ_TIMEOUT_SEC) * time.Second))
		if err != nil {
			router.wsError("ws.SetReadDeadline error: "+err.Error(), ws)
			Xlog.Errorf("wsHandler SetReadDeadline break, error:%s", err.Error())
			break
		}
		// 读取websocket消息数据
		err = websocket.Message.Receive(ws, &message)
		if err != nil {
			if err.Error() != "EOF" {
				router.wsError("websocket.JSON.Receive error: "+err.Error(), ws)
				Xlog.Error("websocket.JSON.Receive error: "+err.Error())
			}
			//Xlog.Info("websocket normal exit, code = ", err.Error()) // 正常断开则不打印
			break
		}

		jsonData := jsoniter.Get(message, "cmd")
		cmd = jsonData.ToString()
		if  (cmd != model.CMD_KEEP_LIVE  &&
			cmd != model.CMD_CANDIDATE &&
			cmd != model.CMD_OFFER &&
			cmd != model.CMD_ANSWER &&
			cmd != model.CMD_REPORT_INFO &&
			cmd != model.CMD_REPORT_STATS) {
			Xlog.Info("message: ", string(message))
		}
		
		beginTime := util.GetMillisecond()
		switch cmd {
		case model.CMD_JOIN:
			ret := router.collier.HandleJoin(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_LEAVE:
			ret := router.collier.HandleLeave(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_OFFER:
			ret := router.collier.HandleOffer(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_ANSWER:
			ret := router.collier.HandleAnswer(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_CANDIDATE:
			ret := router.collier.HandleCandidate(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_KEEP_LIVE:
			ret := router.collier.HandleKeepLive(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_TURN_TALK_TYPE:
			ret := router.collier.HandleTurnTalkType(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_PEER_CONNECTED:
			ret := router.collier.HandlePeerConnected(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_REPORT_INFO:
			ret := router.collier.HandleReportInfo(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		case model.CMD_REPORT_STATS:
			ret := router.collier.HandleReportStats(message, ws)
			if ret != model.ROOM_ERROR_SUCCESS {
				Xlog.Warnf("Cmd:%s error, it is %s", cmd, model.ROOM_ERROR_MSG[ret])
			}
			break
		}
		endTime := util.GetMillisecond()
		// 减少间隔性的打印
		if  cmd != model.CMD_KEEP_LIVE &&  cmd != model.CMD_CANDIDATE  &&  cmd != model.CMD_REPORT_STATS{
			Xlog.Infof("Handle cmd:%s consume %dms", cmd, endTime-beginTime)
		}
	}
	//Xlog.Infof("wsHandler leave WsCount = %d, %p", router.WsCount, ws)

	// This should be unnecessary but just be safe.
	ws.Close()
}

func (c *WebRouter) httpError(msg string, w http.ResponseWriter) {
	err := errors.New(msg)
	http.Error(w, err.Error(), http.StatusInternalServerError)
	//router.dash.onHttpErr(err)
}

func (c *WebRouter) wsError(msg string, ws *websocket.Conn) {
	//err := errors.New(msg)
	//sendServerErr(ws, msg)
	//router.dash.onWsErr(err)
}
