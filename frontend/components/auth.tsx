import type { PropsWithChildren } from 'react'
import { useSession } from '@/session/context'
import { useRouter } from 'next/router'

export default function WithAuth(props: PropsWithChildren) {
    const { isLoading, isLoggedIn } = useSession()
    const router = useRouter()

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('/login')
        return null
    }

    return props
}
