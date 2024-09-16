import { View, Text, TouchableOpacity, ScrollView } from 'react-native'
import { router, useFocusEffect } from 'expo-router'
import React, { useEffect, useState, useCallback } from 'react'
import { SafeAreaView } from 'react-native'

import SearchInput from '../../components/SearchInput'
import { accountService } from '../../service'
import { ContactTransferCashPickup } from '../../constants/contacts'
import string_util from '../../util/string_util'
import { formatName } from '../../util/foree_util'

const ContactTab = () => {
  const [searchText, setSearchText] = useState("")
  const [contacts, setContacts] = useState([])
  const [showContacts, setShowContacts] = useState([])

  useFocusEffect(useCallback(() => {
    const  controller = new AbortController()
    const getAllContacts = async (signal) => {
      try {
        const resp = await accountService.getAllContactAccounts({signal})

        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get all active contacts", resp.status, resp.data)
        } else {
          //How do this: because there is cache in getAllContactAccounts
          //TODO: redesign the cache?
          setContacts([...resp.data.data])
          setSearchText("")
        }
      } catch (e) {
        console.error(e)
      }
      return controller
    }
    getAllContacts(controller.signal)
    return () => {
      controller.abort()
    }
  }, []))

  useEffect(() => {
    if ( !searchText ) setShowContacts(contacts)
    else {
      setShowContacts(contacts.filter(c => string_util.containSubsequence(`${c.firstName}${c.middleName ?? ""}${c.lastName}`, searchText, {caseInsensitive:true})))
    }
  }, [searchText, contacts])

  return (
    <SafeAreaView className="">
      <View className="flex h-full px-4 pt-4">
        <View className="pb-4 mb-4 border-b-[1px] border-slate-300">
          <View className="flex flex-row items-center">
            <Text className="flex-1 font-pbold text-2xl">Contacts</Text>
            <TouchableOpacity
              onPress={()=> {router.push("/contact/create")}}
              activeOpacity={0.7}
              className="bg-[#1A6B54] py-2 px-4 rounded-full"
              disabled={false}
            >
              <Text className="font-pextrabold text-white">Add</Text>
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
        <ScrollView 
          className="flex-1 mb-2"
          showsVerticalScrollIndicator={false}
        >
          {
            showContacts.map(contact => {
              return (
                <TouchableOpacity 
                  key={contact.id}
                  onPress={() => router.push(`/contact/${contact.id}`)}
                  activeOpacity={0.7}
                >
                  <ContactListItem contact={contact}/>
                </TouchableOpacity>
              )
            })
          }
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

const ContactListItem = ({
  contact
}) => {
  if ( ! contact ) return <></>

  return (
    <View className="mb-2 p-2 rounded-lg bg-[#ccded6]">
      <Text className="font-bold">{formatName(contact)}</Text>
      <FormatContactTransferInfo {...contact}/>
      <FormatContactTransferRecentActivity {...contact}/>
    </View>
  )
}


const FormatContactTransferInfo = ({transferMethod, bankName, accountNoOrIBAN}) => {
  if ( transferMethod === ContactTransferCashPickup ) 
    return <Text className="font-semibold text-slate-700">Cash Pickup</Text>
  return (
    <Text className="font-semibold text-slate-700">
      {!!bankName ? bankName.slice(0, 14) + (bankName.length > 14 ? "..." : "") : ""}
      <Text className="italic">
        ({!!accountNoOrIBAN ? accountNoOrIBAN.slice(0, 7) + (accountNoOrIBAN.length > 7 ? "..." : "") : ""})
      </Text>
    </Text>
  )
}

const FormatContactTransferRecentActivity = ({latestActiveAt}) => {
  if ( !latestActiveAt ) return <Text className="text-slate-600 italic">Last sent: -</Text>
  return <Text className="text-slate-600 italic">Last sent: {
    new Date(latestActiveAt).toLocaleString()
  }</Text>
}
export default ContactTab