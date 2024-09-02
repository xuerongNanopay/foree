import React from 'react'
import { Tabs, Redirect } from 'expo-router'
import { Image, View, Text } from 'react-native'

import { icons } from '../../constants'

const TabIcon = ({ icon, color, name, focused }) => {
  return (
    <View className="items-center justify-center">
      <Image
        source={icon}
        resizeMethod='contain'
        tintColor={color}
        className="w-5 h-5"
      />
      <Text className={`${focused ? 'font-psemibold' : 'font-pregular'} text-xs`}>
        {name}
      </Text>
    </View>
  )
}
const TabsLayout = () => {
  return (
    <>
      <Tabs
        screenOptions={{
          tabBarShowLabel: false
        }}
      >
        <Tabs.Screen 
          name="home"
          options={{
            title: 'Home',
            headerShown: false,
            tabBarIcon:({ color, forced}) => (
              <TabIcon 
                icon={icons.home}
                color={color}
                name="Home"
                forced={forced}
              />
            )
          }}
        />
      </Tabs>
    </>
  )
}

export default TabsLayout