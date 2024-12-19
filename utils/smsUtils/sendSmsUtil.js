'use strict';
const Dysmsapi20170525 = require('@alicloud/dysmsapi20170525');
const OpenApi = require('@alicloud/openapi-client');
const Util = require('@alicloud/tea-util');
const Tea = require('@alicloud/tea-typescript');

class Client {

    static createClient() {
        let config = new OpenApi.Config({
            accessKeyId: process.env.ALIBABA_CLOUD_ACCESS_KEY_ID,
            accessKeySecret: process.env.ALIBABA_CLOUD_ACCESS_KEY_SECRET,
        });
        config.endpoint = `dysmsapi.aliyuncs.com`;
        return new Dysmsapi20170525.default(config);
    }

    static async main(args) {
        let client = Client.createClient();
        let sendSmsRequest = new Dysmsapi20170525.SendSmsRequest({
            signName: process.env.ALIBABA_CLOUD_SIGN_NAME,
            templateCode: process.env.ALIBABA_CLOUD_TEMPLATE_CODE,
            phoneNumbers: args.phoneNumbers,
            templateParam: args.templateParam,
        });
        let runtime = new Util.RuntimeOptions({});
        try {
            await client.sendSmsWithOptions(sendSmsRequest, runtime);
        } catch (error) {
            // 错误 message
            console.log(error.message);
            // 诊断地址
            console.log(error.data["Recommend"]);
            Util.default.assertAsString(error.message);
            throw new Error(error.message);
        }
    }
}

async function sendSms(phoneNumbers, Code) {
    try {
        const templateParam = JSON.stringify({
            code: Code
        });
        await Client.main({
            phoneNumbers: phoneNumbers,
            templateParam: templateParam
        });
        return { code: 0, msg: "短信验证码发送成功，请注意查收" };
    } catch (error) {
        console.log(error);
        return { code: 1, msg: "短信验证码发送失败" + error.message };
    }
}

module.exports = sendSms;