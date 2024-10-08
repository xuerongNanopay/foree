import { View, Text, SafeAreaView, TouchableOpacity, ScrollView, Alert } from 'react-native'
import { router, useLocalSearchParams } from 'expo-router'
import React, { useEffect, useState } from 'react'

import { accountService } from '../../service'
import { formatName } from '../../util/foree_util'
import Regions from '../../constants/region'
import Countries from '../../constants/country'
import { ContactTransferCashPickup, ContactTransferMethods } from '../../constants/contacts'

const ContactDetail = () => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const {contactId} = useLocalSearchParams()
  const [contact, setContact] = useState(null)

  useEffect(() => {
    const controller = new AbortController()
    const getContactDetail = async () => {
      try {
        const resp = await accountService.getContactAccount(contactId, {signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get contact detail", resp.status, resp.data)
          router.replace('/contact')
        } else {
          setContact(resp.data.data)
        }
      } catch (e) {
        console.error(e)
        router.replace('/contact')
      }
    }
    getContactDetail()
    return () => {
      controller.abort()
    }
  }, [])

  onDelete = () => {
    setIsSubmitting(true)
    Alert.alert("Delete", "Are you sure?",[
      {text: 'Continue', onPress: async () => {
        try {
          const resp = await accountService.deleteContact(contact.id)
          if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
            console.error("delete contact", resp.status, resp.data)
            router.replace('/contact')
          } else {
            router.replace('/contact')
          }
        } catch (e) {
          console.error(e)
        } finally {
          setIsSubmitting(false)
        }
      }},
      {text: 'Cancel', onPress: () => {
        setIsSubmitting(false)
      }}
    ])
  }

  return (
    <SafeAreaView>
      <View className="flex h-full px-2 py-4">
        <View className="flex flex-row items-center pb-4 border-b-[1px] border-slate-300">
          <Text className="flex-1 font-pbold text-lg">{!!contact ? formatName(contact, 20) : ""}</Text>
          <TouchableOpacity
            activeOpacity={0.7}
            className="bg-[#1A6B54] py-2 px-4 rounded-full"
            disabled={false}
            onPress={() => {
              router.push({ pathname: "/transaction/create", params: { contactId } })
            }}
          >
            <Text className="font-pextrabold text-white">Send</Text>
          </TouchableOpacity>
        </View>
        {/* Contact Detail */}
        <ScrollView className="flex-1 mt-4 mb-2">
          {/* Contact Information Details */}
          <View className="bg-slate-200 border-[1px] border-slate-500 rounded-lg py-6 px-2">
            {/* Name */}
            <View className="border-b-[1px] border-slate-400 pb-2">
              <Text className="font-semibold mb-2">
                First Name: <Text className="text-slate-600">{contact?.firstName}</Text>
              </Text>
              { !!contact?.middleName ? <>
                  <Text className="font-semibold mb-2">
                    Middle Name: <Text className="text-slate-600">{contact?.middleName}</Text>
                  </Text>
                </> : null
              }
              <Text className="font-semibold">
                Last Name: <Text className="text-slate-600">{contact?.lastName}</Text>
              </Text>
            </View>
            {/* Address */}
            <View className="border-b-[1px] border-slate-400 pt-2">
              <Text className="font-semibold mb-2">
                Address Line 1: <Text className="text-slate-600">{contact?.address1}</Text>
              </Text>
              { !!contact?.address2 ? <>
                  <Text className="font-semibold mb-2">
                    Address Line 2: <Text className="text-slate-600">{contact?.address2}</Text>
                  </Text>
                </> : null
              }
              <Text className="font-semibold mb-2">
                City: <Text className="text-slate-600">{contact?.city}</Text>
              </Text>
              <Text className="font-semibold mb-2">
                Province: <Text className="text-slate-600">{!!contact?.province ? Regions["PK"]?.[contact.province].name : ""}</Text>
              </Text>
              <Text className="font-semibold mb-2">
                Country: <Text className="text-slate-600">{!!contact?.country ? Countries[contact.country]?.name : ""}</Text>
              </Text>
              { !!contact?.postalCode ? <>
                  <Text className="font-semibold mb-2">
                    Postal Code: <Text className="text-slate-600">{contact?.postalCode}</Text>
                  </Text>
                </> : null
              }
              { !!contact?.phoneNumber ? <>
                  <Text className="font-semibold mb-2">
                    Phone Number: <Text className="text-slate-600">{contact?.phoneNumber}</Text>
                  </Text>
                </> : null
              }
            </View>
            {/* Destination */}
            <View className="mt-2">
              <Text className="font-semibold mb-2">
                Contact Relationship: <Text className="text-slate-600">{contact?.relationshipToContact}</Text>
              </Text>
              <Text className="font-semibold">
                Transfer Method: <Text className="text-slate-600">{!!contact?.transferMethod ? (ContactTransferMethods.find(x => x.value === contact.transferMethod)?.name || contact.transferMethod): ""}</Text>
              </Text>
              {
                !!contact?.transferMethod && contact.transferMethod !== ContactTransferCashPickup ?
                <>
                  <Text className="font-semibold mt-2">
                    Bank Name: <Text className="text-slate-600">{contact?.bankName}</Text>
                  </Text>
                  <Text className="font-semibold mt-2">
                    Account NO: <Text className="text-slate-600">{contact?.accountNoOrIBAN}</Text>
                  </Text>
                </> : null
              }
            </View>
          </View>
          {/* Delete button */}
          <View className="mt-4 flex flex-row">
            <TouchableOpacity
              activeOpacity={0.7}
              className="py-2 px-4 rounded-xl border-[1px] border-red-600 bg-red-100"
              disabled={isSubmitting}
              onPress={onDelete}
            >
              <Text className="font-psemibold text-red-600 text-lg">Delete</Text>
            </TouchableOpacity>
          </View>
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

export default ContactDetail