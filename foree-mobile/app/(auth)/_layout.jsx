import { Button, Text } from 'react-native'
import { router } from 'expo-router'
import { Stack } from 'expo-router'
import React from 'react'

// flow:
// sign in -> sign up -> verify email -> onboarding
// sign in -> sign up -> verify email -b-> sign in
// sign in -> sign up -> verify email -> onboarding -b-> sign in
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
            headerTintColor:"#004d40",
            headerLeft: () => <Button onPress={() => router.replace("sign_in")} title="Sign In" color="#004d40"/>
          }}
        />
        <Stack.Screen
          name="onboarding"
          options={{
            headerShown: true,
            title:"",
            headerTintColor:"#004d40",
            headerLeft: () => <Button onPress={() => router.replace("sign_in")} title="Sign In" color="#004d40"/>
          }}
        />
      </Stack>
      {/* TODO: investigate why StatusBar.backgroundColor not working. */}
      {/* <StatusBar backgroundColor='#004d40' style='light'/> */}
    </>
  )
}

export default AuthLayout