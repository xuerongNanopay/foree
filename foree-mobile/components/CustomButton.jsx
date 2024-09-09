import { TouchableOpacity, Text } from 'react-native'
import React from 'react'

const variants = {
  solid: {
    buttonStyle: "bg-secondary",
    textStyle: "text-white"
  },
  bordered: {
    buttonStyle: "border-2 border-secondary",
    textStyle: "text-secondary"
  },
}

const CustomButton = ({
  title, 
  handlePress, 
  containerStyles, 
  textStyles, 
  disabled,
  variant='solid'
}) => {
  return (
    <TouchableOpacity
        onPress={handlePress}
        activeOpacity={0.7}
        className={`
            ${variants[variant].buttonStyle} rounded-xl min-h-[48px] justify-center items-center
            ${containerStyles}
            ${disabled ? 'opacity-50' : ''}
        `}
        disabled={disabled}
    >
    <Text 
        className={`${variants[variant].textStyle} font-psemibold text-lg ${textStyles}`}
    >{title}</Text>
    </TouchableOpacity>
  )
}

export default CustomButton