import {assert, assertEquals, assertExists} from "https://deno.land/std@0.122.0/testing/asserts.ts";
import {fetchData, fetchJson} from "./http-fetch.ts";
import {setup, retry} from "./setup.ts";

await setup('oidc-stub is running', retry({ retry: 10, interval: 2}, async () => {
    const oidcConfig = await fetchJson('http://localhost:8080/.well-known/openid-configuration');
    assertEquals(oidcConfig.statusCode, 200, 'oidcConfig is running');
}));

await setup('loginapp is running', retry({ retry: 10, interval: 2}, async () => {
    const loginapp = await fetchData('http://localhost:8082/loginapp/internal/isAlive');
    assertEquals(loginapp.statusCode, 200, 'loginapp is running');
}));

await setup('frontendapp is running', retry({ retry: 10, interval: 2}, async () => {
    const frontendapp = await fetchData('http://localhost:8083/frontend/internal/isAlive');
    assertEquals(frontendapp.statusCode, 200, 'frontendapp is running');
}));

Deno.test("oidc-stub provides jwks", async () => {
    const jwks = await fetchJson('http://localhost:8080/.well-known/jwks.json');
    assertExists(jwks.body, "jwks.json returns a json body");
    assertEquals(jwks.body.keys.length, 1, "jwks has one key");
});

/**
 * Trailing slash is handles by default. No redirect required
 * Test should be rewritten to showcase loading of html with/without trailing slash
 */
// Deno.test('attempts to get frontend resource without trailing slash', async () => {
//     const initial = await fetchData('http://localhost:8083/frontend');
//     assertEquals(initial.statusCode, 301, '/frontend returns 301');
//     assertEquals(initial.redirectURI?.path, 'frontend/', 'appends trailing slash');
// });

Deno.test('attempts to get frontend resource should result in login-flow', async () => {
    const initial = await fetchData('http://localhost:8083/frontend/');

    assertEquals(initial.statusCode, 302, '/frontend returns 302');
    assertEquals(
        initial.redirectURI?.path,
        'http://localhost:8082/loginapp/api/start',
        '/frontend redirects to /loginapp/api/start'
    );
    assertEquals(
        initial.redirectURI?.queryParams?.url,
        encodeURIComponent('http://localhost:8083/frontend/'),
        '/frontend redirect passes original url encoded in queryparameter'
    );

    const startLogin = await fetchData(initial.redirectURI?.uri!!);

    assertEquals(initial.statusCode, 302, '/loginapp/api/start returns 302');
    assertEquals(
        startLogin.redirectURI?.path,
        'http://localhost:8080/authorize',
        '/loginapp/api/start redirects to oidc-stub/authorize'
    );

    const state = startLogin.redirectURI?.queryParams?.state;
    const stateCookie = asArray(startLogin.headers['set-cookie']);
    assertExists(state, '/loginapp/api/start state query-param is present')
    assertEquals(stateCookie.length, 1, '/loginapp/api/start should set state-cookie')
    assertEquals(
        startLogin.redirectURI?.queryParams,
        {
            session: 'winssochain',
            authIndexType: 'service',
            authIndexValue: 'winssochain',
            response_type: 'code',
            scope: 'openid',
            client_id: 'foo',
            state,
            redirect_uri: encodeURIComponent('http://localhost:8082/loginapp/api/login'),
        },
        '/loginapp/api/start passes correct queryParams to idp'
    )

    const authorize = await fetchData(startLogin.redirectURI?.uri!!);
    assertEquals(authorize.statusCode, 302, '/oidc-stub/authorize returns 302');
    assertEquals(
        authorize.redirectURI?.path,
        'http://localhost:8082/loginapp/api/login',
        '/oidc-stub/authorize redirects to loginapp/login'
    );
    const code = authorize.redirectURI?.queryParams?.code;
    assertExists(code, '/oidc-stub/authorize code query-param is present');
    assertEquals(
        authorize.redirectURI?.queryParams?.state,
        state,
        '/oidc-stub/authorize state query-param matches state sent in from loginapp'
    );

    const login = await fetchData(authorize.redirectURI?.uri!!, {
        'Cookie': stateCookie[0]
    });
    assertEquals(login.statusCode, 302, '/loginapp/api/login returns 302');
    assertEquals(
        login.redirectURI?.path,
        'http://localhost:8083/frontend/',
        '/loginapp/login redirects to /frontend'
    );
    const loginCookies = login.headers['set-cookie'];
    const idtoken = asArray(loginCookies).find((cookie: string) => cookie.startsWith('loginapp_ID_token'));
    const refreshtoken = asArray(loginCookies).find((cookie: string) => cookie.startsWith('loginapp_refresh_token'));
    const removeStateCookie = asArray(loginCookies).find((cookie: string) => cookie.startsWith(state));

    assert(idtoken?.startsWith('loginapp_ID_token'), '/loginapp/api/login sets loginapp_ID_token cookie');
    assert(idtoken!!.length > 80, 'loginapp_ID_token has some content');
    assert(idtoken?.includes("Max-Age=3600"), 'loginapp_ID_token is valid for 1 hour')

    assert(refreshtoken?.startsWith('loginapp_refresh_token'), '/loginapp/api/login sets loginapp_refresh_token cookie');
    assert(refreshtoken!!.length > 80, 'loginapp_refresh_token has some content');
    assert(refreshtoken?.includes("Max-Age=72000;"), 'loginapp_ID_token is valid for 24 hours');

    assert(removeStateCookie?.startsWith(state), '/loginapp/api/login sets loginapp_ID_token cookie');
    assert(removeStateCookie?.includes('01 Jan 1970'), '/loginapp/api/login removes state cookie');

    const pageLoadAfterLogin = await fetchData('http://localhost:8083/frontend/', {
        'Cookie': idtoken!!
    });
    assertEquals(pageLoadAfterLogin.statusCode, 200, '/frontend returns 200');
    assert(pageLoadAfterLogin.body.includes('<!DOCTYPE html>'), '/frontend returns HTML')
});

