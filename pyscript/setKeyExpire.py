from redis import StrictRedis
def setExpiredKeys():
    try:
        if redis_pass == 'none':
            redisclient = StrictRedis(host=redis_host, port=redis_port,db=1)
        else:
            redisclient = StrictRedis(host=redis_host, port=redis_port, password=redis_pass, db=1)
        for key in redisclient.scan_iter(match="cd_session_*",count=10):
            print(key)
            #执行redis命令
            #keyttl = redisclient.execute_command('ttl', key)
            #if keyttl == -1: # 此处扫到key，可以进行导出处理或者执行命令
            #redisclient.expire(key, 7200)
    except Exception as e:
        raise e
if __name__ == '__main__':
    redis_host = ''
    redis_port = 6379
    redis_pass = ''
    setExpiredKeys()