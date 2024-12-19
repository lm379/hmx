Cookie.set('Admin', false)
var app = new Vue({
    el: '#app',
    data: {
        UserInfo: null,
        ItemList: [],
    },
    mounted: function () {
        // 初始化
        this.getUserInfo()
        this.getProduct()
        this.getFavorites()
    },
    methods: {
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

        // 获取戏曲列表
        getProduct: function (type) {
            if (type && type != '全部戏曲') { options.Type = type }
            // 小标题显示当前正在浏览的内容类型
            if (type) { this.MenuActive = type }

            HTTP.get('list_opera', null, 'info').then(json => {
                if (json.code === 0) {
                    this.ItemList = json.data;
                    this.updateFavoriteStatus();
                }
                else { this.ItemList = [] }
            })
        },

        // 获取用户收藏的戏曲
        getFavorites: function () {
            HTTP.get('get_favorites', null, 'info').then(json => {
                if (json.code === 0) {
                    this.Favorites = json.data;
                    this.updateFavoriteStatus();
                }
            })
        },

        // 更新收藏状态
        updateFavoriteStatus: function () {
            if (!this.Favorites) return;
            this.ItemList.forEach(item => {
                item.isFavorite = this.Favorites.some(fav => fav.Opera_id === item.ID)
            })
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
    }
});


