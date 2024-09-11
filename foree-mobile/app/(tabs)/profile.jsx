import { SafeAreaView, StyleSheet, Text, View } from 'react-native'
import React from 'react'
import { Link } from 'expo-router'

const Profile = () => {
  return (
    <SafeAreaView>
      <View>
        <Text>Profile</Text>
        <Link href="/login">Logout</Link>
      </View>
    </SafeAreaView>
  )
}

export default Profile

const styles = StyleSheet.create({})