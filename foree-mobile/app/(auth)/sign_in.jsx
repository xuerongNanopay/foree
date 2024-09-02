import { ScrollView, StyleSheet, Text, View, Image } from 'react-native'
import React from 'react'
import { SafeAreaView } from 'react-native-safe-area-context'

import { images } from '../../constants'

const SignIn = () => {
  return (
    <SafeAreaView className="h-full">
      <ScrollView>
        <View class="w-full justify-center min-h-[85vh] px-2 my-6">
          <Image
            source={images.logo}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default SignIn