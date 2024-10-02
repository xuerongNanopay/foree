import { View, Text, Button } from 'react-native'
import React from 'react'
import { router, Stack } from 'expo-router'

const SettingLayout = () => {
  return (
  <Stack
    screenOptions={{
      headerTintColor:"#004d40",
      headerTitleStyle: {
        fontWeight: 'bold',
      },
    }}
  >
    <Stack.Screen
      name="update_passwd"
      options={{
        headerShown: true,
        title:"Update Password",
        headerTintColor:"#004d40",
        headerLeft: () => <Button onPress={() => {
          if ( router.canGoBack() ) {
            router.back()
          } else {
            router.replace("/settings_tab")
          }
        }} title="Back" color="#004d40"/>
      }}
    />
    <Stack.Screen
      name="profile"
      options={{
        headerShown: true,
        title:"Profile",
        headerTintColor:"#004d40",
        headerLeft: () => <Button onPress={() => {
          if ( router.canGoBack() ) {
            router.back()
          } else {
            router.replace("/settings_tab")
          }
        }} title="Back" color="#004d40"/>
      }}
    />
    <Stack.Screen
      name="notification_settings"
      options={{
        headerShown: true,
        title:"Notication",
        headerTintColor:"#004d40",
        headerLeft: () => <Button onPress={() => {
          if ( router.canGoBack() ) {
            router.back()
          } else {
            router.replace("/settings_tab")
          }
        }} title="Back" color="#004d40"/>
      }}
    />
  </Stack>
  )
}

export default SettingLayout