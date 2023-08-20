package com.darren.rtcsdk.ws;

import com.darren.rtcsdk.NetState;

/**
 * Created by darren on 2019/1/3.
 * 592407834@qq.com
 */
public interface IWebSocket {
    boolean isOpen();

    void reConnect();

    void close();

    // 加入房间
    int joinRoom(String room);
    int exitRoom();
    //处理回调消息
    void handleMessage(String message);

    int sendIceCandidate(String desc, String remoteUid);

    int sendAnswer(String desc, String remoteUid);

    int sendOffer(String desc, String remoteUid);
    int sendKeepLive();
    int sendStats(String stats);

    int sendTurnTalkType(int i, boolean isMute);

    int sendPeerConnected(String remoteUid, String connectType);

    NetState getNetState();
}
