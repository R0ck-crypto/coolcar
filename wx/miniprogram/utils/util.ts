export const formatTime = (date: Date) => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return (
    [year, month, day].map(formatNumber).join('/') +
    ' ' +
    [hour, minute, second].map(formatNumber).join(':')
  )
}

const formatNumber = (n: number) => {
  const s = n.toString()
  return s[1] ? s : '0' + s
}


export namespace routing{
  export interface DrivingOpts{
    trip_id:string
  }

  export function driving(o:DrivingOpts){
    return `/pages/driving/driving?trip_id=${o.trip_id}`
  }

  export interface LockOpts{
    car_id:string
  }

  export function lock(o:LockOpts){
    return `/pages/lock/lock?car_id=${o.car_id}`
  }

  export interface RegisterOpts{
    redirect?:string
  }

  export interface RegisterParms{
    redirectURL:string
  }

  export function register(p?:RegisterParms){
    const page = '/pages/register/register'
    if(!p){
      return page
    }
    return `${page}?redirect=${encodeURIComponent(p.redirectURL)}`
  }

  export function mytrips(){

    const page = '/pages/mytrips/mytrips'
    return page
  }
  
}

