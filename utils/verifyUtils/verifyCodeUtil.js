const { getVerificationCode, deleteVerificationCode } = require("../redisUtils/redisAccountUtil");

const verifyCode = async (Account, Code) => {
    try {
        const storedData = await getVerificationCode(Account);
        if (!storedData) {
            return { code: 1, msg: '验证码已过期，请重新接收验证码' };
        }

        const storedVerificationCode = storedData.verificationCode;

        if (Code !== storedVerificationCode) {
            return { code: 1, msg: '验证码错误' };
        } else {
            // 验证成功，删除redis中的验证码
            await deleteVerificationCode(Account);
            return { code: 0, msg: '验证成功' };
        }
    } catch {
        return { code: 1, msg: '服务器错误' };
    }
};

module.exports = verifyCode;