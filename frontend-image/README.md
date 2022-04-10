# Frontend-image

Docker image that secures your frontend resources. Used with an instance of `login-app` for delegated login. 

## Configuration
| Property               | Required | Description                                                                                                                                                                           |
|------------------------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| APP_NAME               | Yes      | Name of application. Used as the app's context-path.                                                                                                                                  | 
| APP_VERSION            | Yes      | Version of application. Used on selftest page.                                                                                                                                        | 
| IDP_DISCOVERY_URL      | Yes      | Url to discovery-url for idp (e.g something ending with .well-known/openid-configuration)                                                                                             |
| IDP_CLIENT_ID          | Yes      | Username for usage against IDP                                                                                                                                                        |
| DELEGATED_LOGIN_URL    | Yes      | Url to `login-app`, e.g `http://www.domain.no/loginapp/api/start`                                                                                                                     |
| AUTH_TOKEN_RESOLVER    | Yes      | Where the application should check for a token. E.g `ID_token` (cookie name) eller `header` (Authorization header)                                                                    |
| CSP_DIRECTIVES         | No       | CSP-header. Default: `default-src: 'self'`                                                                                                                                            | 
| CSP_REPORT_ONLY        | No       | `true` or `false`. Sets `Report-Only`.                                                                                                                                                |
| REFERRER_POLICY        | No       | Default `origin`. Prevent url to be passed in http-header when navigating with links. [Read more](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy#examples) |

**NB** If `AUTH_TOKEN_RESOLVER` is set to `header` then the application will except to receive the token via the http-header `Authorization: Bearer <token>`.

## Proxying

The application supports proxy setups. This is done by adding your own nginx-files into the `/nginx` folder.
For example;
```nginx
# file: proxy.nginx
# NB: trailing-slashes prevents the whole path being forwared
location  /frontend/proxy/open-endpoint/ {
    proxy_pass http://echo-server/;
}
location  /frontend/proxy/authenticated-endpoint/ {
    access_by_lua_file oidc_protected.lua;
    proxy_pass http://echo-server/;
}

# file: Dockerfile
# Add the newly created nginx-file, which is pick up by the image during startup.
COPY proxy.nginx /nginx
```

Other types of customazations:
- Remove cookies;  `proxy_set_header Cookie "";`
- Rename cookies before doing proxy-call; `proxy_set_header Cookie "new_cookie_name=$cookie_ID_token;";`  
- Add cookie; `proxy_set_header Cookie "new_cookie=value; $http_cookie";`  
- And every other nginx-directive you can think of