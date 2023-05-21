import { useState } from 'react'
import { ReactNode } from 'react'
import { useEffect } from 'react'
import { createContext } from 'react'
import { useContext } from 'react'

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

export function SessionProvider(props: { children: ReactNode }) {
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)
    const [isLoading, setIsLoading] = useState<boolean>(true)

    function sync() {
        const user = async () => {
            const response = await fetch('http://mylocal.com:8000/users/current', { credentials: 'include' })
            setIsLoggedIn(response.ok)
            setIsLoading(false)
        }

        user().catch(() => {
            setIsLoading(false)
        })
    }

    useEffect(() => {
        sync()
    }, [])

    const session: Session = {
        isLoggedIn: isLoggedIn,
        isLoading: isLoading,
        updateSession: sync,
    }

    return <SessionContext.Provider value={session}>{props.children}</SessionContext.Provider>
}
