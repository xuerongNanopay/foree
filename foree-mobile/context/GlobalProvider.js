import { createContext, useContext, useState, useEffect, Children } from "react"

const GlobalContext = createContext();

export const useGlobalContext = () => useContext(GlobalContext)

const GlobalProvider = ({ children }) => {
  return (
    <GlobalContext.Provider>
      {children}
    </GlobalContext.Provider>
  )
}