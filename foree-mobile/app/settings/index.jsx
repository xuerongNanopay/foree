import { View, Text } from 'react-native'
import React from 'react'
import { Redirect } from 'expo-router'

const SettingsTab = () => {
  return (
    <View>
      <Redirect href="/settings_tab"/>
    </View>
  )
}

export default SettingsTab