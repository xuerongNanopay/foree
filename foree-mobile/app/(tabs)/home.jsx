import { View, Text, SafeAreaView, FlatList, ScrollView } from 'react-native'
import { Link } from 'expo-router'
import React from 'react'
import { useGlobalContext } from '../../context/GlobalProvider'

const Home = () => {
  const { user } = useGlobalContext()

  return (
    <SafeAreaView>
      <View className="px-4 pt-4">
        <View className="mb-4">
          <Text className="font-pregular text-xl">Welcome Back</Text>
          <Text className="font-pbold text-2xl text-slate-700">{user?.firstName} {user?.lastName}</Text>
        </View>
        <ScrollView>
        {/* <FlatList
          data={[{id: 1}]}
        /> */}
          <View className="bg-[#f2f7f5] rounded-2xl p-4">
            <View className="">
              <View className="flex-1">
                <Text className="font-pbold text-lg">Current Rate</Text>
                <Text className="font-psemibold text-lg">🇨🇦 $1.00 CAD = 🇵🇰 $208.00 PKR</Text>
              </View>
              <View>
                <View className="mt-4 p-2 rounded-xl bg-[#1A6B54]">
                  <Link href="/profile">
                    <Text className="text-lg text-center font-semibold text-white">Send Money</Text>
                  </Link>
                </View>
              </View>
            </View>
          </View>
        </ScrollView>
      </View>
    </SafeAreaView>
  )
}

export default Home