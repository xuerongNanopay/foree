import { View, Text } from 'react-native'
import { useLocalSearchParams } from 'expo-router'
import React, { useEffect } from 'react'
import { accountService } from '../../service'

const ContactDetail = () => {
  const {contactId} = useLocalSearchParams()

  useEffect(() => {
    const controller = new AbortController()
    const getContactDetail = async () => {
      try {
        const resp = await accountService.getContactAccount(contactId, {signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get all active contacts", resp.status, resp.data)
        } else {
          //How do this: because there is cache in getAllContactAccounts
          //TODO: redesign the cache?
          console.log(resp.data.data)
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
    <View>
      <Text>ContactDetail</Text>
    </View>
  )
}

export default ContactDetail