import { View, Text, SafeAreaView, Switch } from 'react-native'
import React, { useState } from 'react'

const NotificationSettings = () => {
  // const [notificationSettings, setNotificationSettings] = useState({})
  // const toggleSwitch = (isInAppNotifications) => {

  // }

  const [isEnabled, setIsEnabled] = useState(false);
  const toggleSwitch = () => setIsEnabled(previousState => !previousState);

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