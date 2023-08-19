package com.darren;

import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.View;
import android.widget.EditText;

import com.darren.webrtc.R;
import com.darren.rtcsdk.bean.MediaType;

import java.util.UUID;


/**
 * Created by darren on 2019/11/7.
 * 592407834@qq.com
 */
public class MainActivity extends AppCompatActivity {
    private EditText et_room;
    private EditText et_username;
    //获取32位uuid工具类    ，此类事java自带的，不需要导包
    public static String get32UUID() {
        String uuid = UUID.randomUUID().toString().trim().replaceAll("-", "");
        return uuid;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        Toolbar toolbar = findViewById(R.id.toolbar);
        setSupportActionBar(toolbar);
        initView();
        initVar();

    }

    private void initView() {
        et_room = findViewById(R.id.et_room);
        et_username = findViewById(R.id.et_username);
    }

    private void initVar() {
        et_room.setText("100");
        et_username.setText("Android_01");
    }



    public void RandomUser(View view) {
        int randomValue = (int) ((Math.random()*9+1)*100);        // 生成2位随机整数
        et_username.setText("Android_" + String.valueOf(randomValue));
    }
    public void RandomRoom(View view) {
        int randomValue = (int) ((Math.random()*9+1)*10000);        // 生成5位随机整数
        et_room.setText(String.valueOf(randomValue));
    }


    public void JoinVideoRoom(View view) {
        WebrtcUtil.call(this, "10000", et_room.getText().toString().trim(),
                "零声学院", get32UUID(), String.valueOf(et_username.getText()), MediaType.TYPE_VIDEO);

    }
    public void JoinAudioRoom(View view) {
        WebrtcUtil.call(this, "10000", et_room.getText().toString().trim(),
                "零声学院", get32UUID(), String.valueOf(et_username.getText()), MediaType.TYPE_AUDIO);
    }
    public void JoinSingleAudioRoom(View view) {
        WebrtcUtil.callSingleRoom(this, "10000", et_room.getText().toString().trim(),
                "零声学院", get32UUID(), String.valueOf(et_username.getText()), MediaType.TYPE_AUDIO);
    }

    public void JoinSingleVideoRoom(View view) {
        WebrtcUtil.callSingleRoom(this, "10000", et_room.getText().toString().trim(),
                "零声学院", get32UUID(), String.valueOf(et_username.getText()), MediaType.TYPE_VIDEO);

    }
}
