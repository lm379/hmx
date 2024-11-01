const express = require('express');
const app = express.Router();
const { ListUsers, ListOpera, getUserInfo, getOperaInfo, getFavorites, getPlayHistory, getPlayCount, getComments, postComment, addOpera, deleteOpera, updateOpera, addFavorite ,deleteFavorite} = require('../dataBase/api');
const jwt = require('jsonwebtoken');
const { secretKey } = require('./account');

// 列出用户列表
app.get('/list_users', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    // 只有管理员可以查看用户列表
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, message: 'permission denied' });
    }
    try {
        const data = await ListUsers();
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 列出曲目列表
app.get('/list_opera', async (req, res, next) => {
    try {
        const data = await ListOpera();
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 获取用户信息
app.get('/get_user_info', async (req, res, next) => {
    const User_id = req.query.user_id;
    const Phone = req.query.phone;
    const Email = req.query.email;
    const Username = req.query.username;
    if (!User_id && !Phone && !Email && !Username) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await getUserInfo(User_id, Phone, Email, Username);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 获取曲目信息
app.get('/get_opera_info', async (req, res, next) => {
    const Opera_id = req.query.opera_id;
    const Opera_title = req.query.opera_title;
    if (!Opera_id && !Opera_title) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await getOperaInfo(Opera_id, Opera_title);
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 获取单用户收藏夹信息
app.get('/get_favorites', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    };
    if (!User_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await getFavorites(User_id,null);
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 获取单用户播放历史
app.get('/get_playhistory', async (req, res, next) => {
    const User_id = req.query.user_id;
    const Opera_id = req.query.opera_id;
    if (!User_id && !Opera_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await getPlayHistory(User_id, Opera_id);
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 统计每首曲目的播放次数
app.get('/get_playcount', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    let user_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        user_id = decoded.user_id;
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    // 只有管理员可以查看播放次数
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, message: 'permission denied' });
    }
    try {
        const data = await getPlayCount();
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 根据用户或曲目列出所有评论
app.get('/get_comments', async (req, res, next) => {
    const User_id = req.query.user_id;
    const Opera_id = req.query.opera_id;
    if (!User_id && !Opera_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await getComments(Opera_id, User_id, null);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }
});

// 发表评论
app.post('/post_comment', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    const Opera_id = req.body.opera_id;
    const Comment_text = req.body.text;
    if (!User_id || !Opera_id || !Comment_text) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await postComment(User_id, Opera_id, Comment_text);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'post failed' });
    }
});

// 删除评论
app.post('/delete_comment', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    let User_id, Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }

    const Comment_id = req.body.comment_id;
    if (!Comment_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }

    try {
        const comment = await getComments(null, User_id, Comment_id);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'query failed' });
    }

    // 非管理员只可以删除自己的评论
    if (Role !== 'Administrator' && comment.code !== 0) {
        return res.status(403).json({ code: 1, message: 'permission denied' });
    }
    try {
        const data = await deleteComment(User_id, Comment_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'delete failed' });
    }
});

// 新增戏曲
app.post('/add_opera', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameter' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    // 仅管理员调用
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, message: 'permission denied' });
    }
    const Opera_title = req.body.opera_title;
    const Artist = req.body.artist;
    const Release_date = req.body.release_date;
    const Description = req.body.description;
    const Duration = req.body.duration;
    const Music_path = req.body.music_path;
    const Video_path = req.body.video_path;
    if (!Opera_title || !Music_path || !Video_path) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await addOpera(Opera_title, Artist, Release_date, Duration, Music_path, Video_path, Description);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'add failed' });
    }
});

// 删除戏曲
app.post('/delete_opera', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameter' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    // 仅管理员调用
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, message: 'permission denied' });
    }
    const Opera_id = req.body.opera_id;
    if (!Opera_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    try {
        const data = await deleteOpera(Opera_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'delete failed' });
    }
});

// 修改戏曲信息
app.post('/update_opera', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameter' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    // 仅管理员调用
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, message: 'permission denied' });
    }
    const Opera_id = req.body.opera_id;
    const UpdateFields = req.body.updateFields;
    const Value = req.body.value;
    if (!Opera_id || !UpdateFields || !Value) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    if (UpdateFields !== 'Opera_title' && UpdateFields !== 'Artist' && UpdateFields !== 'Release_date' && UpdateFields !== 'Duration' && UpdateFields !== 'Music_path' && UpdateFields !== 'Video_path' && UpdateFields !== 'Description') {
        return res.status(400).json({ code: 1, message: 'invalid field' });
    }
    try {
        const data = await updateOpera(Opera_id, UpdateFields, Value);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'update failed' });
    }
});

// 添加收藏
app.post('/add_favorite', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameter' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    const Opera_id = req.body.opera_id;
    if (!Opera_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    // 检查是否存在该戏曲
    let opera = await getOperaInfo(Opera_id, null);
    if (opera.code !== 0) {
        return res.status(400).json({ code: 1, message: 'Opera not exist' });
    }
    // 检查是否已收藏
    opera = await getFavorites(User_id, Opera_id);
    if (opera.code === 0) {
        return res.status(400).json({ code: 1, message: 'already in favorites' });
    }
    try {
        const data = await addFavorite(User_id, Opera_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'add failed' });
    }
});

// 取消收藏
app.post('/delete_favorite', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, message: 'invalid parameter' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, message: 'invalid token' });
    }
    const Opera_id = req.body.opera_id;
    if (!Opera_id) {
        return res.status(400).json({ code: 1, message: 'invalid parameters' });
    }
    // 检查是否存在该戏曲
    let opera = await getOperaInfo(Opera_id, null);
    if (opera.code !== 0) {
        return res.status(400).json({ code: 1, message: 'Opera not exist' });
    }
    // 检查是否已收藏
    opera = await getFavorites(User_id, Opera_id);
    if (opera.code !== 0) {
        return res.status(400).json({ code: 1, message: 'not in favorites' });
    }
    // 取消收藏
    try {
        const data = await deleteFavorite(User_id, Opera_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, message: 'delete failed' });
    }
});

module.exports = app;