counter = 1

request = function()
  local account = 10000000 + counter
  counter = counter + 1

  local body = string.format([[
  {
    "deposit_account": "%d",
    "amount": 1000,
    "currency": "USD"
  }
  ]], account)

  return wrk.format("POST", "/deposits", {
    ["Content-Type"] = "application/json"
  }, body)
end