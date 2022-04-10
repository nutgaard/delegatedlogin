# Login-app

login-app is a standalone application that implements delegated login for your users.
Authentization happens through an external identity provider (idp), configureted by the `IDP_DISCOVERY_URL` property.

## Konfigurasjon
| Name                   | Required | Description                                                                               | Default         |
|------------------------|----------|-------------------------------------------------------------------------------------------|-----------------|
| APP_NAME               | Yes      | Name of application. Used as the app's context-path.                                      |                 |
| APP_VERSION            | Yes      | Version of application. Used on selftest page.                                            |                 |
| IDP_DISCOVERY_URL      | Yes      | Url to discovery-url for idp (e.g something ending with .well-known/openid-configuration) |                 |
| IDP_CLIENT_ID          | Yes      | Username for usage against IDP                                                            |                 |
| IDP_CLIENT_SECRET      | Yes      | Password for usage against IDP                                                            |                 |
| AUTH_TOKEN_RESOLVER    | No       | Cookie name for your id token                                                             | `ID_token`      |
| REFRESH_TOKEN_RESOLVER | No       | Cookie name for your refresh token                                                        | `refresh_token` |