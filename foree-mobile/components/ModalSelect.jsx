import { View, Image, Text, TextInput, TouchableOpacity, Modal, SafeAreaView } from 'react-native'
import React, { useState } from 'react'

import { icons } from '../constants'
const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
  
}

const Touchable = ({text="select a country"}) => {
  const TouchableComponent = () => (
    <TouchableOpacity
      onPress={() => {
        alert('touch')
      }}
    >
      <Text>{text}</Text>
    </TouchableOpacity>
  )
  return {TouchableComponent}
}
const ModalSelect = ({
  title="Nationality",
  value,
  containerStyles,
  inputStyles,
  inputContainerStyles,
  variant='bordered',
  placeholder="select an option"
}) => {
  // const { TouchableComponent }  = Touchable(placeHolder)
  const [visible, setVisible] = useState(false)
  return (
    <View>
      {/* <TouchableComponent>ModalSelect</TouchableComponent> */}
      <View className={`space-y-2 ${containerStyles}`}>
      { !!title ? (<Text className="test-base test-gray-100 font-pmedium">{title}</Text>) : null }
        <View
          className={`
            w-full h-12 px-4 bg-slate-100 ${variants[variant] ?? variants.bordered}
            border-slate-400 focus:border-secondary-200
            items-center flex-row ${inputContainerStyles}
          `}
        >
          <TextInput
            onPress={() => setVisible(true)}
            className={`flex-1 font-psemibold text-base ${inputStyles}`}
            value={value}
            placeholder={placeholder}
            placeholderTextColor="#BDBDBD"
            // onChangeText={handleChangeText}
            editable={false}
          />

          <View>
            <Text 
              className="text-2xl font-bold text-[#BDBDBD]"
            >
              &gt;
            </Text>
          </View>
        </View>
      </View>
      <Modal 
        visible={visible} 
        onTouchCancel={() => setVisible(false)} 
        animationType='slide'
      >
        <SafeAreaView>
          <View>
            <Text>Title</Text>
          </View>
          <View>
            <Text>Can Create</Text>
          </View>
          <View>
            <Text>Can search</Text>
          </View>
          <View>
            <Text>List</Text>
          </View>
        </SafeAreaView>
      </Modal>
    </View>
  )
}

export default ModalSelect