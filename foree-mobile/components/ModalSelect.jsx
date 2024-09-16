import { View, Text, TextInput, TouchableOpacity, Modal, ScrollView, SafeAreaView } from 'react-native'
import React, { useEffect, useState } from 'react'
import string_util from '../util/string_util'

const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
  
}

const SelectCountryItem =(country) => (
  <Text className="font-pregular py-3 text-xl">
    {`${country["unicodeIcon"]}`} {country["name"]}
  </Text>
)

//TODO: Redesign value currently is bonded by value Extractor
// the real value for the form may different compare with value that show in UI.
const ModalSelect = ({
  title,
  modalTitle,
  value,
  containerStyles,
  inputStyles,
  allowSearch=true,
  searchKey,
  keyExtractor,
  showExtractor,
  valueExtractor,
  list=[],
  listView,
  uselistSeperator=true,
  allowAdd=true,
  addHandler,
  inputContainerStyles,
  onPress=()=>{},
  variant='bordered',
  placeholder,
  errorMessage,
  isEditable=true
}) => {
  const [showList, setShowList] = useState(list)
  const [visible, setVisible] = useState(false)
  const [cache, setCache] = useState(new Map(list.map((obj) => [obj[valueExtractor], obj])))

  useEffect(() => {
    setCache(new Map(list.map((obj) => [obj[valueExtractor], obj])))
    setShowList(list)
  }, [list])

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
            onPress={() => {
              if (isEditable) setVisible(true)
              else {}
            }}
            className={`h-full w-full px-4 flex-1 font-psemibold text-base ${inputStyles}`}
            value={typeof value === "function" ? value() : cache.get(value)?.[showExtractor]}
            placeholder={placeholder}
            placeholderTextColor="#BDBDBD"
            // onChangeText={handleChangeText}
            editable={false}
          />
        </View>
          {
            !!errorMessage ?
            <View>
              <Text className="text-red-600">{errorMessage}</Text>
            </View> : null
          }
      </View>
      <Modal 
        visible={visible}
        animationType='slide'
      >
        <SafeAreaView className="h-full flex flex-col">
          <View
            className="flex flex-row items-center border-b-[1px] border-slate-400 pr-4"
          >
            <View className="flex-1 flex flex-row items-center">
              <Text
                onPress={() => {
                  setVisible(false)
                  setShowList(list)
                }}
                className="py-2 px-4 text-2xl font-bold text-slate-600"
              >
                &#8592;
              </Text>
              <Text
                className="font-psemibold text-xl text-slate-600"
              >{modalTitle}</Text>
            </View>
            <View>
              {
                allowAdd ?
                  <TouchableOpacity
                    onPress={()=> {
                      addHandler()
                      setVisible(false)
                    }}
                    activeOpacity={0.7}
                    className="bg-[#1A6B54] py-2 px-4 rounded-full"
                    disabled={false}
                  >
                    <Text className="font-pextrabold text-white">New</Text>
                  </TouchableOpacity>
                : null
              }
              
            </View>
          </View>
          <View 
            className="px-2 flex-1"
          >
            {
              allowSearch ? <View
                className="w-full h-14 my-2 border-2 border-secondary rounded-2xl flex-row items-center"
              >
                <Text className="px-2">&#128270;</Text>
                <TextInput
                  className={`flex-1 h-full font-pregular text-base`}
                  // value={"aaa"}
                  placeholder="searching..."
                  editable={true}
                  keyboardType="default"
                  placeholderTextColor="#BDBDBD"
                  onChangeText={(t)=>{
                    if ( !t ) setShowList(list)
                    else {
                      setShowList(list.filter(v => string_util.containSubsequence(typeof searchKey === "function" ? searchKey(v) : v[searchKey], t, {caseInsensitive:true})))
                    }
                  }}
                />
                  
                {/* <Text>Search</Text> */}
              </View> : null
            }
            <View className="flex-1">
              <ScrollView 
                className=""
                showsVerticalScrollIndicator={false}
              >
                { !showList || showList.length === 0 ? 
                  <View className="w-full">
                    <Text 
                      className="font-psemibold text-center py-4 text-xl"
                    >🚫 Empty</Text> 
                  </View>
                  :
                  showList.map((v) => 
                  (
                    <TouchableOpacity
                      onPress={() => {
                        onPress(v[valueExtractor])
                        setVisible(false)
                        setShowList(list)
                      }}
                      className={`w-full ${uselistSeperator ? "border-b-[1px] border-slate-300" : ""}`}
                      key={v[keyExtractor]}
                    >
                      {listView(v)}
                    </TouchableOpacity>
                  ))
                }
              </ScrollView>
            </View>
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