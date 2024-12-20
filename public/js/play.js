// Cookie.set('Admin', false);
var app = new Vue({
    el: '#app',
    data: {
        UserInfo: null,
        ItemList: [],
        pageTitle: null,
        loading: true,
        CommentsCount: 0,
        aiResult: '',
    },
    mounted: function () {
        // 初始化
        this.getUserInfo();
        this.pageTitle = '加载中';
        // 从url获取戏曲id
        let Opera_id = this.getQuery('opt');
        this.Opera_id = Opera_id;

        // 查询该歌曲信息
        this.getOperaInfo(Opera_id);
        this.getComment(Opera_id);
        this.updateCommentsCount(Opera_id);

    },
    watch: {
        pageTitle: function (newTitle) {
            document.title = newTitle;
        }
    },
    methods: {
        // 从url获取参数
        getQuery: function (value) {
            let query = window.location.search.substring(1);
            let vars = query.split("&");
            for (const i of vars) {
                let pair = i.split("=");
                if (pair[0] == value) { return pair[1]; }
            }
            return false;
        },

        // 获取用户信息
        getUserInfo: function () {
            let token = Cookie.get('Token'),
                type = Cookie.get('UserType');

            if (token && type != 'admin') {
                HTTP.get('get_user_info', null, 'info').then(json => {
                    if (json.code === 0) {
                        this.UserInfo = json.data[0];
                        if (!json.data[0].Icon) { this.UserInfo.Icon = './img/logo.png' }
                    } else {
                        Cookie.del('Token')
                        Cookie.del('UserType')
                    }
                }).catch(error => {
                    Cookie.del('Token');
                    Cookie.del('UserType');
                });
            }
        },

        // 获取当前戏曲信息
        getOperaInfo: function (Opera_id) {
            HTTP.get('get_opera_info', { opera_id: Opera_id }, 'info').then(json => {
                if (json.code === 0) {
                    this.ItemList = json.data[0];
                    this.pageTitle = this.ItemList.Opera_title;
                    this.loading = false;
                    this.$nextTick(() => {
                        const player = new Plyr('#player');
                    });
                } else {
                    this.ItemList = [];
                    this.loading = false;
                }
            })
        },

        // 获取评论数
        getCommentsCount: function (Opera_id) {
            if (!Opera_id) { return Promise.resolve(0); }
            return HTTP.get('get_comments_count', { params: "Opera_id", value: Opera_id }, 'info').then(json => {
                if (json.code === 0) {
                    let count = json.data.Comments_count;
                    return count;
                } else {
                    return 0;
                }
            });
        },

        // 更新评论数
        updateCommentsCount: function (Opera_id) {
            this.getCommentsCount(Opera_id).then(count => {
                this.CommentsCount = count;
            }).catch(err => {
                console.log(err);
            });
        },

        // 获取评论
        getComment: function (Opera_id) {
            CommentsList = [];
            Promise.all([
                HTTP.get('get_comments', { opera_id: Opera_id }, 'info'),
            ])
                .then(results => {
                    const json = results[0];
                    if (json.code === 0) {
                        // if (!json.data.Icon) { json.data.Icon = './img/logo.png' }
                        json.data.forEach(item => {
                            if (!item.Icon) {
                                item.Icon = './img/logo.png'
                            }
                        });
                        this.CommentsList = json.data;
                    } else {
                        this.CommentsList = [];
                    }
                });
        },

        // 格式化时间戳
        formatDate: function (dateString) {
            const options = { year: 'numeric', month: '2-digit', day: '2-digit' };
            return new Date(dateString).toLocaleDateString('zh-CN', options);
        },

        // 搜索戏曲
        search: function (event) {
            let value = event.target.value;
            if (value == '') { this.getProduct() }
            else {
                HTTP.get('search', { keywords: value }, 'info').then(json => {
                    if (json.code == 0) { this.ItemList = json.data }
                    else { this.ItemList = [] }
                })
            }
        },

        // 输入后隐藏提示
        handleInput: function (event) {
            const editor = event.target;
            const placeholder = editor.previousElementSibling;
            if (editor.textContent.trim() === '') {
                placeholder.style.display = 'block';
            } else {
                placeholder.style.display = 'none';
            }
        },

        // 跳转页面
        jump: function (event) {
            // 获取需要跳转的页面
            let page = event.target.dataset.page + '.html';
            // 获取跳转参数
            let option = event.target.dataset.opt;
            if (option) { page += '?opt=' + option }
            window.location.href = page;
        },

        // 退出登录
        logout: function () {
            // 弹出提示框
            let check = confirm("确定退出登录吗?");
            if (check) {
                Cookie.del('Token')
                window.location.href = 'login.html'
            }
        },
        // 退出登录
        logout: function () {
            // 弹出提示框
            let check = confirm("确定退出登录吗?");
            if (check) {
                Cookie.del('Token')
                window.location.href = 'login.html'
            }
        },

        // 异步发表评论
        async postComment() {
            const editor = document.querySelector('.editor');
            const commentText = editor.textContent.trim();
            if (commentText === '') {
                alert('评论内容不能为空');
                return;
            }

            try {
                // 发送评论请求
                const res = await HTTP.post('post_comment', { opera_id: this.Opera_id, text: commentText }, 'info');
                if (res.code === 0) {
                    // 清空编辑器内容
                    editor.textContent = '';
                    this.handleInput({ target: editor });
                    // 更新评论列表
                    await this.getComment(this.Opera_id);
                    await this.getCommentsCount(this.Opera_id);
                } else {
                    alert('评论发表失败：' + res.msg);
                }
            } catch (error) {
                console.error('Failed to post comment:', error);
            }
        },

        // 处理ai请求
        getAiResult: function () {
            HTTP.get('get_ai_result', { opera_id: this.Opera_id }, 'info').then(json => {
                if (json.code === 0) {
                    this.aiResult = '';
                    this.typeWriter(json.data[0].text, 0);
                    // this.aiResult = json.data[0].text;
                } else {
                    this.aiResult = '';
                    setTimeout(() => {
                        this.typeWriter('数据还在生成中，当前戏曲排在第0位', 0);
                    }, 1000);
                }
            }).catch(error => {
                console.error('Failed to get AI result:', error);
            });
        },

        // 逐字动效
        typeWriter: function (text, index) {
            if (index < text.length) {
                this.aiResult += text.charAt(index);
                setTimeout(() => {
                    this.typeWriter(text, index + 1);
                }, 50); // 调整速度，100ms 显示一个字
            }
        },
    }
});