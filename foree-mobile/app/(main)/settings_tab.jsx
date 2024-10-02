import { Alert, Image, SafeAreaView, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native'
import React from 'react'
import { Link, router } from 'expo-router'
import { icons } from '../../constants'

const ProfileTab = () => {
  return (
    <SafeAreaView
      className=""
    >
      <View
        className="h-full px-4"
      >
        <View className="mt-4 mb-6 flex flex-row items-center justify-between">
          <Text className="font-pbold text-2xl">Settings</Text>
            <TouchableOpacity
              onPress={() => {
                Alert.alert("Logout", "Are you sure?", [
                  {text: 'Confirm', onPress: () => {router.replace("/login")}},
                  {text: 'Cancel', onPress: () => {}},
                ])
              }}
              className=""
            >
              <Image
                source={icons.signOutOutline}
                className="h-[20px] w-[20px] mr-2"
                resizeMode='contain'
                tintColor={"#005a32"}
              />  
            </TouchableOpacity>
        </View>
        <ScrollView
          className="h-full"
        >
          <View>
            <View
              className="mb-6"
            >
              <View className="mb-2 border-b-2 border-slate-300 flex flex-row items-center">
                <Image
                  source={icons.userOutline}
                  className="h-[14px] w-[14px] mr-2"
                  resizeMode='contain'
                  tintColor={"#005a32"}
                />
                <Text className="font-psemibold text-xl">Account</Text>
              </View>
              <TouchableOpacity
                onPress={() => {
                  router.push("/settings/profile")
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Personal Settings</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={() => {
                  router.push("/settings/update_passwd")
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Change Password</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={() => {
                  router.push("/settings/notification_settings")
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Notification Settings</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={() => {
                  router.push("/settings/profile")
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Refer Others</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
            </View>
            <View>
              <View className="mb-2 border-b-2 border-slate-300 flex flex-row items-center">
                <Image
                  source={icons.infoOutline}
                  className="h-[14px] w-[14px] mr-2"
                  resizeMode='contain'
                  tintColor={"#005a32"}
                />
                <Text className="font-psemibold text-xl">Info</Text>
              </View>
              <TouchableOpacity
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Contacts</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <TouchableOpacity
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Terms and Conditions</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <TouchableOpacity
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Privacy Policy</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
            </View>
          </View>
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

export default ProfileTab

const styles = StyleSheet.create({})