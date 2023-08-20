#encoding: utf-8
import time
import os
import sys
import fcntl
import urllib2
from urllib import quote
import re
import logging
from logging.handlers import TimedRotatingFileHandler
import ConfigParser
import math
import datetime #导入日期时间模块
from bottle import template
import smtplib
from bottle import template
from email.mime.text import MIMEText
from email.utils import formataddr
import string

SERVER_PATH='/data/rtc/ice_server/'
CONFIG_FILE = '/data/rtc/ice_server/ice_server_monitor.ini'
LOG_FILE_PATH = '/data/rtc/ice_server/monitor_log/'
INTERVAL = 5
SECONDS_OF_MIN = 60 #便于调试!!!
REPORT_ADDR = 'http://184.173.216.6:8094/sendmail?sender={sender}&password={passwd}&receiver={receiver}&mode={mode}&descr={desc}&address={addr}&title={title}'



def processInstance():
    ret = os.popen('ps -ef |grep "python.*ice_server_monitor.py" |grep -v grep |grep -v sudo |wc -l').read().strip()
    if ret != str(1):
        print 'another instance running,ret:' + str(ret)
        sys.exit(1)

def isProcessAlive(aliveCmd):
    ret = os.system(aliveCmd)
    return ret == 0

def startProcess(startCmd):
    ret = os.system(startCmd)
    logger.info('system cmd {} ret:{}'.format(startCmd, ret))

def sendHttpGetReq(url, timeout=10, headers={}):
    try:
        # logger.info("url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers))
        req = urllib2.Request(url, headers=headers)
        response = urllib2.urlopen(req, None, timeout)
        # body:response.read(), head:response.info()
        return 0, response
    except urllib2.HTTPError, e:  # https://docs.python.org/3/howto/urllib2.html#httperror
        # logger.error("http error|url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers) + "|code=" + str(e.code) + "|reason=" + str(e.reason))
        return e.code, str(e.reason)
    except urllib2.URLError, e:
        # logger.error("url error|url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers) + "|reason=" + str(e.reason))
        if "timed out" in str(e.reason):
            return 504, str(e.reason)
        else:
            return 1, str(e.reason)
    except Exception, e:
        # logger.error("other error|url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers) + "|reason=" + str(e))
        if "timed out" in str(e):
            return 504, str(e)
        else:
            return 2, str(e)

def creat_html(info_items, the_day):
    html = """
    <html>
    <title>消息统计</title>
    <head></head>
    <body>
    <br></br>
    <h2 align=center>报警消息({{date}})</h2>
    <table width="90%" align=center border="0" bgcolor="#666666" cellpadding="8">
        <tr bgcolor="#DDDDDD">
            <th>IP地址</th>
            <th>进程名称</th>
            <th>所属类型</th>
            <th>详细描述</th>
            <th>时间</th>
        </tr>
        <tr align=center bgcolor="#FFFFFF">
            <td><font color="#33CC00">{{items[0]}}</font></td>
            <td><font color="#33CC00">{{items[1]}}</font></td>
            <td><font color="#33CC00">{{items[2]}}</font></td>
            <td><font color="#33CC00">{{items[3]}}</font></td>
            <td><font color="#33CC00">{{items[4]}}</font></td>
        </tr>
    </table> 
    </body>
    </html>
    """
    return template(html, items=info_items, date=the_day)

def sendMail(cfg, process, content):
        ret = True
        ipAddr = cfg.get('DEFAULT', 'ipAddr')
        receivers = cfg.get('DEFAULT', 'receivers')
        sender = cfg.get('DEFAULT', 'sender')
        passwd = cfg.get('DEFAULT', 'passwd')
        smtp_server = cfg.get('DEFAULT', 'smtp_server')
        smtp_port = cfg.get('DEFAULT', 'smtp_port')

        # today = datetime.date.today() #获得今天的日期
        nowTime=datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        
        item = [ipAddr, process, 'process_monitor', content, nowTime]
        # print item
        # items = []
        # items.apend(item)
        html = creat_html(item, nowTime)
        # print [receivers]
        # msg = MIMEText(html, 'html', 'utf-8') # 网页格式
        msg = MIMEText(html, 'html', 'utf-8')
        msg['From']=formataddr([sender,sender])   #括号里的对应发件人邮箱昵称、发件人邮箱账号
        msg['To']=";".join([receivers])   #收件人列表
        msg['Subject']="实时音视频ICE服务器报警信息" #邮件的主题，也可以说是标题
        receivers = string.splitfields(receivers, ",")
        try:
            server = smtplib.SMTP(smtp_server, smtp_port)  #发件人邮箱中的SMTP服务器，端口是25
            server.login(sender, passwd)    #括号中对应的是发件人邮箱账号、邮箱密码
            server.sendmail(sender, receivers, msg.as_string())   #括号中对应的是发件人邮箱账号、收件人邮箱账号、发送邮件
            server.quit()   #这句是关闭连接的意思
        except Exception as e:
            print e
            logger.error(e)
        return ret


