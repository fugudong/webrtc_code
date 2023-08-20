package com.darren;

import android.app.Activity;

import com.darren.rtcsdk.ui.ChatRoomActivity;
import com.darren.rtcsdk.ui.ChatSingleActivity;

/**
 * Created by darren on 2019/1/7.
 * 592407834@qq.com
 */
public class WebrtcUtil {
    private final static String TAG = "WebrtcUtil";

    // Videoconferencing
    public static void call(Activity activity, String appId, String roomId, String roomName, String uid, String uname, int mediaType) {
        ChatRoomActivity.openActivity(activity, appId, roomId, roomName, uid, uname, mediaType);
    }

    public static void callSingleRoom(Activity activity, String appId, String roomId, String roomName, String uid, String uname, int mediaType) {
        ChatSingleActivity.openActivity(activity, appId, roomId, roomName, uid, uname, mediaType);;
    }
}
