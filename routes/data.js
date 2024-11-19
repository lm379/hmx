const express = require('express');
const router = express.Router();
const { ListUsers, ListOpera, getUserInfo, getOperaInfo, getFavorites, getPlayHistory, getPlayCount, getComments, postComment, addOpera, deleteOpera, updateOpera, addFavorite, deleteFavorite, getMenu, search } = require('../dataBase/api');
const jwt = require('jsonwebtoken');
const secretKey = require('../utils/secretKey');

// 列出用户列表
router.get('/list_users', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    // 只有管理员可以查看用户列表
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, msg: 'permission denied' });
    }
    try {
        const data = await ListUsers();
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 列出曲目列表
router.get('/list_opera', async (req, res, next) => {
    try {
        const data = await ListOpera();
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 获取用户信息
router.get('/get_user_info', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    let Role, decoded;
    try {
        decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    // 只有管理员可以查看任意用户信息
    let User_id, Phone, Email, Username;
    if (Role === 'Administrator') {
        User_id = req.query.user_id;
        Phone = req.query.phone;
        Email = req.query.email;
        Username = req.query.username;
    } else {
        User_id = decoded.User_id;
        Phone = decoded.Phone;
        Email = decoded.Email;
        Username = decoded.Username;
    }
    if (!User_id && !Phone && !Email && !Username) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await getUserInfo(User_id, Phone, Email, Username);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 获取曲目信息
router.get('/get_opera_info', async (req, res, next) => {
    const Opera_id = req.query.opera_id;
    const Opera_title = req.query.opera_title;
    if (!Opera_id && !Opera_title) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await getOperaInfo(Opera_id, Opera_title);
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 获取单用户收藏夹信息
router.get('/get_favorites', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    };
    if (!User_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await getFavorites(User_id, null);
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 获取单用户播放历史
router.get('/get_playhistory', async (req, res, next) => {
    const User_id = req.query.user_id;
    const Opera_id = req.query.opera_id;
    if (!User_id && !Opera_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await getPlayHistory(User_id, Opera_id);
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 统计每首曲目的播放次数
router.get('/get_playcount', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    let user_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        user_id = decoded.user_id;
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    // 只有管理员可以查看播放次数
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, msg: 'permission denied' });
    }
    try {
        const data = await getPlayCount();
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 根据用户或曲目列出所有评论
router.get('/get_comments', async (req, res, next) => {
    const User_id = req.query.user_id;
    const Opera_id = req.query.opera_id;
    if (!User_id && !Opera_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await getComments(Opera_id, User_id, null);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 发表评论
router.post('/post_comment', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    const Opera_id = req.body.opera_id;
    const Comment_text = req.body.text;
    if (!User_id || !Opera_id || !Comment_text) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await postComment(User_id, Opera_id, Comment_text);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'post failed' });
    }
});

// 删除评论
router.post('/delete_comment', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    let User_id, Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }

    const Comment_id = req.body.comment_id;
    if (!Comment_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }

    try {
        const comment = await getComments(null, User_id, Comment_id);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }

    // 非管理员只可以删除自己的评论
    if (Role !== 'Administrator' && comment.code !== 0) {
        return res.status(403).json({ code: 1, msg: 'permission denied' });
    }
    try {
        const data = await deleteComment(User_id, Comment_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'delete failed' });
    }
});

// 新增戏曲
router.post('/add_opera', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameter' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    // 仅管理员调用
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, msg: 'permission denied' });
    }
    const Opera_title = req.body.opera_title;
    const Artist = req.body.artist;
    const Release_date = req.body.release_date;
    const Description = req.body.description;
    const Duration = req.body.duration;
    const Music_path = req.body.music_path;
    const Video_path = req.body.video_path;
    if (!Opera_title || !Music_path || !Video_path) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await addOpera(Opera_title, Artist, Release_date, Duration, Music_path, Video_path, Description);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'add failed' });
    }
});

// 删除戏曲
router.post('/delete_opera', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameter' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    // 仅管理员调用
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, msg: 'permission denied' });
    }
    const Opera_id = req.body.opera_id;
    if (!Opera_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await deleteOpera(Opera_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'delete failed' });
    }
});

