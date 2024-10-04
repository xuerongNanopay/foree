import { View, Text, SafeAreaView, ScrollView, Image, TouchableOpacity } from 'react-native'
import React from 'react'
import { icons } from '../../constants'

const profile = () => {
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
              <Text className="font-pregular text-lg text-slate-800">Addfdsa dasfd</Text>
            </View>
            {
              1==1 ? <></> :
              <View className="mt-1">
                <Text className="font-light text-xs text-slate-600">Middle Name</Text>
                <Text className="font-pregular text-lg text-slate-800">Addfdsa dasfd</Text>
              </View>
            }
            <View className="mt-1">
              <Text className="font-light text-xs text-slate-600">Last Name</Text>
              <Text className="font-pregular text-lg text-slate-800">Addfdsa dasfd</Text>
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
              className="p-1"
            >
              <Image
                source={icons.composeFull}
                className="w-[20px] h-[20px]"
                tintColor={"#475569"}
              />
            </TouchableOpacity>
          </View>
          <View className="mt-3 px-4">
            <View>
              <Text className="font-light text-xs text-slate-600">Address Line 1</Text>
              <Text className="font-pregular text-lg text-slate-800">56 Colonsay Rd</Text>
            </View>
            {
              1 === 1 ? <></> :
              <View className="mt-1">
                <Text className="font-light text-xs text-slate-600">Address Line 2</Text>
                <Text className="font-pregular text-lg text-slate-800"></Text>
              </View>
            }
            <View className="mt-1 flex flex-row">
              <View className="flex-1">
                <Text className="font-light text-xs text-slate-600">City</Text>
                <Text className="font-pregular text-lg text-slate-800">Thornhill</Text>
              </View>
              <View className="mr-12">
                <Text className="font-light text-xs text-slate-600">Province</Text>
                <Text className="font-pregular text-lg text-slate-800">ON</Text>
              </View>
            </View>
            <View className="mt-1 flex flex-row">
              <View className="flex-1">
                <Text className="font-light text-xs text-slate-600">Postal Code</Text>
                <Text className="font-pregular text-lg text-slate-800">L3T 3E8</Text>
              </View>
              <View className="mr-12">
                <Text className="font-light text-xs text-slate-600">Country</Text>
                <Text className="font-pregular text-lg text-slate-800">CA</Text>
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
              className="p-1"
            >
              <Image
                source={icons.composeFull}
                className="w-[20px] h-[20px]"
                tintColor={"#475569"}
              />
            </TouchableOpacity>
          </View>
          <View className="mt-3 px-4">
            <View>
              <Text className="font-light text-xs text-slate-600">Phone Number</Text>
              <Text className="font-pregular text-lg text-slate-800">3065022222</Text>
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}

export default profile