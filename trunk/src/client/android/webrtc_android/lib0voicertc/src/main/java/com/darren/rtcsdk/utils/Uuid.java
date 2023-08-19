package com.darren.rtcsdk.utils;

import java.util.Random;
import java.util.UUID;

public class Uuid {
    /**
     * 获得指定数目的UUID
     * @param number int 需要获得的UUID数量
     * @return String[] UUID数组
     */
    public static String[] getUUID(int number){
        if(number < 1){
            return null;
        }
        String[] retArray = new String[number];
        for(int i=0;i<number;i++){
            retArray[i] = getUUID();
        }
        return retArray;
    }

    /**
     * 获得一个UUID
     * @return String UUID
     */
    public static String getUUID(){
        String uuid = UUID.randomUUID().toString();
        //去掉“-”符号
        return uuid.replaceAll("-", "");
    }

    /**
     *
     * java通过UUID生成16位唯一订单号
     *
     *
     * */
    public static String get16BitUUId() {
        int first = new Random(10).nextInt(8) + 1;
        System.out.println(first);
        int hashCodeV = UUID.randomUUID().toString().hashCode();
        if (hashCodeV < 0) {//有可能是负数
            hashCodeV = -hashCodeV;
        }
        // 0 代表前面补充0
        // 4 代表长度为4
        // d 代表参数为正数型
        return first + String.format("%015d", hashCodeV);
    }

    /**
     * 生成指定位数的随机数
     * @param length
     * @return
     */
    public static String getRandom(int length){
        String val = "";
        Random random = new Random();
        for (int i = 0; i < length; i++) {
            val += String.valueOf(random.nextInt(10));
        }
        return val;
    }

    /**
     * 生成随机的用户ID
     * @return
     */
    public static String generateUserUUID() {
        return get16BitUUId();
    }
    /**
     * 生成随机的用户名
     * @return
     */
    public static String generateUserName(){
        return "client_" + getRandom(4);
    }

    /**
     * 生成随机的房间ID
     * @return
     */
    public static String  generateRoomUUID() {
        return get16BitUUId();
    }

    /**
     * 生成随机的房间名
     * @return
     */
    public static String generateRoomName(){
        return "room_" + getRandom(4);
    }
}