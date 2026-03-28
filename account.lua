wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

counter = 1

request = function()
  local user_id = "user_" .. counter
  counter = counter + 1

  local body = string.format([[
  {
    "user_id": "%s",
    "currency": "USD"
  }
  ]], user_id)

  return wrk.format("POST", "/accounts", {
    ["Content-Type"] = "application/json"
  }, body)

end