import { View, Text, SafeAreaView, Switch } from 'react-native'
import React, { useState } from 'react'
import { authService } from '../../service';

const NotificationSettings = () => {
  const [form, setForm] = useState({
    isInAppNotificationEnable: true,
    isPushNotificationEnable: true,
    isEmailNotificationsEnable: true,
  })


  const [isEnabled, setIsEnabled] = useState(false);
  const toggleSwitch = () => setIsEnabled(previousState => !previousState);

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
          console.log('vvvvvv', userSetting)
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
            onValueChange={toggleSwitch}
            value={isEnabled}
          />
        </View>
        <View className="mt-4 flex flex-row justify-between items-center">
          <Text className="font-psemibold text-lg">Push Notifications</Text>
          <Switch
            trackColor={{false: '#cbd5e1', true: '#005a32'}}
            ios_backgroundColor="#cbd5e1"
            onValueChange={toggleSwitch}
            value={isEnabled}
          />
        </View>
        <View className="mt-4 flex flex-row justify-between items-center">
          <Text className="font-psemibold text-lg">Email Notifications</Text>
          <Switch
            trackColor={{false: '#cbd5e1', true: '#005a32'}}
            ios_backgroundColor="#cbd5e1"
            onValueChange={toggleSwitch}
            value={isEnabled}
          />
        </View>
      </View>
    </SafeAreaView>
  )
}

export default NotificationSettings