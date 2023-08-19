package com.darren.rtcsdk.utils;

import android.content.Context;

/**
 * Created by darren on 2019/1/2.
 * 592407834@qq.com
 */
public class Utils {

    public static int dip2px(Context context, float dpValue) {
        final float scale = context.getResources().getDisplayMetrics().density;
        return (int) (dpValue * scale + 0.5f);
    }

    public static String wssUrl = "wss://www.0voice.com:8088/ws";
    public static String videoUrl="https://www.0voice.com:8088";
}