// 修改戏曲信息
router.post('/update_opera', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameter' });
    }
    let Role;
    try {
        const decoded = jwt.verify(token, secretKey);
        Role = decoded.Role;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    // 仅管理员调用
    if (Role !== 'Administrator') {
        return res.status(403).json({ code: 1, msg: 'permission denied' });
    }
    const Opera_id = req.body.opera_id;
    const UpdateFields = req.body.updateFields;
    const Value = req.body.value;
    if (!Opera_id || !UpdateFields || !Value) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    if (UpdateFields !== 'Opera_title' && UpdateFields !== 'Artist' && UpdateFields !== 'Release_date' && UpdateFields !== 'Duration' && UpdateFields !== 'Music_path' && UpdateFields !== 'Video_path' && UpdateFields !== 'Description') {
        return res.status(400).json({ code: 1, msg: 'invalid field' });
    }
    try {
        const data = await updateOpera(Opera_id, UpdateFields, Value);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'update failed' });
    }
});

// 添加收藏
router.post('/add_favorite', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameter' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    const Opera_id = req.body.opera_id;
    if (!Opera_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    // 检查是否存在该戏曲
    let opera = await getOperaInfo(Opera_id, null);
    if (opera.code !== 0) {
        return res.status(400).json({ code: 1, msg: 'Opera not exist' });
    }
    // 检查是否已收藏
    opera = await getFavorites(User_id, Opera_id);
    if (opera.code === 0) {
        return res.status(400).json({ code: 1, msg: 'already in favorites' });
    }
    try {
        const data = await addFavorite(User_id, Opera_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'add failed' });
    }
});

// 取消收藏
router.post('/delete_favorite', async (req, res, next) => {
    const token = req.headers['authorization'];
    if (!token) {
        return res.status(400).json({ code: 1, msg: 'invalid parameter' });
    }
    let User_id;
    try {
        const decoded = jwt.verify(token, secretKey);
        User_id = decoded.User_id;
    } catch (err) {
        return res.status(401).json({ code: 1, msg: 'invalid token' });
    }
    const Opera_id = req.body.opera_id;
    if (!Opera_id) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    // 检查是否存在该戏曲
    let opera = await getOperaInfo(Opera_id, null);
    if (opera.code !== 0) {
        return res.status(400).json({ code: 1, msg: 'Opera not exist' });
    }
    // 检查是否已收藏
    opera = await getFavorites(User_id, Opera_id);
    if (opera.code !== 0) {
        return res.status(400).json({ code: 1, msg: 'not in favorites' });
    }
    // 取消收藏
    try {
        const data = await deleteFavorite(User_id, Opera_id);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'delete failed' });
    }
});

// 获取菜单
router.get('/get_menu', async (req, res, next) => {
    try {
        const data = await getMenu();
        res.json(data)
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'query failed' });
    }
});

// 模糊搜索
router.get('/search', async (req, res, next) => {
    const token = req.headers['authorization'];
    let Table, Field;
    // 无token仅允许搜索曲目
    if (!token) {
        Table = 'Opera';
        Field = 'Opera_title';
        Column = ['*'];
    } else {
        // 预留接口，暂未实现
        Table = 'Opera';
        Field = 'Opera_title';
        Column = ['*'];
        // let Role;
        // try {
        //     const decoded = jwt.verify(token, secretKey);
        //     Role = decoded.Role;
        // } catch (err) {
        //     return res.status(401).json({ code: 1, msg: 'invalid token' });
        // }
        // if (Role === 'Administrator') {
        //     Table = req.query.table;
        //     Field = req.query.field;
        // } else {
        //     Table = 'Opera';
        //     Field = 'Opera_title';
        // }
    }
    const Keywords = req.query.keywords;
    if (!Keywords) {
        return res.status(400).json({ code: 1, msg: 'invalid parameters' });
    }
    try {
        const data = await search(Keywords, Table, Field, Column);
        res.json(data);
    } catch (err) {
        console.error(err); // 记录错误日志
        res.status(500).json({ code: 1, msg: 'search failed' });
    }
});

module.exports = router;