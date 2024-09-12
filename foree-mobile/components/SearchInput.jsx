import { View, Text, TextInput, TouchableOpacity, Image } from 'react-native'
import React, { useState } from 'react'

import { icons } from '../constants'

const variants = {
  bordered: "border-[1px] rounded",
  flat: "border-b-[1px]"
  
}

const SearchInput = ({
  value, 
  placeholder, 
  handleChangeText, 
  inputStyles,
  containerStyles,
  variant='bordered',
  keyboardType='default', 
  editable=true,
  errorMessage,
  ...props
}) => {
  const [showPassword, setShowPassword] = useState(false)


  return (
      <View 
        className={`
          w-full h-12 ${variants[variant] ?? variants.bordered}
          border-slate-400 items-center flex-row ${containerStyles}
        `}
      >
        <Text className="px-2">🔍</Text>
        <TextInput
          className={`flex-1 h-full font-psemibold text-base ${inputStyles}`}
          autoCorrect={false}
          spellCheck={false}
          value={value}
          placeholder={placeholder}
          placeholderTextColor="#BDBDBD"
          onChangeText={handleChangeText}
          keyboardType={keyboardType} 
          editable={editable}
        />

        {/* <TouchableOpacity
        onPress={()=> setShowPassword(!showPassword)}
        >
        <Image 
            source={icons.search}
            className="w-7 h-7"
            resizeMode='contain'
        />
        </TouchableOpacity> */}
      </View>
  )
}

export default SearchInput