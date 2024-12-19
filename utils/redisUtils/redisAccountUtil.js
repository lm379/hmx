const { setData, getData, deleteData } = require('./redisBaseUtil');

/**
 * 存储验证码到 Redis 中，并设置过期时间
 * @param {string} account 用户账号
 * @param {string} verificationCode 验证码
 * @param {string} verifyType 验证方式
 * @param {number} expiresIn 过期时间（秒） 默认5分钟
 * @returns {Promise} 存储操作的结果
 */
async function storeVerificationCode(account, verificationCode, verifyType, expiresIn = 300) {
    try {
        const data = JSON.stringify({ verificationCode, verifyType });
        await setData(account, data, expiresIn);
        return true;
    } catch (error) {
        console.error('Error storing verification code:', error);
        throw error;
    }
}

/**
 * 获取存储的验证码和签名
 * @param {string} account 用户账号
 * @returns {Promise} 存储的验证码和签名
 */
async function getVerificationCode(account) {
    try {
        const data = await getData(account);
        return data ? JSON.parse(data) : null;
    } catch (error) {
        console.error('Error getting verification code:', error);
        throw error;
    }
}

/**
 * 删除存储的验证码和签名
 * @param {string} account 用户账号
 * @returns {Promise} 删除操作的结果
 */
async function deleteVerificationCode(account) {
    try {
        await deleteData(account);
        return true;
    } catch (error) {
        console.error('Error deleting verification code:', error);
        throw error;
    }
}

module.exports = {
    storeVerificationCode,
    getVerificationCode,
    deleteVerificationCode,
};