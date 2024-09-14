import { View, Text } from 'react-native'
import { useLocalSearchParams } from 'expo-router'

import React from 'react'

const ContactDetail = () => {
  const {contactId} = useLocalSearchParams()
  console.log(contactId)
  return (
    <View>
      <Text>ContactDetail</Text>
    </View>
  )
}

export default ContactDetail