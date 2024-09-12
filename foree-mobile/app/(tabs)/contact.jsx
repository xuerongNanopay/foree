import { View, Text, TouchableOpacity } from 'react-native'
import { useFocusEffect } from 'expo-router'
import React, { useEffect, useState, useCallback } from 'react'

import { SafeAreaView } from 'react-native-safe-area-context'
import SearchInput from '../../components/SearchInput'

const Contact = () => {

  useFocusEffect(
    useCallback(() => {
      console.log('Hello, I am focused!');

      return () => {
        console.log('This route is now unfocused.');
      }
    }, [])
  )

  return (
    <SafeAreaView className="border-2 border-red-600">
      <View className="px-4 pt-4">
        <View className="pb-4 border-b-[1px] border-slate-300">
          <View className="flex flex-row items-center">
            <Text className="flex-1 font-pbold text-2xl">Contacts</Text>
            <TouchableOpacity
              onPress={()=> {console.log("TODO: contact creation")}}
              activeOpacity={0.7}
              className="bg-[#1A6B54] py-2 px-4 rounded-full"
              disabled={false}
            >
              <Text className="font-pbold text-white">Add</Text>
            </TouchableOpacity>
          </View>
          <View>
            <SearchInput
              placeholder="search name..."
              variant='bordered'
              containerStyles="mt-4"
            />
          </View>
        </View>
      </View>
    </SafeAreaView>
  )
}

export default Contact