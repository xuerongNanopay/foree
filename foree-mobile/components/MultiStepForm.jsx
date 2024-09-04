import { View, Text, ScrollView } from 'react-native'
import React, { useState, useEffect } from 'react'
import CustomButton from './CustomButton'

const MultiStepForm = ({
  containerStyle,
  showProgress=true,
  submitTintTitle = 'Submit'
}) => {
  const [progress, setProgress] = useState(0.30)
  const [progressCss, setProgressCss] = useState('-100%')
  console.log(progressCss)
  useEffect(()=>{
    setProgressCss('-' + Math.round((1-progress)*100) + '%')
  }, [progress])

  return (
    <View className={`${containerStyle}`}>
      { showProgress ? 
        (
          <View
            className="h-[3px] bg-slate-300 relative"
          >
            <View 
              className={`border-t-[3px] border-secondary top-[0px]`}
              style={{
                left: progressCss
              }}
            />
          </View>
        ) : null
      }
      <ScrollView
        className="h-full"
        automaticallyAdjustKeyboardInsets
      >
        <View>

        
        </View>
        <View className="px-2">
          <View>
            <CustomButton
              title={"Previous"}
              handlePress={()=>{}}
              containerStyles="mb-2"
                // isLoading={isSubmitting}
            />
            <CustomButton
              title={"Next"}
              handlePress={()=>{}}
              containerStyles="mb-2"
                // isLoading={isSubmitting}
            />
            <CustomButton
              title={submitTintTitle}
              handlePress={()=>{}}
              containerStyles="mb-4"
                // isLoading={isSubmitting}
            />
          </View>
        </View>
      </ScrollView>
    </View>
  )
}

export default MultiStepForm