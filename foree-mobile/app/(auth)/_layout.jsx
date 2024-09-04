import { StatusBar } from 'expo-status-bar'
import { Button, Text, Image } from 'react-native'
import { router, Stack } from 'expo-router'
import React from 'react'
import { images } from '../../constants'

// flow:
// sign in -> sign up -> verify email -> onboarding
// sign in -> sign up -> verify email -b-> sign in
// sign in -> sign up -> verify email -> onboarding -b-> sign in
// foget password -> update password -> sign in
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
          name="login"
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
            headerTitle: props => (<Image source={images.logoSmall} resizeMode='contain' className="w-[24px] h-[24px]" />),
            headerBackTitle:"Back",
          }}
        />
        <Stack.Screen
          name="update_password"
          options={{
            headerShown: true,
            headerTitle: props => (<Image source={images.logoSmall} resizeMode='contain' className="w-[24px] h-[24px]" />),
            headerLeft: () => <Button onPress={() => router.replace("/login")} title="Login" color="#004d40"/>
          }}
        />
        <Stack.Screen
          name="verify_email"
          options={{
            headerShown: true,
            title:"",
            contentStyle:{top:0, bottom:0},
            headerTitle: props => (<Image source={images.logoSmall} resizeMode='contain' className="w-[24px] h-[24px]" />),
            headerLeft: () => <Button onPress={() => router.replace("/login")} title="Logout" color="#004d40"/>
          }}
        />
        <Stack.Screen
          name="onboarding"
          options={{
            headerShown: true,
            title:"",
            headerTintColor:"#004d40",
            headerTitle: props => (<Image source={images.logoSmall} resizeMode='contain' className="w-[24px] h-[24px]" />),
            headerLeft: () => <Button onPress={() => router.replace("/login")} title="Logout" color="#004d40"/>
          }}
        />
      </Stack>
      {/* TODO: investigate why StatusBar.backgroundColor not working. */}
      {/* <StatusBar backgroundColor='#000000' translucent={true} style='auto'/> */}
    </>
  )
}

export default AuthLayout