import { Button, Text, Image, View } from 'react-native'
import { router, Stack } from 'expo-router'
import React from 'react'
import { images } from '../../constants'

const ContactLayout = () => {
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
        name="create"
        options={{
          headerShown: true,
          title:"Create Contact",
          headerTintColor:"#004d40",
          headerLeft: () => <Button onPress={() => {
            if ( router.canGoBack() ) {
              router.back()
            } else {
              router.replace("/contact_tab")
            }
          }} title="Back" color="#004d40"/>
        }}
      />
      </Stack>
    </>
  )
}

export default ContactLayout