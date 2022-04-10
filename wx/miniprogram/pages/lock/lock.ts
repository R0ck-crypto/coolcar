import { routing } from "../../utils/util"
import {TripService} from "../../service/trip";

const ShareLocationKey  = "share _location"


Page({
    carID:'',
    data: {
      ShareLocation:false,
        userInfo: {},
      hasUserInfo: false,
      canIUseGetUserProfile: false,
    },
    onLoad(opt:Record<'car_id',string>) {
      const o:routing.LockOpts = opt
      this.carID = o.car_id
      if (wx.getUserProfile) {
        this.setData({
          canIUseGetUserProfile: true
        })
      }
    },
    getUserProfile(e) {
      // 推荐使用wx.getUserProfile获取用户信息，开发者每次通过该接口获取用户个人信息均需用户确认
      // 开发者妥善保管用户快速填写的头像昵称，避免重复弹窗
      wx.getUserProfile({
        desc: '用于完善会员资料', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
        success: (res) => {
          this.setData({
            userInfo: res.userInfo,
            hasUserInfo: true,
            ShareLocation:wx.getStorageSync(ShareLocationKey) || false,
          })
        }
      })
    },
    onShareLocation(e:any){
        const ShareLocation:boolean = e.detail.value
        wx.setStorageSync(ShareLocationKey,ShareLocation)
    },
    onUnlockTap(){
        wx.getLocation({
          type:"gcj02",
          success: async loc =>{
            console.log('starting a trip',{
              location:{
                latitude:loc.latitude,
                longitude:loc.longitude,
              },
            })
          if(!this.carID){
              console.log('no carID specified')
              return
          }
          const  trip = await TripService.CreateTrip({
              start:loc,
              carId:this.carID,
          })

          if(!trip.id){
              console.log('no tripid specified')
              return
          }

            wx.showLoading({
              title:'开锁中',
              mask:true,
          })

          setTimeout(() => {
            wx.redirectTo({
                url:routing.driving({
                  trip_id:trip.id!,
                }),
                complete:() => {
                    wx.hideLoading()
                }
            })
        }, 2000);
            
          },
          fail: () => {
            wx.showToast({
              icon:'none',
              title:'请前往设置页授权位置信息',
            })
          }
        })
       

       
    }
    
  })