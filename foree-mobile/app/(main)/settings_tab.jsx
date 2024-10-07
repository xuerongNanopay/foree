import { Alert, Image, Modal, SafeAreaView, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'
import { Link, router, useFocusEffect } from 'expo-router'
import { icons } from '../../constants'
import { generalService } from '../../service'

const ProfileTab = () => {
  const [showContactSupport, setShowContactSupport] = useState(false)
  

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
                  {text: 'Continue', onPress: () => {router.replace("/login")}},
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
                  router.push("/settings/personal_settings")
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
                  router.push("/settings/invitation")
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Invitation</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <TouchableOpacity
                onPress={() => {
                  router.push("/settings/close_account")
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Close My Account</Text>
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
                onPress={() => {
                  setShowContactSupport(true)
                }}
                className="flex flex-row items-center justify-between py-2"
              >
                <Text className="font-semibold text-lg text-slate-500">Contacts Support</Text>
                <Image
                  source={icons.rightArrowDark}
                  className="h-[14px] w-[14px]"
                  resizeMode='contain'
                  tintColor={"#adb5bd"}
                />
              </TouchableOpacity>
              <ContactSupportModal 
                visible={showContactSupport}
                closeModal={() => setShowContactSupport(false)}
              />
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

const ContactSupportModal = ({
  visible=false,
  closeModal=_=>{}
}) => {

  const [customerSupport, setCustomerSupport] = useState({})
  useEffect(() => {
    const controller = new AbortController()
    const getCustomerSupport = async() => {
      try {
        console.log("customer support run")
        const resp = await generalService.cusomterSupport({signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("customer support", resp.status, resp.data)
        } else {
          console.log(resp.data.data)
          setCustomerSupport(resp.data.data)
        }
      } catch (e) {
        console.error("customer support", e, e.response, e.response?.status, e.response?.data)
      }
    }
    getCustomerSupport()
    return () => {
      controller.abort()
    }
  }, [])

  return (
  <Modal
    visible={visible}
    animationType='fade'
  >
    <SafeAreaView className="h-full bg-slate-200">
      <TouchableOpacity 
        className="mt-4 mb-4 px-4 flex flex-row items-center"
        onPress={_ => closeModal()}
      >
        <Image
          source={icons.leftArrowDark}
          className="w-[14px] h-[14px] mr-2"
          tintColor={"#475569"}
        />
        <Text
          className="font-psemibold text-2xl text-slate-600"
        >Contact Us</Text>
      </TouchableOpacity>
      <View className="px-2">
        <Text className="px-2 font-light text-slate-600">You can get in touch with us through below platforms. Our Team will reach out to you as soon as it would be possible</Text>
        <View
          className="mt-6 py-4 px-4 rounded-2xl bg-white shadow-lg"
        >
          <Text className="text-slate-400 font-pregular">Customer Support</Text>
          {
            !customerSupport.supportPhoneNumber ? <></> :          
            <View
              className="mt-4 flex flex-row items-center"
            >
              <View
                className="p-2 bg-slate-200 rounded-full mr-2"
              >
                <Image
                  source={icons.phoneOutline}
                  className="w-[14px] h-[14px]"
                  tintColor={"#475569"}
                />
              </View>
              <View>
                <Text className="font-light text-xs text-slate-400">Customer Number</Text>
                <Text className="font-pregular text-lg text-slate-800">{customerSupport.supportPhoneNumber}</Text>
              </View>
            </View>
          }
          {
            !customerSupport.supportEmail ? <></>:
            <View
              className="mt-4 flex flex-row items-center"
            >
              <View
                className="p-2 bg-slate-200 rounded-full mr-2"
              >
                <Image
                  source={icons.mailOutline}
                  className="w-[14px] h-[14px]"
                  tintColor={"#475569"}
                />
              </View>
              <View>
                <Text className="font-light text-xs text-slate-400">Email Address</Text>
                <Text className="font-pregular text-lg text-slate-800">{customerSupport.supportEmail}</Text>
              </View>
            </View>
          }
        </View>
        <View
          className="mt-6 py-4 px-4 rounded-2xl bg-white shadow-lg"
        >
          <Text className="text-slate-400 font-pregular">Social Media</Text>
          {
            !customerSupport.instagram ? <></> :
            <View
              className="mt-4 flex flex-row items-center"
            >
              <View
                className="bg-slate-200 rounded-full mr-2"
              >
                <Image
                  source={icons.instagramColor}
                  className="w-[30px] h-[30px]"
                />
              </View>
              <View>
                <Text className="font-light text-xs text-slate-400">Instagram</Text>
                <Text className="font-pregular text-lg text-slate-800">{customerSupport.instagram}</Text>
              </View>
            </View>
          }
          {
            !customerSupport.twitter ? <></> :
            <View
              className="mt-4 flex flex-row items-center"
            >
              <View
                className="bg-slate-200 rounded-full mr-2"
              >
                <Image
                  source={icons.twitterColor}
                  className="w-[30px] h-[30px]"
                />
              </View>
              <View>
                <Text className="font-light text-xs text-slate-400">Twitter</Text>
                <Text className="font-pregular text-lg text-slate-800">{customerSupport.twitter}</Text>
              </View>
            </View>
          }
          {
            !customerSupport.facebook ? <></> :
            <View
              className="mt-4 flex flex-row items-center"
            >
              <View
                className="bg-slate-200 rounded-full mr-2"
              >
                <Image
                  source={icons.facebookColor}
                  className="w-[30px] h-[30px]"
                />
              </View>
              <View>
                <Text className="font-light text-xs text-slate-400">Facebook</Text>
                <Text className="font-pregular text-lg text-slate-800">{customerSupport.facebook}</Text>
              </View>
            </View>
          }
        </View>
      </View>
    </SafeAreaView>
  </Modal>
  )
}
export default ProfileTab