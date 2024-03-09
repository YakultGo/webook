local key = KEYS[1]
-- 用户输入的验证码
local expectCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 重试次数用完或者验证码失效
    return -1
elseif expectCode == code then
    -- 将验证码置为无效
    redis.call("set", cntKey, -1)
    return 0
else
    -- 验证码错误，重试次数减一
    redis.call("decr", cntKey)
    return -2
end
