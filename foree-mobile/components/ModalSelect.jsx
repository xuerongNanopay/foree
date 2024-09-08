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
  modalTitle="select a country",
  value,
  containerStyles,
  inputStyles,
  allowSearch=true,
  searchKey,
  list,
  listView,
  allowAdd=true,
  addTitle="Add New Contact",
  addHandler,
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
            w-full h-12 bg-slate-100 ${variants[variant] ?? variants.bordered}
            border-slate-400 focus:border-secondary-200
            items-center flex-row ${inputContainerStyles}
          `}
        >
          <TextInput
            onPress={() => setVisible(true)}
            className={`h-full w-full px-4 flex-1 font-psemibold text-base ${inputStyles}`}
            value={value}
            placeholder={placeholder}
            placeholderTextColor="#BDBDBD"
            // onChangeText={handleChangeText}
            editable={false}
          />

          {/* <View>
            <Text 
              className="mx-4 text-2xl font-bold text-[#BDBDBD]"
            >
              &gt;
            </Text>
          </View> */}
        </View>
      </View>
      <Modal 
        visible={visible}
        onTouchCancel={() => setVisible(false)} 
        animationType='slide'
      >
        <SafeAreaView>
          <View
            className="flex flex-row items-center border-b-[1px] border-slate-400"
          >
            <Text
              onPress={() => setVisible(false)}
              className="py-2 pl-4 pr-8 text-2xl font-bold text-slate-600"
            >
              &#8592;
            </Text>
            <Text
              className="font-psemibold text-xl text-slate-600"
            >{modalTitle}</Text>
          </View>
          <View className="px-2">
            {
              allowSearch ? <View
                className="w-full h-14 my-2 border-2 border-secondary rounded-full flex-row items-center"
              >
                <Text className="px-2">&#128270;</Text>
                <TextInput
                  className={`flex-1 h-full font-pregular text-base`}
                  // value={"aaa"}
                  placeholder="searching..."
                  editable={true}
                  keyboardType="default"
                  placeholderTextColor="#BDBDBD"
                  onChangeText={()=>{}}
                />
                  
                {/* <Text>Search</Text> */}
              </View> : null
            }
            <View>
              {
                allowAdd ? <View
                  className="w-full h-14 my-2 border-2 border-secondary rounded-full flex-row items-center"
                >
                  <TouchableOpacity
                    className="w-full"
                  >
                    <Text className="text-center font-semibold text-secondary text-xl"><Text className="text-2xl">+</Text> {addTitle}</Text>
                  </TouchableOpacity>
                </View> : null
              }
              
            </View>
            <View>
              <Text>List</Text>
            </View>
          </View>
        </SafeAreaView>
      </Modal>
    </View>
  )
}

export default ModalSelect