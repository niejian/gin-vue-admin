import service from '@/utils/request'

// @Summary 环境初始化
// @Produce  application/json
// @Router/watchdog/init [post]
export const initEnv = (data) => {
  return service({
    url: "/watchdog/init",
    method: 'post',
    data
  })
}