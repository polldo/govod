import { useState } from 'react'
import { ReactNode } from 'react'
import { useEffect } from 'react'
import { createContext } from 'react'
import { useContext } from 'react'
import Cookies from 'js-cookie'

type Session = {
    isLoggedIn: boolean
    isLoading: boolean
    login: () => void
    logout: () => void
}

export const SessionContext = createContext<Session>({
    isLoggedIn: false,
    isLoading: true,
    login: () => {},
    logout: () => {},
})

export function useSession() {
    return useContext(SessionContext)
}

export function SessionProvider(props: { children: ReactNode }) {
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)
    const [isLoading, setIsLoading] = useState<boolean>(true)

    // When the user wants to logout or the session expires, this function must be called.
    function logout() {
        Cookies.remove('session')
        setIsLoggedIn(false)
    }

    // Upon login, the state is automatically added in cookies,
    // so we only need to set the session state to true.
    function login() {
        setIsLoggedIn(true)
    }

    useEffect(() => {
        const user = async () => {
            const response = await fetch('http://mylocal.com:8000/users', { credentials: 'include' })
            setIsLoggedIn(response.ok)
            setIsLoading(false)
        }

        user().catch(() => {
            setIsLoggedIn(false)
            setIsLoading(false)
        })
    }, [])

    const session: Session = {
        isLoggedIn: isLoggedIn,
        isLoading: isLoading,
        login: login,
        logout: logout,
    }

    return <SessionContext.Provider value={session}>{props.children}</SessionContext.Provider>
}
