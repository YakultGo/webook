-- 验证码对应的redis中的key
-- phone_code:login:xxxxxx
local key = KEYS[1]
-- 验证码次数
-- phone_code:login:xxxxxx:cnt
local cntKey = key..":cnt"
-- 验证码
-- xxxxxx
local val = ARGV[1]
-- 过期时间
local ttl = redis.call("ttl", key)
if ttl == -1 then
    --   key存在，但是没有过期时间
    return -2
    --   -2表示不存在该key，540 = 600 - 60
elseif tt1 == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    return -1
end