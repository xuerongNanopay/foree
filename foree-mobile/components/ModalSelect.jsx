import { View, Text, TextInput, TouchableOpacity, Modal, SafeAreaView, ScrollView } from 'react-native'
import React, { useState } from 'react'

import Countries from '../constants/country'

const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
  
}
const aaa = Object.values(Countries)
const SelectCountryItem =(country) => (
  <Text className="font-pregular py-3 text-xl">
    {`${country["unicodeIcon"]}`} {country["name"]}
  </Text>
)

const ModalSelect = ({
  title,
  modalTitle="select a country",
  value,
  containerStyles,
  inputStyles,
  allowSearch=true,
  searchKey="name",
  list,
  listView,
  allowAdd=true,
  addTitle,
  addHandler,
  inputContainerStyles,
  onPress=()=>{},
  variant='bordered',
  placeholder
}) => {
  const [showList, setShowList] = useState(list)
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
            <ScrollView className="h-full">
              { !showList || showList.length === 0 ? 
                <View className="w-full border-b-[1px] border-slate-300">
                  <Text 
                    className="font-psemibold text-center py-4 text-xl"
                  >🚫 Empty</Text> 
                </View>
                :
                showList.map((v) => 
                (
                  <TouchableOpacity
                    onPress={() => {
                      onPress(v)
                      setVisible(false)
                    }}
                    className="w-full border-b-[1px] border-slate-300"
                    key={v[searchKey]}
                  >
                    {SelectCountryItem(v)}
                  </TouchableOpacity>
                ))
              }
            </ScrollView>
          </View>
        </SafeAreaView>
      </Modal>
    </View>
  )
}

export default ModalSelect
export {
  SelectCountryItem
}