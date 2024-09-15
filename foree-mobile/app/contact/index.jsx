import { View, Text } from 'react-native'
import React from 'react'
import { Redirect } from 'expo-router'

const ContracIndex = () => {
  return (
    <View>
      <Redirect href="/contact_tab"/>
    </View>
  )
}

export default ContracIndex