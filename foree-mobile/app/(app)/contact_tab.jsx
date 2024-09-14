import { View, Text, TouchableOpacity } from 'react-native'
import { router, useFocusEffect } from 'expo-router'
import React, { useEffect, useState, useCallback } from 'react'

import { SafeAreaView } from 'react-native-safe-area-context'
import SearchInput from '../../components/SearchInput'
import { accountService } from '../../service'

const ContactTab = () => {

  const [searchText, setSearchText] = useState("")
  useFocusEffect(useCallback(() => {
    const controller = new AbortController()
    const getAllContacts = async() => {
      try {
        const resp = await accountService.getAllContactAccounts()
        console.log(resp.data)
      } catch (e) {
        console.error(e)
      }
    }
    getAllContacts()
    return () => {
      controller.abort()
    }
  }, []))

  return (
    <SafeAreaView className="border-2 border-red-600">
      <View className="px-4 pt-4">
        <View className="pb-4 border-b-[1px] border-slate-300">
          <View className="flex flex-row items-center">
            <Text className="flex-1 font-pbold text-2xl">Contacts</Text>
            <TouchableOpacity
              onPress={()=> {router.push("/contact/create")}}
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
              value={searchText}
              handleChangeText={(e) => setSearchText(e.toLowerCase())}
              containerStyles="mt-4"
            />
          </View>
        </View>
      </View>
    </SafeAreaView>
  )
}

const ContactListItem = ({
  contact
}) => {
  if ( ! contact ) return <></>

}

export default ContactTab