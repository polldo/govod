import { useState } from 'react'
import { ReactNode } from 'react'
import { useEffect } from 'react'
import { createContext } from 'react'
import Cookies from 'js-cookie'

type Session = {
    isLoggedIn: boolean
    login: () => void
    logout: () => void
}

// export const SessionContext = createContext<Session | undefined>(undefined)
export const SessionContext = createContext<Session>({ isLoggedIn: false, login: () => {}, logout: () => {} })

export function SessionProvider(props: { children: ReactNode }) {
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)

    // At startup check if there is a session cookie.
    // If it is, then let's consider the user logged in.
    useEffect(() => {
        const session = Cookies.get('session')
        setIsLoggedIn(!!session)
    }, [])

    // When the user wants to logout or the session expires, this function must be called.
    const logout = () => {
        Cookies.remove('session')
        setIsLoggedIn(false)
    }

    // Upon login, the state is automatically added in cookies,
    // so we only need to set the session state to true.
    const login = () => {
        setIsLoggedIn(true)
    }

    const session: Session = {
        isLoggedIn: isLoggedIn,
        login: login,
        logout: logout,
    }

    return <SessionContext.Provider value={session}>{props.children}</SessionContext.Provider>
}
