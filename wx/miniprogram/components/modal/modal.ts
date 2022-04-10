// pages/components/modal/modal.ts
import {ModalResult} from "./types";

Component({
    /**
     * 组件的属性列表
     */
    properties: {
        showModal:Boolean,
        title:String,
        contents:String,
        showCancel:Boolean,
    },

    /**
     * 组件的初始数据
     */
    data: {
        resolve:undefined as ((r:ModalResult) => void) | undefined
    },

    /**
     * 组件的方法列表
     */
    methods: {
        onCancel(){
            this.hideModal('cancel')
        },

        onOk(){
            this.hideModal('ok')
        },

        hideModal(res:ModalResult){
            this.setData({
                showModal:false,
            })
            this.triggerEvent(res)
            if(this.data.resolve){
                this.data.resolve(res)
            }
        },

        showModal():Promise<ModalResult>{
            this.setData({
                showModal:true,
            })

            return new Promise((resolve) => {
                this.data.resolve = resolve
            })
        }

    },
    options:{
        addGlobalClass:true,
    }
});
