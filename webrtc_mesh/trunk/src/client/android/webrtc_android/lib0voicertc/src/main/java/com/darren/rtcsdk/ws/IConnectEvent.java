package com.darren.rtcsdk.ws;

/**
 * Created by darren on 2019/4/5.
 * 592407834@qq.com
 */
public interface IConnectEvent {
    void onSuccess();
    void onFailed(String msg);
}
