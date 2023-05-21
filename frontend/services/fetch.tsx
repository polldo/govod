import { useSession } from '@/session/context'

// useFetch returns a wrapped fetch function that:
// - adds the 'credentials' header, to allow sending the session cookie to the backend.
// - calls the logout function when a 401 error is returned.
export function useFetch() {
    const { updateSession } = useSession()

    async function customFetch(url: RequestInfo, options: RequestInit = {}): Promise<Response> {
        if (options.method !== 'OPTIONS') {
            options.credentials = 'include'
        }

        const response = await fetch(url, options)
        if (response.status === 401) {
            updateSession()
        }

        return response
    }

    return customFetch
}
