import { useSession } from '@/session/context'
import { useFetch } from '@/services/fetch'

export default function Logout() {
    const { updateSession } = useSession()
    const fetch = useFetch()

    function handleLogout() {
        fetch('http://mylocal.com:8000/auth/logout', { method: 'POST', credentials: 'include' })
            .then((response) => {
                if (!response.ok) {
                    throw new Error()
                }
                updateSession()
            })
            .catch((err) => {
                console.log(err)
            })
    }

    return (
        <button onClick={handleLogout} className="w-full rounded bg-red-500 p-2 font-semibold text-white">
            Logout
        </button>
    )
}
