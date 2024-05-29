from redis import StrictRedis
def getKeys():
    try:
        if redis_pass == 'none':
            redisclient = StrictRedis(host=redis_host, port=redis_port,db=1)
        else:
            redisclient = StrictRedis(host=redis_host, port=redis_port, password=redis_pass, db=1)
        for key in redisclient.scan_iter(match="cd_session_*",count=10):
            value = redisclient.get(key)
            print(key,value)
    except Exception as e:
        raise e
if __name__ == '__main__':
    redis_host = ''
    redis_port = 6379
    redis_pass = ''
    getKeys()