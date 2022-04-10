local idpDiscoveryUrl = ngx.var.idp_discovery_url       -- "http://oidc-stub:8080/.well-known/openid-configuration"
local acceptedAudience = ngx.var.idp_client_id          -- "loginapp-p"
local delegatedLoginUrl = ngx.var.delegated_login_url   -- "http://localhost:8081/loginapp/api/start"
local authTokenLocation = ngx.var.auth_token_location   -- "cookie:loginapp_ID_token"

local opts = {
    discovery = idpDiscoveryUrl,
    token_signing_alg_values_expected = { "RS256" },
    jwk_expires_in = 24 * 60 * 60,
    access_token_expires_leeway = 240,
    auth_accept_token_as = authTokenLocation == 'header' and authTokenLocation or ('cookie:'..authTokenLocation),
    ssl_verify = "no"
}

local contains = function (haystack, needle)
    for key, value in pairs(haystack) do
        if (value == needle) then
            return true
        end
    end
    return false
end

-- Verify token;
-- Fetches jwks from `discovery`, and caches the keys for `jwk_expires_in`.
-- Checkes that the request has a token at `auth_accept_token_as`,
-- that the token is not expired or expires within `access_token_expires_leeway`, and
-- that the token is signed using an appropriate algorithm `token_signing_alg_values_expected`.
local res, err = require("resty.openidc").bearer_jwt_verify(opts)

local hasAcceptedAudience = false
if res and res.aud and type(res.aud) == "string" and res.aud == acceptedAudience then
    hasAcceptedAudience = true
elseif res and res.aud and type(res.aud) == "table" and contains(res.aud, acceptedAudience) then
    hasAcceptedAudience = true
end

if err or not res or not hasAcceptedAudience then
    ngx.status = 302
    if err then
        ngx.header['Redirect-Reason'] = err
    elseif not res then
        ngx.log(ngx.INFO, "Redirect-Reason: no result")
        ngx.header['Redirect-Reason'] = "no result"
    elseif not res.aud then
        ngx.log(ngx.INFO, "Redirect-Reason: missing audience")
        ngx.header['Redirect-Reason'] = "missing audience"
    elseif not hasAcceptedAudience then
        ngx.log(ngx.INFO, "Redirect-Reason: invalid audience "..res.aud)
        ngx.header['Redirect-Reason'] = "invalid audience "..res.aud
    else
        ngx.log(ngx.INFO, "Redirect-Reason: unknown reason")
        ngx.header['Redirect-Reason'] = "unknown reason"
    end

    local full_url = ngx.var.real_scheme.."://"..ngx.var.http_host..ngx.var.request_uri
    ngx.redirect(delegatedLoginUrl.."?url="..ngx.escape_uri(full_url))
end