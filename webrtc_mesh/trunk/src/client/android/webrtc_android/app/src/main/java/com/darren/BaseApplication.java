package com.darren;

import android.app.Application;

import com.squareup.leakcanary.LeakCanary;

/**
 * Created by darren on 2019/4/5.
 * 592407834@qq.com
 */
public class BaseApplication extends Application {

    @Override
    public void onCreate() {
        super.onCreate();
        // 先关闭内存泄露检测
//        if (LeakCanary.isInAnalyzerProcess(this)) {
//            return;
//        }
//        LeakCanary.install(this);
    }
}


