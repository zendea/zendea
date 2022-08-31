import qs from 'qs'

export default function({ $axios, $toast, app }) {
  $axios.onRequest((config) => {
    config.headers.common['X-Client'] = 'zendea'
    config.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'
    config.transformRequest = [
      function(data) {
        if (process.client && data instanceof FormData) {
          // 如果是FormData就不转换
          return data
        }
        data = qs.stringify(data)
        return data
      }
    ]
  })

  $axios.onResponse((response) => {
    if (response.status !== 200) {
      return Promise.reject(response)
    }
    const jsonResult = response.data
    if (jsonResult.code === 0) {
      return Promise.resolve(jsonResult.data)
    } else if (jsonResult.code === 200) {
      return Promise.resolve(jsonResult)
    } else {
      return Promise.reject(jsonResult)
    }
  })
}
