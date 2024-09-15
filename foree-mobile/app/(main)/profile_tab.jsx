import { SafeAreaView, StyleSheet, Text, View } from 'react-native'
import React from 'react'
import { Link } from 'expo-router'

const ProfileTab = () => {
  return (
    <SafeAreaView>
      <View>
        <Text>ProfileTab</Text>
        <Link href="/login">Logout</Link>
      </View>
    </SafeAreaView>
  )
}

export default ProfileTab

const styles = StyleSheet.create({})