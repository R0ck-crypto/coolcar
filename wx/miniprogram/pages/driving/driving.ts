import { routing } from "../../utils/util"
import {TripService} from "../../service/trip";

// pages/driving/driving.ts
const centPerSec = 1



function formatDuration(sec:number){

    const padString = (n:number ) => 
    n < 10 ? '0' + n.toFixed(0):n.toFixed(0)
    const  h = Math.floor(sec/3600)
    sec -= 3600 * h
    const m  = Math.floor(sec / 60)
    sec -= 60 * m
    const s = Math.floor(sec)
    return `${padString(h)}:${padString(m)}:${padString(s)}`
}

function formatFee(cents:number){
    return (cents / 100).toFixed(2)
}


Page({

    /**
     * 页面的初始数据
     */

    tripID:'',
    data: {
        timer:undefined as number | undefined,
        tripID:'',
        location:{
            latitude:40.003304,
            longitude:116.326759,
        },
        scale:14,
        elapsed:'00:00:00',
        fee:'0.0'
    },


    onLoad(opt:Record<'trip_id',string>){
        const o:routing.DrivingOpts = opt
        this.tripID = o.trip_id
        TripService.GetTrip(o.trip_id).then(console.log)
        this.setupLocationUpdator()
        this.setupTimer()
    },

    onUnload(){
        wx.stopLocationUpdate()
        if(this.timer){
            clearInterval(this.timer)
        }
    },

    setupLocationUpdator(){
        wx.startLocationUpdate({
            fail:console.error,
        })

        wx.onLocationChange(loc => {
            this.setData({
                location:{
                    latitude:loc.latitude,
                    longitude:loc.longitude,
                }
            })
        })
    },

    setupTimer(){
        let elapsedSec = 0
        let cents = 0
        this.timer = setInterval(() =>{
            elapsedSec++
            cents += centPerSec
            this.setData({
                elapsed:formatDuration(elapsedSec),
                fee:formatFee(cents)
            })
        },1000)

    },

    onEndTripTap(){
        TripService.finishTrip(this.tripID).then(() => {
            wx.redirectTo({
                url:routing.mytrips(),
            }).catch(err => {
                console.error(err)
                wx.showToast({
                    title:'结束行程失败',
                    icon:'none',
                })
            })
        })
    }
})