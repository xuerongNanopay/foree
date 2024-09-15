import { View, Text, SafeAreaView, TouchableOpacity } from 'react-native'
import { useLocalSearchParams } from 'expo-router'
import React, { useEffect, useState } from 'react'

import { accountService } from '../../service'
import { formatContactName } from '../../util/contact_util'

const ContactDetail = () => {
  const {contactId} = useLocalSearchParams()
  const [contact, setContact] = useState(null)

  useEffect(() => {
    const controller = new AbortController()
    const getContactDetail = async () => {
      try {
        const resp = await accountService.getContactAccount(contactId, {signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get all active contacts", resp.status, resp.data)
        } else {
          setContact(resp.data.data)
        }
      } catch (e) {
        console.error(e)
        //TODO: route back
      }
    }
    getContactDetail()
    return () => {
      controller.abort()
    }
  }, [])
  return (
    <SafeAreaView>
      <View className="px-2 py-4">
        <View className="flex flex-row items-center pb-4 border-b-[1px] border-slate-300">
          <Text className="flex-1 font-pbold text-lg">{!!contact ? formatContactName(contact, 20) : ""}</Text>
          <TouchableOpacity
            activeOpacity={0.7}
            className="bg-[#1A6B54] py-2 px-4 rounded-full"
            disabled={false}
            onPress={() => {console.log("TODO: route to send")}}
          >
            <Text className="font-pextrabold text-white">Send</Text>
          </TouchableOpacity>
        </View>
      </View>
    </SafeAreaView>
  )
}

export default ContactDetail