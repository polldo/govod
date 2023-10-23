import { useSession } from '@/session/context'
import { fetcher } from '@/services/fetch'
import { toast } from 'react-hot-toast'

export default function Logout() {
    const { updateSession } = useSession()

    function handleLogout() {
        fetcher
            .fetch('/auth/logout', { method: 'POST' })
            .then(() => {
                updateSession()
            })
            .catch(() => {
                toast.error('Something went wrong')
            })
    }

    return (
        <button onClick={handleLogout} className="w-full rounded bg-red-600 p-2 font-semibold text-white">
            Logout
        </button>
    )
}
