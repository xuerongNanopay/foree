import { View, Text, TouchableOpacity, Alert } from 'react-native'
import React from 'react'
import { SafeAreaView } from 'react-native'
import FormField from '../../components/FormField'

const CloseAccount = () => {
  return (
    <SafeAreaView className="h-full">
      <View className="w-full mt-4 px-2">
        <Text className="text font-pbold text-center m-3">Close My Foree Remittance Account</Text>
        <Text className="text-center text-sm text-slate-600">
          Clost your Foree Remittance account will permanently revoke access to your account. As a regulated financial entity, we are required by law to maintain some personal data associated with your transactions.
        </Text>
        <Text className="mt-1 text-center text-sm text-slate-600">To see how we treat your data, please refer to our Privacy Policy</Text>
        <Text className="mt-2 mb-4 font-psemibold text-center text-red-600">You can not undo this action.</Text>
        <FormField
          title="Closing Reason(optional)"
          multiline={true}
          numberOfLines={4}
          inputContainerStyles="h-44"
        />
        <TouchableOpacity
          className="mt-4 py-2 border-2 border-red-800 bg-red-200 rounded-xl"
          onPress={() => {
            Alert.alert("Close Account", "Are you sure?", [
              {text: 'Continue', onPress: () => {console.log("TODO: close account")}},
              {text: 'Cancel', onPress: () => {}},
            ])
          }}
        >
          <Text className="font-pbold text-lg text-red-800 text-center">Close my account</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  )
}

export default CloseAccount