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
      name="personal_settings"
      options={{
        headerShown: true,
        title:"Personal Settings",
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
      name="update_address"
      options={{
        headerShown: true,
        title:"Update Address",
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
      name="update_phone_number"
      options={{
        headerShown: true,
        title:"Update Phone Number",
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
        title:"Notication Settings",
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
      name="invitation"
      options={{
        headerShown: true,
        title:"Invitation",
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
      name="close_account"
      options={{
        headerShown: true,
        title:"Close Account",
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