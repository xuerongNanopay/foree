import { View, Text, SafeAreaView, Switch } from 'react-native'
import React, { useCallback, useEffect, useState } from 'react'
import { authService } from '../../service';
import { router, useFocusEffect } from 'expo-router';

const NotificationSettings = () => {
  const [form, setForm] = useState({
    isInAppNotificationEnable: true,
    isPushNotificationEnable: true,
    isEmailNotificationsEnable: true,
  })

  useFocusEffect(useCallback(() => {
    const controller = new AbortController()
    const getUserSetting = async() => {
      try {
        const resp = await authService.getUserSetting({signal: controller.signal})
        if ( resp.status / 100 !== 2 &&  !resp?.data?.data) {
          console.error("get userSetting", resp.status, resp.data)
          router.replace("/personal_settings")
        } else {
          const userSetting = resp.data.data
          setForm({
            isInAppNotificationEnable: userSetting.isInAppNotificationEnable,
            isPushNotificationEnable: userSetting.isPushNotificationEnable,
            isEmailNotificationsEnable: userSetting.isEmailNotificationsEnable,
          })
        }
      } catch (e) {
        console.error("get userSetting", e, e.response, e.response?.status, e.response?.data)
        router.replace("/personal_settings")
      }
    }
    getUserSetting()
    return () => {
      controller.abort()
    }
  }, []))

  useEffect(() => {
    const updateUserNotification = async() => {
      try {
        const resp = await authService.updateUserSetting(form)
        if ( resp.status / 100 !== 2 ) {
          console.warn("update notification", resp.status, resp.data)
          if ( router.canGoBack ) {
            router.back()
          } else {
            router.replace("/personal_settings")
          }
        }
      } catch (err) {
        console.error("update notification", err)
      }
    }
    updateUserNotification()
  }, [form])

  return (
    <SafeAreaView className="h-full">
      <View
        className="px-4 mt-6"
      >
        <View className="flex flex-row justify-between items-center">
          <Text className="font-psemibold text-lg">In-App Notifications</Text>
          <Switch
            trackColor={{false: '#cbd5e1', true: '#005a32'}}
            ios_backgroundColor="#cbd5e1"
            onValueChange={(v) => {
              setForm((form) => ({
                ...form,
                isInAppNotificationEnable: v,
              }))
            }}
            value={!!form.isInAppNotificationEnable}
          />
        </View>
        <View className="mt-4 flex flex-row justify-between items-center">
          <Text className="font-psemibold text-lg">Push Notifications</Text>
          <Switch
            trackColor={{false: '#cbd5e1', true: '#005a32'}}
            ios_backgroundColor="#cbd5e1"
            onValueChange={(v) => {
              setForm((form) => ({
                ...form,
                isPushNotificationEnable: v,
              }))
            }}
            value={!!form.isPushNotificationEnable}
          />
        </View>
        <View className="mt-4 flex flex-row justify-between items-center">
          <Text className="font-psemibold text-lg">Email Notifications</Text>
          <Switch
            trackColor={{false: '#cbd5e1', true: '#005a32'}}
            ios_backgroundColor="#cbd5e1"
            onValueChange={(v) => {
              setForm((form) => ({
                ...form,
                isEmailNotificationsEnable: v,
              }))
            }}
            value={!!form.isEmailNotificationsEnable}
          />
        </View>
      </View>
    </SafeAreaView>
  )
}

export default NotificationSettings