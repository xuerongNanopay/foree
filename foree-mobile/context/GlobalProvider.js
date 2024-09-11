import { createContext, useContext, useState, useEffect } from "react"

import { authService } from "../service"

const GlobalContext = createContext()

export const useGlobalContext = () => useContext(GlobalContext)

export const GlobalProvider = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [user, setUser] = useState(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const getUser = async () => {
      setIsLoading(true)
      try {
        const resp = await authService.getUser()
        if ( resp.status / 100 !== 2 ) {
          setIsLoggedIn(false)
          setUser(null)
        } else {
          setIsLoggedIn(true)
          setUser(resp.data.data)
        }
      } catch (e) {
        setIsLoggedIn(false)
      } finally {
        setIsLoading(false)
      }

    }
    getUser()
  }, [])

  return (
    <GlobalContext.Provider
      value={{
        isLoggedIn,
        setIsLoading,
        user,
        setUser,
        isLoading
      }}
    >
      {children}
    </GlobalContext.Provider>
  )
}