Deno.test('static resources returns 302 login redirect, if not logged in', async () => {
    const staticResource = await fetchData('http://localhost:8083/frontend/static/css/index.css');
    assertEquals(staticResource.statusCode, 302, '/frontend returns 302');
    assert(!staticResource.body.includes('<!DOCTYPE html>'), 'css-file is not HTML')
});

Deno.test('static resources returns 200 ok if logged in', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const staticResource = await fetchData('http://localhost:8083/frontend/static/css/index.css', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(staticResource.statusCode, 200, '/frontend returns 302');
    assertEquals(staticResource.headers['referrer-policy'], 'no-referrer', '/frontend has referrer-policy');
    assert(!staticResource.body.includes('<!DOCTYPE html>'), 'css-file is not HTML')
});

Deno.test('frontend routing should return index.html', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const root = await fetchData('http://localhost:8083/frontend', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    const rootTrailingSlash = await fetchData('http://localhost:8083/frontend/', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    const path = await fetchData('http://localhost:8083/frontend/with/longer/path', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    const pathTrailingSlash = await fetchData('http://localhost:8083/frontend/with/longer/path/', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });

    assertEquals(root.statusCode, 200, '/frontend returns 200');
    assert(root.body.includes('<!DOCTYPE html>'), '/frontend returns HTML');
    assertEquals(rootTrailingSlash.statusCode, 200, '/frontend/ returns 200');
    assert(rootTrailingSlash.body.includes('<!DOCTYPE html>'), '/frontend/ returns HTML');
    assertEquals(path.statusCode, 200, '/frontend/with/longer/path returns 200');
    assert(path.body.includes('<!DOCTYPE html>'), '/frontend/with/longer/path returns HTML');
    assertEquals(pathTrailingSlash.statusCode, 200, '/frontend/with/longer/path/ returns 200');
    assert(pathTrailingSlash.body.includes('<!DOCTYPE html>'), '/frontend/with/longer/path/ returns HTML');
});

Deno.test('frontend routing should return 302 if not logged in', async () => {
    const root = await fetchData('http://localhost:8083/frontend');
    const rootTrailingSlash = await fetchData('http://localhost:8083/frontend/');
    const path = await fetchData('http://localhost:8083/frontend/with/longer/path');
    const pathTrailingSlash = await fetchData('http://localhost:8083/frontend/with/longer/path/');

    assertEquals(root.statusCode, 302, '/frontend returns 302');
    assertEquals(rootTrailingSlash.statusCode, 302, '/frontend returns 302');
    assertEquals(path.statusCode, 302, '/frontend returns 302');
    assertEquals(pathTrailingSlash.statusCode, 302, '/frontend returns 302');
});

Deno.test('missing static resource returns 404 instead of fallback to index.html', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const staticResource = await fetchData('http://localhost:8083/frontend/static/css/missing.css',{
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(staticResource.statusCode, 404, '/frontend returns 404');
    assert(!staticResource.body.includes('<!DOCTYPE html>'), 'css-file is not HTML')
});

/**
 * Feature no longer supported.
 * All request that are proxyied must be by an authenticated user
 */
// Deno.test('proxying to open endpoint when not logged in', async () => {
//     const openEndpointWithoutCookie = await fetchJson('http://localhost:8083/frontend/proxy/open-endpoint/data');
//     assertEquals(openEndpointWithoutCookie.statusCode, 200, '/frontend proxied to open endpoint');
//     assertEquals(openEndpointWithoutCookie.body.path, '/data', '/frontend removed url prefix');
//     assert(openEndpointWithoutCookie.body.headers['cookie'] === undefined, '/frontend did not send cookie');
// });

Deno.test('proxying to open endpoint when logged in', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const openEndpointWithCookie = await fetchJson('http://localhost:8083/frontend/proxy/open-endpoint/data', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(openEndpointWithCookie.statusCode, 200, '/frontend proxied to open endpoint');
    assertEquals(openEndpointWithCookie.body.path, '/data', '/frontend removed url prefix');
    assertExists(openEndpointWithCookie.body.headers['cookie'], '/frontend did send cookie');
});

Deno.test('proxying to protected endpoint when not logged in', async () => {
    const protectedEndpoint = await fetchData('http://localhost:8083/frontend/proxy/protected-endpoint/data');
    assertEquals(protectedEndpoint.statusCode, 302, '/frontend returns 302');
    assertEquals(
        protectedEndpoint.redirectURI?.path,
        'http://localhost:8082/loginapp/api/start',
        '/frontend redirects to /loginapp/api/start'
    );
});

Deno.test('proxying to protected endpoint when logged in', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const protectedEndpoint = await fetchJson('http://localhost:8083/frontend/proxy/protected-endpoint/data', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(protectedEndpoint.statusCode, 200, '/frontend returns 200');
    assertEquals(protectedEndpoint.body.path, '/data', '/frontend removed url prefix');
    assertExists(protectedEndpoint.body.headers['cookie'], '/frontend did send cookie');
    assert(protectedEndpoint.body.headers['cookie'].startsWith('loginapp_ID_token'), '/frontend sent loginapp_ID_token cookie');
});

Deno.test('proxying to open endpoint that removes cookie when logged in', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const openEndpointWithCookie = await fetchJson('http://localhost:8083/frontend/proxy/open-endpoint-no-cookie/data', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(openEndpointWithCookie.statusCode, 200, '/frontend proxied to open endpoint');
    assertEquals(openEndpointWithCookie.body.path, '/data', '/frontend removed url prefix');
    assert(openEndpointWithCookie.body.headers['cookie'] === undefined, '/frontend did not send cookie');
});

Deno.test('proxying to protected endpoint when logged in, and rewriting cookie name', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const protectedEndpoint = await fetchJson('http://localhost:8083/frontend/proxy/protected-endpoint-with-cookie-rewrite/data', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(protectedEndpoint.statusCode, 200, '/frontend returns 200');
    assertEquals(protectedEndpoint.body.path, '/data', '/frontend removed url prefix');
    assert(protectedEndpoint.body.headers['cookie'].startsWith('ID_token'), '/frontend sent ID_token cookie');
    assert(!protectedEndpoint.body.headers['cookie'].includes('loginapp_ID_token'), '/frontend did not send loginapp_ID_token cookie');
});

Deno.test('environments variables are injected into nginx config', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const page = await fetchData('http://localhost:8083/frontend/env-data', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assertEquals(page.body, 'APP_NAME: frontend', 'Page contains environmentvariable value')
});

Deno.test('environments variables are injected into html config', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const page = await fetchData('http://localhost:8083/frontend/', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });
    assert(page.body.includes('&#36;env{APP_NAME}: frontend'), 'Page contains environmentvariable value')
});

Deno.test('csp directive is added to request', async () => {
    const tokens = await fetchJson('http://localhost:8080/oauth/token', {}, {});
    const page = await fetchData('http://localhost:8083/frontend/', {
        'Cookie': `loginapp_ID_token=${tokens.body['id_token']};`
    });

    const cspPolicy = page.headers['content-security-policy-report-only'];
    assertExists(cspPolicy, '/frontend has report-only CSP-policy');
    assert(cspPolicy.includes('script-src'), '/frontend has report-only CSP-policy');
});

function asArray<T>(t: T | Array<T>): Array<T> {
    return Array.isArray(t) ? t : [t];
}