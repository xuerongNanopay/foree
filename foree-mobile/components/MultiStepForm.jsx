import { View, Text, ScrollView } from 'react-native'
import React, { useState, useEffect } from 'react'
import CustomButton from './CustomButton'

const MultiStepForm = ({
  containerStyle,
  showProgress=true,
  submitTintTitle = 'Submit',
  onSumbit = () => {},
  steps = () => [],
  isSubmitting=false
}) => {
  const formStep = steps()
  const [curIdx, setCurIdx] = useState(0)
  const [isFirst, setIsFirst] = useState(curIdx===0)
  const [isLast, setIsLast] =  useState(curIdx===formStep.length-1)
  const [progress, setProgress] = useState(curIdx+1/formStep.length)
  const [progressCss, setProgressCss] = useState('-100%')

  useEffect(() => {
    setIsFirst(curIdx===0)
    setIsLast(curIdx===formStep.length-1)
    setProgress(curIdx+1/formStep.length)
  }, [curIdx])

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
          {
            formStep[curIdx].formView()
          }
          <View className="mt-10">
            {
              !isFirst ?
              <CustomButton
                title={"Previous"}
                isLoading={isSubmitting}
                handlePress={()=>{setCurIdx(curIdx-1)}}
                containerStyles="mb-2"
                variant="bordered"
                  // isLoading={isSubmitting}
              /> : null
            }
            {
              !isLast ? 
              <CustomButton
                title={"Next"}
                isLoading={isSubmitting}
                handlePress={()=>{setCurIdx(curIdx+1)}}
                containerStyles="mb-2"
                  // isLoading={isSubmitting}
              /> : null
            }
            {
              isLast ? 
              <CustomButton
                title={submitTintTitle}
                isLoading={isSubmitting}
                handlePress={onSumbit}
                containerStyles="mb-4"
                  // isLoading={isSubmitting}
              /> : null
            }
          </View>
        </View>
      </ScrollView>
    </View>
  )
}

export default MultiStepForm