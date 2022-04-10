import { routing } from "../../utils/util"
import {rental} from "../../service/proto_gen/rental/rental_pb";
import {padString} from "../../utils/format";
import {ProfileService} from "../../service/profile";
import {Coolcar} from "../../service/request";




function formatDate(millis:number){
    const dt = new Date(millis)

    const y = dt.getFullYear()
    const m = dt.getMonth() + 1
    const d = dt.getDate()

    return `${padString(y)}-${padString(m)}-${padString(d)}`
}


// pages/register/register.ts
Page({
    redirectURL:'',
    profileRefresher:0,
    /**
     * 页面的初始数据
     */
    data: {
        licNo:'',
        name:'',
        genderIndex:0,
        genders:['未知','男','女','其他'],
        LicImgURL:'',
        birthDate:'1990-01-01',
        state :rental.v1.IdentityStatus[rental.v1.IdentityStatus.UNSUBMITTED],
    },

    renderProfile(p:rental.v1.IProfile){
        this.renderIdentity(p.identity!)
        this.setData({
            state:rental.v1.IdentityStatus[p.identityStatus || 0],
        })
    },

    renderIdentity(i?:rental.v1.IIdentity){
        this.setData({
            licNo:i?.licNumber || '',
            name:i?.name || '',
            genderIndex:i?.gender || 0,
            birthDate:formatDate(i?.birthDateMillis ||  0),
        })

    },

    onLoad(opt:Record<'redirect',string>){
        const o:routing.RegisterOpts = opt
        if(o.redirect){
            this.redirectURL = decodeURIComponent(o.redirect)
        }
        ProfileService.getProfile().then(p => this.renderProfile(p))
        ProfileService.getProfilePhoto().then(p => {
            this.setData({
                LicImgURL:p.url || '',
            })
        })
    },

    onUploadLic: function () {
        wx.chooseImage({
                success: async res => {
                    if (res.tempFilePaths[0].length === 0) {
                        return
                    }
                    this.setData({
                        LicImgURL: res.tempFilePaths[0],
                    })

                    const photoRes = await ProfileService.createProfilePhoto()

                    if (!photoRes.uploadUrl){
                        return
                    }
                    await Coolcar.uploadfile({
                        localPath:res.tempFilePaths[0],
                        url:photoRes.uploadUrl,
                    })

                    const identity = await ProfileService.completeProfilePhoto()
                    this.renderIdentity(identity)
                }
            }
        )
    },

    onGenderChange(e:any){
        this.setData({
            genderIndex : parseInt(e.detail.value),
        })
    },

    onBirthDateChange(e:any){
        this.setData({
            birthDate : e.detail.value,
        })
    },

    onSubmit(){
        ProfileService.submitProfile({
            licNumber:this.data.licNo,
            name:this.data.name,
            gender:this.data.genderIndex,
            birthDateMillis:Date.parse(this.data.birthDate)
        }).then(p => {
            this.renderProfile(p)
            this.scheduleProfileRefresher()
        }
        )
    },

    onUnload() {
        this.clearProfilerefresher()
    },


    scheduleProfileRefresher(){
        this.profileRefresher = setInterval(()=>{
            ProfileService.getProfile().then(p => {
                this.renderProfile(p)
                if (p.identityStatus !== rental.v1.IdentityStatus.PENDING){
                    this.clearProfilerefresher()
                }
                if (p.identityStatus === rental.v1.IdentityStatus.VERIFIED){
                    this.onLicVerified()
                }
            })
        }, 1000)
    },

    clearProfilerefresher(){
        if (this.profileRefresher){
            clearInterval(this.profileRefresher)
            this.profileRefresher = 0
        }
    },


    onReSubmit() {
        ProfileService.clearProfile().then(p => this.renderProfile(p))
        ProfileService.clearProfilePhoto().then(() => {
            this.setData({
                LicImgURL:'',
            })
        })

    },

    onLicVerified(){
        if(this.redirectURL){
            wx.redirectTo({
                url:this.redirectURL,
            })
        }
    },



    


 

    // /**
    //  * 生命周期函数--监听页面加载
    //  */
    // onLoad() {

    // },

    // /**
    //  * 生命周期函数--监听页面初次渲染完成
    //  */
    // onReady() {

    // },

    // /**
    //  * 生命周期函数--监听页面显示
    //  */
    // onShow() {

    // },

    // /**
    //  * 生命周期函数--监听页面隐藏
    //  */
    // onHide() {

    // },

    // /**
    //  * 生命周期函数--监听页面卸载
    //  */
    // onUnload() {

    // },

    // /**
    //  * 页面相关事件处理函数--监听用户下拉动作
    //  */
    // onPullDownRefresh() {

    // },

    // /**
    //  * 页面上拉触底事件的处理函数
    //  */
    // onReachBottom() {

    // },

    // /**
    //  * 用户点击右上角分享
    //  */
    // onShareAppMessage() {

    // }
})