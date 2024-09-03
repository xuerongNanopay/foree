import { StatusBar } from 'expo-status-bar'
import { Stack } from 'expo-router'
import React from 'react'

const AuthLayout = () => {
  return (
    <>
      <Stack
        screenOptions={{
          headerTintColor:"#004d40",
          headerTitleStyle: {
            fontWeight: 'bold',
          },
        }}
      >
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
        <Stack.Screen
          name="forget_password"
          options={{
            headerShown: true,
            title:"",
            headerBackTitle:"Sign In",
          }}
        />
        <Stack.Screen
          name="verify_email"
          options={{
            headerShown: true,
            title:"",
            headerBackTitle:"Sign In",
          }}
        />
        <Stack.Screen
          name="onboarding"
          options={{
            headerShown: true,
            title:"",
            headerBackTitle:"Sign In",
          }}
        />
      </Stack>
      {/* TODO: investigate why StatusBar.backgroundColor not working. */}
      {/* <StatusBar backgroundColor='#004d40' style='light'/> */}
    </>
  )
}

export default AuthLayout