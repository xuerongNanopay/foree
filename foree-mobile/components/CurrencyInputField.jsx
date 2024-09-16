import { View, Text, TouchableOpacity, TextInput, Modal } from 'react-native'
import React, { useState } from 'react'

const variants = {
  bordered: "border-2 rounded-2xl",
  flat: "border-b-2"
}

const CurrencyInputField = ({
  title,
  containerStyles,
  variant='bordered',
  inputContainerStyles,
  supportCurrency=[],
  value,
  onCurrencyChange=()=>{}
}) => {
  const [visible, setVisible] = useState(false)
  const [currency, setCurrency] = useState("")
  const [amount, setAmount] = useState(0.0)

  return (
    <>
      <View className={`${containerStyles}`}>
        <Text className="test-base test-gray-100 font-pmedium mb-2">{title}</Text>
        <View
          className={`
            flex flex-row items-center
            h-12 ${variants[variant] ?? variants.bordered}
            border-slate-400 focus:border-secondary-200 ${inputContainerStyles}
          `}
        >
          <TouchableOpacity
            onPress={()=>{console.log("press currency")}}
            activeOpacity={0.7}
            className="bg-slate-200 h-full px-2 rounded-l-2xl flex justify-center"
          >
            <Text className="font-semibold">🇨🇦 CAD 🔽</Text>
          </TouchableOpacity>
          <TextInput
            keyboardType='decimal-pad'
            className="flex-1 h-full text-right px-2 font-semibold"
          />
        </View>
      </View>
      {/* Currency picker */}
      <Modal
        visible={visible}
        animationType='slide'
      >
      
      </Modal>
    </>
  )
}

export default CurrencyInputField