export type DataResponse = {
    statusCode: number;
    statusMessage: string;
    redirectURI: RedirectURI | null;
    headers: MapHeaders;
    body: any;
};
export type RedirectURI = {
    uri: string;
    path: string;
    queryParams: {
        [key: string]: string;
    }
}

function getRedirectURI(response: Response): RedirectURI | null {
    const location: string | null = response.headers.get("Location");
    if (location === null) {
        return null;
    }
    const fragments = location.split("?");
    const path = fragments[0];
    const queryString = fragments[1] || '';
    const queryParamStrings = queryString.split("&");
    const queryParams = queryParamStrings
        .filter((str) => str.length > 0)
        .map((str) => str.split("="))
        .reduce((acc, [key, value]) => ({...acc, [key]: value}), {});

    return { uri: location, path, queryParams };
}

type MapHeaders = {[key: string]: string | Array<string> };
function getHeadersMap(response: Response): MapHeaders {
    const acc : MapHeaders = {};
    for (const [key, value] of response.headers.entries()) {
        const prevValue = acc[key];
        if (!prevValue) {
            acc[key] = value;
        } else if (Array.isArray(prevValue)) {
            acc[key] = [...prevValue, value];
        } else {
            acc[key] = [prevValue, value];
        }
    }
    return acc;
}

export async function fetchData(url: string, headers: { [key: string]: string } = {}, body?: any): Promise<DataResponse> {
    const response = await fetch(url, {
        method: body ? 'POST' : 'GET',
        headers,
        body: body ? JSON.stringify(body) : undefined,
        redirect: 'manual'
    });
    const responseText = await response.text();
    return {
        statusCode: response.status,
        statusMessage: response.statusText,
        redirectURI: getRedirectURI(response),
        headers: getHeadersMap(response),
        body: responseText
    }
}

export function fetchJson(url: string, headers: { [key: string]: string } = {}, body?: any): Promise<DataResponse> {
    return fetchData(url, headers, body)
        .then(({statusCode, statusMessage, redirectURI, headers,  body}) => {
            if (body) {
                return {statusCode, statusMessage, redirectURI, headers, body: JSON.parse(body)}
            } else {
                throw new Error(`${url} did not return json, body: ${body}`);
            }
        });
}