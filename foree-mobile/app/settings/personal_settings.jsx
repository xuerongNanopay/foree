import { View, Text, SafeAreaView, ScrollView, Image, TouchableOpacity } from 'react-native'
import React, { useCallback, useState } from 'react'
import { icons } from '../../constants'
import { router, useFocusEffect } from 'expo-router'
import { authService } from '../../service'
import Countries from '../../constants/country'
import Regions from '../../constants/region'

const profile = () => {
  const [ userDetail, setUserDetail ] = useState(null)

  useFocusEffect(useCallback(() => {
    const controller = new AbortController()
    const getUserDetail = async() => {
      try {
        const resp = await authService.getUserDetail({signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get userDetail", resp.status, resp.data)
          router.replace("/settings_tab")
        } else {
          setUserDetail(resp.data.data)
        }
      } catch (e) {
        console.error("get userDetail", e, e.response, e.response?.status, e.response?.data)
        router.replace("/settings_tab")
      }
    }
    getUserDetail()
    return () => {
      controller.abort()
    }
  }, []))

  if ( !userDetail ) return<></>
  return (
    <SafeAreaView className="h-full bg-slate-200">
      <ScrollView
        className="px-2"
      >
        <View
          className="mt-6 py-3 border-[1px] border-slate-600 rounded-md bg-white shadow-2xl"
        >
          <View
            className="px-4 pb-2 flex flex-row items-center justify-between border-b-[1px] border-slate-400"
          >
            <Text className="font-psemibold text-lg">Full Name</Text>
            {/* <Image
              source={icons.composeFull}
              className="w-[20px] h-[20px]"
              tintColor={"#475569"}
            /> */}
          </View>
          <View className="mt-3 px-4">
            <View>
              <Text className="font-light text-xs text-slate-600">First Name</Text>
              <Text className="font-pregular text-lg text-slate-800">{userDetail.firstName}</Text>
            </View>
            {
              !userDetail.middleName ? <></> :
              <View className="mt-1">
                <Text className="font-light text-xs text-slate-600">Middle Name</Text>
                <Text className="font-pregular text-lg text-slate-800">{userDetail.middleName}</Text>
              </View>
            }
            <View className="mt-1">
              <Text className="font-light text-xs text-slate-600">Last Name</Text>
              <Text className="font-pregular text-lg text-slate-800">{userDetail.lastName}</Text>
            </View>
          </View>
        </View>
        <View
          className="mt-6 py-3 border-[1px] border-slate-600 rounded-md bg-white shadow-2xl"
        >
          <View
            className="px-4 pb-2 flex flex-row items-center justify-between border-b-[1px] border-slate-400"
          >
            <Text className="font-psemibold text-lg">Your Address</Text>
            <TouchableOpacity
              onPress={() => {router.push("/settings/update_address")}}
              className="p-1"
            >
              <Image
                source={icons.composeFull}
                className="w-[20px] h-[20px]"
                tintColor={"#41ab5d"}
              />
            </TouchableOpacity>
          </View>
          <View className="mt-3 px-4">
            <View>
              <Text className="font-light text-xs text-slate-600">Address Line 1</Text>
              <Text className="font-pregular text-lg text-slate-800">{userDetail.address1}</Text>
            </View>
            {
              !userDetail.address2 ? <></> :
              <View className="mt-1">
                <Text className="font-light text-xs text-slate-600">Address Line 2</Text>
                <Text className="font-pregular text-lg text-slate-800">{userDetail.address2}</Text>
              </View>
            }
            <View className="mt-1 flex flex-row">
              <View className="w-52">
                <Text className="font-light text-xs text-slate-600">City</Text>
                <Text className="font-pregular text-lg text-slate-800">{userDetail.city}</Text>
              </View>
              <View className="">
                <Text className="font-light text-xs text-slate-600">Province</Text>
                <Text className="font-pregular text-lg text-slate-800">{Regions[userDetail.country]?.[userDetail.province]?.code}</Text>
              </View>
            </View>
            <View className="mt-1 flex flex-row">
              <View className="w-52">
                <Text className="font-light text-xs text-slate-600">Postal Code</Text>
                <Text className="font-pregular text-lg text-slate-800">{userDetail.postalCode}</Text>
              </View>
              <View className="">
                <Text className="font-light text-xs text-slate-600">Country</Text>
                <Text className="font-pregular text-lg text-slate-800">{Countries[userDetail.country]?.name}</Text>
              </View>
            </View>
          </View>
        </View>
        <View
          className="mt-6 py-3 border-[1px] border-slate-600 rounded-md bg-white shadow-2xl"
        >
          <View
            className="px-4 pb-2 flex flex-row items-center justify-between border-b-[1px] border-slate-400"
          >
            <Text className="font-psemibold text-lg">Contact Info</Text>
            <TouchableOpacity
              onPress={() => {router.push("/settings/update_phone_number")}}
              className="p-1"
            >
              <Image
                source={icons.composeFull}
                className="w-[20px] h-[20px]"
                tintColor={"#41ab5d"}
              />
            </TouchableOpacity>
          </View>
          <View className="mt-3 px-4">
            <View>
              <Text className="font-light text-xs text-slate-600">Phone Number</Text>
              <Text className="font-pregular text-lg text-slate-800">{userDetail.phoneNumber}</Text>
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default profile