# coding=UTF-8
#!/usr/bin/env Python
import time
import sys
import urllib2
import json
import ConfigParser
import logging 
from logging.handlers import TimedRotatingFileHandler
from urllib import quote
import re
import smtplib
from bottle import template
from email.mime.text import MIMEText
from email.utils import formataddr
import string
import datetime #导入日期时间模块

CONFIG_FILE = '/data/rtc/ice_server/ice_net_monitor.ini'
LOG_FILE_PATH = '/data/rtc/ice_server/monitor_log/'

RegisterIceDone = False
ReportIceDone = False

# log函数
def getLogger():
    logger = logging.getLogger('ice_monitor')
    # logger.setLevel(logging.DEBUG)
    logger.setLevel(logging.INFO)

    # fh = logging.FileHandler('log/updateCache.log')
    # fh.setLevel(logging.DEBUG)

    # ch = logging.StreamHandler()
    # ch.setLevel(logging.DEBUG)

    logFileName = LOG_FILE_PATH + 'ice_net_monitor.log'
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

# 配置log
try:
    logger = getLogger()
except Exception, e:
    print e
    sys.exit(1)


STATS = []
def getRx(net_card):
    ifstat = open('/proc/net/dev').readlines()
    for interface in  ifstat:
        if net_card in interface:
            stat = float(interface.split()[1])
            STATS[0:] = [stat]

def getTx(net_card):
    ifstat = open('/proc/net/dev').readlines()
    for interface in  ifstat:
        if net_card in interface:
            stat = float(interface.split()[9])
            STATS[1:] = [stat]

def jsonPost(url, data):
    headers = {}
    headers['Content-Type'] = 'application/json; charset=utf-8'
    req = urllib2.Request(url, data, headers)
    page = urllib2.urlopen(req)
    res = page.read()
    page.close()
    return res

