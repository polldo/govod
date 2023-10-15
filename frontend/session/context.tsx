import { useState } from 'react'
import { ReactNode } from 'react'
import { useEffect } from 'react'
import { createContext } from 'react'
import { useContext } from 'react'
import { useCallback } from 'react'

type Session = {
    isLoggedIn: boolean
    isLoading: boolean
    updateSession: () => void
}

export const SessionContext = createContext<Session>({
    isLoggedIn: false,
    isLoading: true,
    updateSession: () => {},
})

export function useSession() {
    return useContext(SessionContext)
}

const base = process.env.NEXT_PUBLIC_BASE_URL || 'localhost:8000'

export function SessionProvider(props: { children: ReactNode }) {
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)
    const [isLoading, setIsLoading] = useState<boolean>(true)

    const sync = useCallback(async () => {
        try {
            const resp = await fetch(base + '/users/current', { credentials: 'include' })
            setIsLoggedIn(resp.ok)
            setIsLoading(false)
        } catch {
            setIsLoading(false)
        }
    }, [])

    useEffect(() => {
        sync()
    }, [sync])

    const session: Session = {
        isLoggedIn: isLoggedIn,
        isLoading: isLoading,
        updateSession: sync,
    }

    return <SessionContext.Provider value={session}>{props.children}</SessionContext.Provider>
}
