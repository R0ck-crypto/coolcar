import { routing } from "../../utils/util"
import {rental} from "../../service/proto_gen/rental/rental_pb";
import {ProfileService} from "../../service/profile";


const licStatusMap = new Map([
    [rental.v1.IdentityStatus.UNSUBMITTED, '未认证'],
    [rental.v1.IdentityStatus.PENDING,'未认证'],
    [rental.v1.IdentityStatus.VERIFIED,'已认证']
])
// pages/mytrips/mytrips.ts
Page({

    /**
     * 页面的初始数据
     */
    data: {
        licStatus:licStatusMap.get(rental.v1.IdentityStatus.UNSUBMITTED),
        hasUserInfo:false,
        userInfo:{},
        autoplay:true,
        promotionItems:[
            {
                img:'https://img3.mukewang.com/61f354e500011d8117920764.jpg',     
                promotionID:1,
            },
            
            {
                img: 'https://img2.mukewang.com/61f5f3520001e43e17920764.jpg',     
                promotionID:2,
            },
            
            {
                img:'https://img3.mukewang.com/61f5f3d70001052f17920764.jpg',     
                promotionID:3,
            },

            {
                img:'https://img3.mukewang.com/61f5f4060001bd5717920764.jpg',     
                promotionID:4,
            },

            {
                img: 'https://img3.mukewang.com/61f5f44800017edb17920764.jpg',    
                promotionID:4,
            },
        ]
    },

    getUserProfile(e){
        wx.getUserProfile({
            desc:'用于获取用户头像',
            success:(res) => {
                this.setData({
                    userInfo:res.userInfo,
                    hasUserInfo:true,
                })
            }
        })
    },
        


    onSwiperChange(e){
       // console.log(e)
        if(e.detail.source){
            // cause by our program
            return
        }
        // process
    },

    onPromotionItemTap(e:any){
        console.log(e)
        const promotionID = e.currentTarget.dataset.promotionID
        if(promotionID){
            //claim the promotionID
        }
    },

    onRegisterTap(){
        wx.navigateTo({
            //url:'/pages/register/register',
            url:routing.register(),
        })
    },
    /**
     * 生命周期函数--监听页面加载
     */
    onLoad() {

    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady() {

    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow() {
        ProfileService.getProfile().then(p =>{
            this.setData({
                licStatus:licStatusMap.get(p.identityStatus || 0)
            })
        })
    },

    /**
     * 生命周期函数--监听页面隐藏
     */
    onHide() {

    },

    /**
     * 生命周期函数--监听页面卸载
     */
    onUnload() {

    },

    /**
     * 页面相关事件处理函数--监听用户下拉动作
     */
    onPullDownRefresh() {

    },

    /**
     * 页面上拉触底事件的处理函数
     */
    onReachBottom() {

    },

    /**
     * 用户点击右上角分享
     */
    onShareAppMessage() {

    }
})