def reportError(reportUrl):
    logger.info('report url:' + reportUrl)
    ret, response = sendHttpGetReq(reportUrl, 10, {})
    if ret == 0:
       logger.info('report succeed')
        #print response.read()
    else:
        logger.info('report failed, ret:%d,%s' % (ret, response))
        ret, response = sendHttpGetReq(reportUrl, 10, {})
        if ret != 0:
            logger.error('report failed again, ret:%d,%s' % (ret, response))

def monitorProcess(process, cfg, state):
    # 进程拉起失败后降低监控频率
    if not state['isAlive'] and state['skipTimes'] < math.pow(2, state['startTimes']) and state['skipTimes'] < 10:
        state['skipTimes'] += 1
        return
    if isProcessAlive(cfg.get(process, 'aliveCmd')):
        state['isAlive'] = True
        state['startTimes'] = 0
        state['skipTimes'] = 0                                      
        return 
    logger.info(process + ' is not alive')                        
    # 已打开的日志文件描述符fd将传递给nginx进程!                  
    startProcess(cfg.get(process, 'startCmd'))                 
    time.sleep(2)
    if isProcessAlive(cfg.get(process, 'aliveCmd')):
        logger.info('start ' + process + ' succeed')
        # reportUrl = REPORT_ADDR.format(sender=cfg.get(process, 'sender'), passwd=quote(cfg.get(process, 'passwd')), \
        #         receiver=cfg.get(process, 'receiver'), mode='process_monitor', title=quote(cfg.get(process, 'title')), \
        #         desc= quote(process + ' died, restart ok'), addr=cfg.get(process, 'ipAddr'))
        # reportError(reportUrl)
        sendMail(cfg, process, process + ' died, restart ok')
        state['isAlive'] = True
        state['startTimes'] = 0
        state['skipTimes'] = 0
    else:
        logger.info('start ' + process + ' failed')
        # reportUrl = REPORT_ADDR.format(sender=cfg.get(process, 'sender'), passwd=quote(cfg.get(process, 'passwd')),\
        # receiver=cfg.get(process, 'receiver'), mode='process_monitor', title=quote(cfg.get(process, 'title')), \
        # desc= quote(process + ' died, restart fail'), addr=cfg.get(process, 'ipAddr'))

        now = int(time.time())
        #print now , state['lastReportTime'] + SECONDS_OF_MIN * cfg.getint(process, 'reportInterval')
        if state['isAlive'] or now > state['lastReportTime'] + SECONDS_OF_MIN * cfg.getint(process, 'reportInterval'):
            # reportError(reportUrl)
            sendMail(cfg, process, process + ' died, restart fail')
            state['lastReportTime'] = now
        
        state['isAlive'] = False
        state['startTimes'] += 1
        state['skipTimes'] = 1

def getLogger():
    logger = logging.getLogger('ice_server_monitor')
    # logger.setLevel(logging.DEBUG)
    logger.setLevel(logging.INFO)

    # fh = logging.FileHandler('log/updateCache.log')
    # fh.setLevel(logging.DEBUG)

    # ch = logging.StreamHandler()
    # ch.setLevel(logging.DEBUG)

    logFileName = LOG_FILE_PATH + 'ice_server_monitor.log'
    dailyFH = TimedRotatingFileHandler(filename=logFileName, when='MIDNIGHT', interval=1, backupCount=30)
    dailyFH.setLevel(logging.INFO)

    formatter = logging.Formatter('%(asctime)s|%(levelname)s|%(process)d|%(funcName)s|%(lineno)d|%(message)s')
    # fh.setFormatter(formatter)
    # ch.setFormatter(formatter)
    dailyFH.setFormatter(formatter)

    logger.propagate = False
    # logger.addHandler(fh)
    # logger.addHandler(ch)
    logger.addHandler(dailyFH)

    return logger

if __name__ == '__main__':
    processInstance()
    try:
        logger = getLogger()
    except Exception, e:
        print e
        sys.exit(1)
    # 初始化配置
    logger.info('loading config...')
    cfg = ConfigParser.SafeConfigParser()
    try:
        cfg.read(CONFIG_FILE)
    except Exception as e:
        logger.error('read config failed:' + str(e))
        exit(1)
    stateDict = dict({})
    for process in cfg.sections():
        stateDict[process] = dict({'isAlive':True, 'startTimes':0, 'skipTimes':0, 'lastReportTime':0})

    os.chdir(SERVER_PATH)

    logger.info('start monitoring...')
    sendMail(cfg, 'ice_server_monitor', 'ice_server_monitor start...')
    #print stateDict
    while True:
        for process in cfg.sections():
            monitorProcess(process, cfg, stateDict[process])
        time.sleep(INTERVAL)
    logger.info('end monitoring')
    sys.exit(0)
