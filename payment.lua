
math.randomseed(os.time())

request = function()
  local sender_account = 10000000 + math.random(1, 1000)
  local receiver_account = 10000000 + math.random(1, 1000)

  if sender_account == receiver_account then
    receiver_account = receiver_account + 1
  end

  local body = string.format([[
  {
    "sender_account": "%d",
    "receiver_account": "%d",
    "amount": 100,
    "currency": "USD"
  }
  ]], sender_account, receiver_account)

  return wrk.format("POST", "/payments", {
    ["Content-Type"] = "application/json"
  }, body)
end