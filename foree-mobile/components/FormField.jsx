import { View, Text, TextInput, TouchableOpacity, Image } from 'react-native'
import React, { useState } from 'react'

import { icons } from '../constants'

const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
  
}

const FormField = ({
  title, 
  value, 
  placeholder, 
  handleChangeText, 
  containerStyles,
  inputStyles,
  inputContainerStyles,
  variant='bordered',
  keyboardType='default', 
  editable=true,
  ...props
}) => {
  const [showPassword, setShowPassword] = useState(false)


  return (
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
          className={`flex-1 h-full font-psemibold text-base ${inputStyles}`}
          value={value}
          placeholder={placeholder}
          placeholderTextColor="#BDBDBD"
          onChangeText={handleChangeText}
          secureTextEntry={title === 'Password' && !showPassword}
          keyboardType={keyboardType} 
          editable={editable}
        />

        {title === 'Password' && (
          <TouchableOpacity
            onPress={()=> setShowPassword(!showPassword)}
          >
            <Image 
              source={!showPassword ? icons.eye : icons.eyeHide}
              className="w-7 h-7"
              resizeMode='contain'
            />
          </TouchableOpacity>
        )}
      </View>
    </View>
  )
}

export default FormField