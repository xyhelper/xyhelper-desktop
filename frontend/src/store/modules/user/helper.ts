// import { v4 as uuidv4 } from 'uuid'
import { ss } from '@/utils/storage'


const LOCAL_NAME = 'userStorage'

export interface UserInfo {
  avatar: string
  name: string
  description: string
  baseURI: string
  accessToken: string
}

export interface UserState {
  userInfo: UserInfo
}

export function defaultSetting(): UserState {
  // 生成随机的字符串
  function generateRandomString(length: number): string {
    let result = ''
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    for (let i = 0; i < length; i++)
      result += characters.charAt(Math.floor(Math.random() * characters.length))

    return result
  }
  const randomString = generateRandomString(10);
  // console.log(randomString); // 输出类似于：lJRObYwExl


  return {
    userInfo: {
      avatar: '/defaultavatar.jpeg',
      name: '攻城狮老李',
      description: '官网: https://xyhelper.cn',
      baseURI: 'https://freechat.lidong.xin',
      accessToken: randomString,
    },
  }
}

export function getLocalState(): UserState {
  const localSetting: UserState | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalState(setting: UserState): void {
  ss.set(LOCAL_NAME, setting)
}
