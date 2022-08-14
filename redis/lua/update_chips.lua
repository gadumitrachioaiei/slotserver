if redis.call('EXISTS', KEYS[1]) == 0 then
  if redis.call('SET', KEYS[1], ARGV[1]) == false then
    return -1
  end
end
local v = redis.call('GET', KEYS[1])
if v == false then
  return -1
end
local chips = v + ARGV[2]
if chips < 0 then
  return -1
end
if redis.call('SET', KEYS[1], chips) == false then
  return -1
end
return chips