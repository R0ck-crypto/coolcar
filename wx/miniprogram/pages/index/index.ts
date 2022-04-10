import {routing} from "../../utils/util"
import {TripService} from "../../service/trip";
import {rental} from "../../service/proto_gen/rental/rental_pb";
import {ProfileService} from "../../service/profile";


Page({
    userInfo: {},
    hasUserInfo: false,
    isPageShow: false,
    data: {
        latitude: 22.625404,
        longitude: 114.060874,
        markers: [
            {
                iconPath: "/resources/car.png",
                id: 0,
                latitude: 22.625404,
                longitude: 114.060874,
                name: '五和',
                width: 50,
                height: 50,
            },
            {
                iconPath: "/resources/car.png",
                id: 1,
                latitude: 23.625404,
                longitude: 114.060874,
                name: '龙岗纪检委',
                width: 50,
                height: 50,
            },
        ],
        scale: 15,
        //   covers: [{
        //     latitude: 22.625404,
        //     longitude: 113.344520,
        //     iconPath: '/image/location.png'
        //   }, {
        //     latitude: 22.625404,
        //     longitude: 113.304520,
        //     iconPath: '/image/location.png'
        //   }]
    },
    onShow() {
        this.isPageShow = true
    },

    onHide() {
        this.isPageShow = false
    },

    onLoad() {

    },

    onMyLocationTap() {
        wx.getLocation(
            {
                type: 'gcj02',
                success: res => {
                    this.setData(
                        {
                            latitude: res.latitude,
                            longitude: res.longitude,
                        },
                    )
                },
                fail: () => {
                    wx.showToast(
                        {
                            icon: 'none',
                            title: '请前往设置页授权',
                        }
                    )
                }
            }
        )
    },

    async onScanTap() {
        const trips = await TripService.GetTrips(rental.v1.TripStatus.IN_PROGRESS)
        if ((trips.trips?.length || 0) > 0) {
            await this.selectComponent('#tripModal').showModal()
            wx.navigateTo({
                url: routing.driving({
                    trip_id: trips.trips![0].id!
                })
            })
            return
        }
        wx.scanCode({
                success: async () => {
                    const carID = '62526826538283d56ac2d8d0'
                    // const redirectURL = `/pages/lock/lock?car_id=${carID}`
                    const lockURL = routing.lock({
                        car_id: carID,
                    })
                    const prof = await ProfileService.getProfile()
                    if (prof.identityStatus == rental.v1.IdentityStatus.VERIFIED) {
                        wx.navigateTo({
                            url: lockURL,
                        })
                    } else {
                        await this.selectComponent('#licModal').showModal()
                        // TODO: get car id from scan result

                        wx.navigateTo({
                            //url:`/pages/register/register?redirect=${encodeURIComponent(redirectURL)}`,
                            url: routing.register({
                                redirectURL: lockURL,
                            })
                        })

                    }
                },
                fail: console.error,
            }
        )
    },
    moveCars() {
        const map = wx.createMapContext("myMap")
        const dest = {
            latitude: 22.625404,
            longitude: 114.060874,
        }
        const moveCar = () => {
            dest.latitude += 0.1
            dest.longitude += 0.1
            map.translateMarker({
                    destination: {
                        latitude: dest.latitude + 1,
                        longitude: dest.longitude + 1,
                    },
                    markerId: 0,
                    autoRotate: false,
                    rotate: 0,
                    duration: 5000,
                    animationEnd: () => {
                        if (this.isPageShow) {
                            moveCar()
                        } else {

                        }
                    }
                }
            )
        }
        moveCar()

    },
    onMyTripsTap() {
        wx.navigateTo({
            url: routing.mytrips(),
        })
    },

})
  