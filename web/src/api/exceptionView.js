import service from '@/utils/request'

export const getExceptionView = () => {
  return service({
    url: "/exception/viewException",
    method: 'post',
  })
}