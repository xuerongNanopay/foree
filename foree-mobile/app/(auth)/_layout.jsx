import { StatusBar } from 'expo-status-bar'
import { Stack } from 'expo-router'
import React from 'react'

const AuthLayout = () => {
  return (
    <>
      <Stack>
        <Stack.Screen
          name="sign_in"
          options={{
            headerShown: false
          }}
        />
        <Stack.Screen
          name="sign_up"
          options={{
            headerShown: false
          }}
        />
      </Stack>
      <StatusBar backgroundColor='#004d40' style='light'/>
    </>
  )
}

export default AuthLayout