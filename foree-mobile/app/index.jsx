import { StatusBar } from 'expo-status-bar'
import { ScrollView, Text, View } from 'react-native'
import { Redirect, router } from 'expo-router'
import { SafeAreaView } from 'react-native'
import { images } from '../constants'
import { Image } from 'react-native'
import CustomButton from '../components/CustomButton'
import { useGlobalContext } from '../context/GlobalProvider'

export default function App() {
  const {isLoading, isLoggedIn} = useGlobalContext()
  if ( !isLoading && isLoggedIn ) return <Redirect href="/home_tab"/>

  return (
    <SafeAreaView className="h-full">
      <ScrollView contentContainerStyle={{ height: '100%' }}>
        <View className="w-full items-center min-h-[85vh] px-2">
          <Image
            source={images.logo}
            className="w-[120px] h-[84px]"
            resizeMode='contain'
          />
          <Image
            source={images.cards}
            className="max-w-[380px] w-full h-[300px]"
            resizeMode='contain'
          />

          <View className="relative mt-5">
            <Text className="text-3xl text-slate-900 font-bold text-center">
              Transfer money with{' '}
              <Text className="text-secondary-200">Foree</Text>
            </Text>
            <Image
              source={images.path}
              className="w-[136px h-[15px] absolute -bottom-2 right-4"
              resizeMode='contain'
            />
          </View>
          <Text
            className="text-sm font-pregular text-gray-600 mt-4 text-center"
          >
            Where creativeit meets invotaion embark fdsa pakkk fdsaf asdf 
          </Text>
          <CustomButton
            title="Continue With Email"
            handlePress={() => router.push('/login')}
            containerStyles="w-full mt-7"
          />
        </View>
      </ScrollView>

      <StatusBar backgroundColor='#161622' style='auto'/>
    </SafeAreaView>
  );
}