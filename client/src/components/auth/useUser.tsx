import useSWR from 'swr'
import { fetcherWithToken } from 'libs/fetch'
import cookie from 'js-cookie'

export default function useUser() {
  const { data, mutate, error } = useSWR(
    `${process.env.NEXT_PUBLIC_API_URL}/user`,
    fetcherWithToken,
    {
      initialData: undefined,
    }
  )

  const loading = !data && !error
  const loggedOut = data?.errors?.body === 'Unauthorized'

  if (!loggedOut) {
    cookie.set('auth', 'true', { expires: 1 })
  } else {
    cookie.remove('auth')
  }

  return {
    loading,
    loggedOut,
    user: loggedOut ? null : data,
    mutate,
  }
}
