import { SafeAreaView, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native'
import React from 'react'
import { Link } from 'expo-router'

const ProfileTab = () => {
  return (
    <SafeAreaView
      className=""
    >
      <View
        className="h-full px-4"
      >
        <View className="mt-4 mb-6">
          <Text className="font-pbold text-2xl">Settings</Text>
        </View>
        <ScrollView
          className="border h-full"
        >
          <View>
            <View
              className="mb-6"
            >
              <View className="border-b-2 border-slate-300">
                <Text className="font-psemibold text-xl">Account</Text>
              </View>
            </View>
            <Link href="/login"
              className="border-2 border-red-700 bg-red-200 p-1 rounded-lg"
            >
              <Text className="font-semibold text-lg text-red-700 text-center">Logout</Text>
            </Link>
          </View>
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

export default ProfileTab

const styles = StyleSheet.create({})