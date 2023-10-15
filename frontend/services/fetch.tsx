export class ResponseError extends Error {
    status: number
    constructor(status: number, message: string) {
        super(message)
        this.status = status
    }
}

// Fetcher wraps the fetch function in order to:
// - add the 'credentials' header, to allow sending the session cookie to the backend.
// - intercept 401 errors and execute a custom callback.
// - throw ResponseError when fetched status code is not OK.
export class Fetcher {
    onUnauthenticated: () => void
    f: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>

    constructor(
        onUnauth: () => void,
        fetchFunction: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>
    ) {
        this.onUnauthenticated = onUnauth
        this.f = fetchFunction
    }

    setUnauthenticated(onUnauth: () => void) {
        this.onUnauthenticated = onUnauth
    }

    async fetch(url: RequestInfo, options: RequestInit = {}): Promise<Response> {
        if (options.method !== 'OPTIONS') {
            options.credentials = 'include'
        }

        const response = await this.f(url, options)
        if (response.status === 401) {
            this.onUnauthenticated()
            const res = await response.json()
            throw new ResponseError(response.status, res.error)
        }

        if (!response.ok) {
            const res = await response.json()
            throw new ResponseError(response.status, res.error)
        }

        return response
    }
}

export const fetcher = new Fetcher(
    () => {
        console.log('unauthenticated')
    },
    (...args) => fetch(...args)
)
