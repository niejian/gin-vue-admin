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

export const try2Connect = (data) => {
  return service({
    url: "/watchdog/try2Connect",
    method: 'post',
    data
  })
}

export const sshService = () => {
  return service({
    url: "/ws/:id",
    method: 'get'
  })
}

// 下载配置文件
export const downloadConfig = (data) => {
  return service({
    url: "/watchdog/downloadConfig",
    method: 'POST',
    data
  })
}