def sendHttpGetReq(url, data,  headers={}, timeout=10):
    try:
		# logger.info("url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers))
		req = urllib2.Request(url, data, headers)
		response = urllib2.urlopen(req, None, timeout)
		return 0, response
    except urllib2.HTTPError, e:  # https://docs.python.org/3/howto/urllib2.html#httperror
        logger.error("http error|url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers) + "|code=" + str(e.code) + "|reason=" + str(e.reason))
        return e.code, str(e.reason)
    except urllib2.URLError, e:
        logger.error("url error|url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers) + "|reason=" + str(e.reason))
        if "timed out" in str(e.reason):
            return 504, str(e.reason)
        else:
            return 1, str(e.reason)
    except Exception, e:
        logger.error("other error|url=" + url + "|timeout=" + str(timeout) + "|headers=" + str(headers) + "|reason=" + str(e))
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

def sendMail(cfg, content):
        ret = True
        ipAddr = cfg.get('mail', 'ipAddr')
        receivers = cfg.get('mail', 'receivers')
        sender = cfg.get('mail', 'sender')
        passwd = cfg.get('mail', 'passwd')
        smtp_server = cfg.get('mail', 'smtp_server')
        smtp_port = cfg.get('mail', 'smtp_port')

        # today = datetime.date.today() #获得今天的日期
        nowTime=datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        
        item = [ipAddr, 'ice_net_monitor(网速实时监控)', 'process_monitor', content, nowTime]
        print item
        # items = []
        # items.apend(item)
        html = creat_html(item, nowTime)
        # print [receivers]
        # msg = MIMEText(html, 'html', 'utf-8') # 网页格式
        msg = MIMEText(html, 'html', 'utf-8')
        msg['From']=formataddr([sender,sender])   #括号里的对应发件人邮箱昵称、发件人邮箱账号
        msg['To']=";".join([receivers])   #收件人列表
        msg['Subject']="ICE monitor ip:%s report" % ipAddr #邮件的主题，也可以说是标题
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

def registerIce(cfg):
    # 封装body
    body = {}
    body['cmd']         	= 'iceRegister'
    body['iceType']         = cfg.get('ice', 'iceType')
    body['bundlePolicy']    = cfg.get('ice', 'bundlePolicy')
    body['rtcpMuxPolicy']   = cfg.get('ice', 'rtcpMuxPolicy')
    body['iceTransportPolicy'] = cfg.get('ice', 'iceTransportPolicy')
    body['ip']              = cfg.get('ice', 'ip')
    body['turnPort']        = cfg.getint('ice', 'turnPort')
    body['stunPort']        = cfg.getint('ice', 'stunPort')
    body['maxBandwidth']    = cfg.getint('ice', 'maxBandwidth')
    body['username']        = cfg.get('ice', 'username')
    body['credential'] = cfg.get('ice', 'credential')
    url = cfg.get('ice', 'roomserverUrl')
    j_body = json.dumps(body)
    # print j_body
    headers = {}
    headers['Content-Type'] = 'application/json; charset=utf-8'
    ret, res = sendHttpGetReq(url, j_body, headers, 5)

    if ret == 0 and res.getcode() == 200:
        content = res.read()
        # print content
        if content == 'successful' or content == 'ice_have_registered':
            return True			# 注册成功
        else:
            logger.error("registerIce failed:%s" % content)
            return False    
    else:
        logger.error("http failed:%s" % res)
        return False

def tryRegisterIce(cfg):
    logger.info('tryRegisterIce....')
    registerCount = 0
    tryMaxregisterCount  = cfg.getint('register', 'tryMaxregisterCount')   # 累计5分钟发邮件报告一次异常
    registerTimeoutSleep = cfg.getint('register', 'registerTimeoutSleep') 
    global RegisterIceDone
    while True:
        ret = registerIce(cfg)
        if ret == True: 
            print "register ice successful" 
            logger.info( "register ice successful")
            if RegisterIceDone == False:
                # 发送邮件
                ret = sendMail(cfg, "ip %s register %s successful" % (cfg.get('ice', 'ip'), cfg.get('ice', 'roomserverUrl')))
                if ret:
                    logger.info('sendMail ok')
                    RegisterIceDone = True
                else:
                    logger.error('sendMail failed')    
            break             
        else:  
            logger.info('wait %d sec, continue try register....' % registerTimeoutSleep)
            if RegisterIceDone == True:
                # 发邮件报错
                print "can't register" 
                logger.error("can't register %s" % (cfg.get('ice', 'roomserverUrl')))
                # 发送邮件
                ret = sendMail(cfg, "can't register %s" % (cfg.get('ice', 'roomserverUrl')))
                if ret:
                    logger.info('sendMail ok')
                    RegisterIceDone = False
                else:
                    logger.error('sendMail failed')    
            time.sleep(registerTimeoutSleep)	# 隔registerTimeoutSleep秒再重新注册 
    
    logger.info('tryRegisterIce finish.')

def reportRxTxRate(cfg, rx, tx):
    # 封装body
    body = {}
    body['cmd'] 	= 'iceReportRxTxRate'
    body['ip']  	= cfg.get('ice', 'ip')
    body['rxRate']	= rx
    body['txRate'] 	= tx

    url = cfg.get('ice', 'roomserverUrl')
    j_body = json.dumps(body)
    # print j_body
    headers = {}
    headers['Content-Type'] = 'application/json; charset=utf-8'
    ret, res = sendHttpGetReq(url, j_body, headers, 5)
    if ret == 0 and res.getcode() == 200:
        content = res.read()
        # print content
        if content == 'ice_no_found': 
            # 重新注册
            tryRegisterIce(cfg)
            return False			# 注册成功
        elif content == 'successful':
            return True    
        else:
            logger.error("reportRxTxRate failed:%s" % content)
            return False    
    else:
        logger.error("http failed:%s" % res)
        return False
#	print RX_RATE ,'KB ',TX_RATE ,'KB'
#一个python的文件有两种使用的方法，第一是直接作为脚本执行，第二是import到其他的python脚本#
#中被调用（模块重用）执行。因此if __name__ == 'main': 
#的作用就是控制这两种情况执行代码的过程，在if __name__ == 'main': 
#下的代码只有在第一种情况下（即文件作为脚本直接执行）才会被执行，而import到其他脚本中是不
#会被执行的。
if __name__ == '__main__':
    # 初始化配置
    logger.info('loading config...')
    cfg = ConfigParser.SafeConfigParser()
    try:
        cfg.read(CONFIG_FILE)
    except Exception as e:
        logger.error('read config failed:' + str(e))
        exit(1)
    # 读取配置文件 
    logger.info('iceType        :%s' % (cfg.get('ice', 'iceType') ))  
    logger.info('bundlePolicy   :%s' % (cfg.get('ice', 'bundlePolicy')  ))
    logger.info('rtcpMuxPolicy  :%s' % ( cfg.get('ice', 'rtcpMuxPolicy')))
    logger.info('iceTransportPolicy:%s' % (cfg.get('ice', 'iceTransportPolicy')))
    logger.info('ip             :%s' % (cfg.get('ice', 'ip')))
    logger.info('turnPort       :%d' % (cfg.getint('ice', 'turnPort')))  
    logger.info('stunPort       :%d' % (cfg.getint('ice', 'stunPort')))  
    logger.info('maxBandwidth   :%s' % (cfg.get('ice', 'maxBandwidth')))  
    logger.info('username       :%s' % (cfg.get('ice', 'username')))
    logger.info('credential     :%s' % (cfg.get('ice', 'credential')))
    logger.info('net_card       :%s' % (cfg.get('ice', 'net_card')))
    logger.info('roomserverUrl  :%s' % (cfg.get('ice', 'roomserverUrl')))
    
    # ret = sendMail(cfg, 'I am Darren')
    # if ret:
    #     print 'sendMail ok'
    # else:
    #     print 'sendMail failed'    
    # 向服务器节点请求注册，直到注册成功，如果不成功则发邮件通知
    
    tryRegisterIce(cfg)
    # 执行到此则说明注册成功
    net_card = cfg.get('ice', 'net_card')
    getRx(net_card)
    getTx(net_card)
    interval = 5
    registerCount = 0
    tryMaxregisterCount  = cfg.getint('register', 'tryMaxregisterCount')   # 累计5分钟发邮件报告一次异常
    registerTimeoutSleep = cfg.getint('register', 'registerTimeoutSleep') 
    interval = cfg.getint('register', 'interval') 
    logger.info('start report tx rx info....')
    while True:
        time.sleep(interval)
        rxstat_o = list(STATS)
        getRx(net_card)
        getTx(net_card)
        RX = float(STATS[0])
        RX_O = rxstat_o[0]
        TX = float(STATS[1])
        TX_O = rxstat_o[1]
        rxRate = int((RX - RX_O)/1024/interval*8)
        txRate = int((TX - TX_O)/1024/interval*8)
        ret = reportRxTxRate(cfg, rxRate, txRate)
        global ReportIceDone
        # 报告带宽信息
        if ret == True:
            if ReportIceDone == False:
                ReportIceDone = True
                logger.info("start report rx tx info to %s" % (cfg.get('ice', 'roomserverUrl')))
                # 发送邮件
                ret = sendMail(cfg, "start report rx tx info to %s" % (cfg.get('ice', 'roomserverUrl')))
                if ret:
                    logger.info('sendMail ok')
                else:
                    logger.error('sendMail failed')      
        else:
            if ReportIceDone == True:
                ReportIceDone = False
                # 邮件报警
                print "can't report rx tx info" 
                logger.error("can't report rx tx info to %s" % (cfg.get('ice', 'roomserverUrl')))
                # 发送邮件
                ret = sendMail(cfg, "can't report rx tx info to %s" % (cfg.get('ice', 'roomserverUrl')))
                if ret:
                    logger.info('sendMail ok')
                else:
                    logger.error('sendMail failed')                       
