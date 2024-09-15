import React from 'react'
import { Tabs } from 'expo-router'
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
      <Text className={`${focused ? 'font-psemibold' : 'font-pregular'} text-xs`} style={{ color: color }}>
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
          tabBarShowLabel: false,
          tabBarActiveTintColor: '#004D40',
          tabBarInactiveTintColor: '#009688',
          tabBarStyle: {
            backgroundColor: '#FAFAFA',
            borderTopWidth: 1,
            borderTopColor: '#E0E0E0',
            height: 84
          }
        }}
      >
        <Tabs.Screen 
          name="home_tab"
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
        <Tabs.Screen 
          name="contact_tab"
          options={{
            title: 'Contact',
            headerShown: false,
            tabBarIcon:({ color, forced}) => (
              <TabIcon 
                icon={icons.bookmark}
                color={color}
                name="Contact"
                forced={forced}
              />
            )
          }}
        />
        <Tabs.Screen 
          name="transaction_tab"
          options={{
            title: 'Transaction',
            headerShown: false,
            tabBarIcon:({ color, forced}) => (
              <TabIcon 
                icon={icons.plus}
                color={color}
                name="Transaction"
                forced={forced}
              />
            )
          }}
        />
        <Tabs.Screen 
          name="profile_tab"
          options={{
            title: 'Profile',
            headerShown: false,
            tabBarIcon:({ color, forced}) => (
              <TabIcon 
                icon={icons.profile}
                color={color}
                name="Profile"
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