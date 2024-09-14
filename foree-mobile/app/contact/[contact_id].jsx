import { View, Text } from 'react-native'
import { useLocalSearchParams } from 'expo-router'

import React from 'react'

const ContactDetail = () => {
  const {contact_id} = useLocalSearchParams()
  return (
    <View>
      <Text>ContactDetail</Text>
    </View>
  )
}

export default ContactDetail