// 简化的艺术家信息（用于作品列表）
export interface SimpleArtist {
    artist_id: number;
    name: string;
    avatar?: string;
}

// 简化的作品信息（用于艺术家详情，不包含艺术家避免冗余）
export interface SimpleOpera {
    opera_id: number;
    opera_title: string;
    avatar?: string;
    duration?: string;
    play_count: number;
    like_count: number;
    created_at: string;
    ai_summary?: string;
}

// 列表页的作品信息（简化版）
export interface OperaListItem {
    opera_id: number;
    opera_title: string;
    artists: SimpleArtist[];
    avatar?: string;
    duration?: string;
    is_hidden: boolean;
    play_count: number;
    like_count: number;
    recommend_reason?: string;
    created_at: string;
    ai_summary?: string;
}

// 详情页的作品信息（完整版）
export interface OperaDetail {
    opera_id: number;
    opera_title: string;
    artists: SimpleArtist[];
    release_date?: string;
    duration?: string;
    music_path?: string;
    video_path: string;
    srt_path?: string;
    description: string;
    avatar?: string;
    ai_summary?: string;
    is_hidden: boolean;
    created_at: string;
    updated_at: string;
    like_count: number;
    favorite_count: number;
    share_count: number;
    play_count: number;
    liked: boolean;
    favorited: boolean;
}

// 通用作品类型（向后兼容）
export interface Opera {
    opera_id: number;
    opera_title: string;
    artists: SimpleArtist[];
    release_date?: string;
    duration?: string;
    music_path?: string;
    video_path: string;
    srt_path?: string;
    description: string;
    avatar?: string;
    ai_summary?: string;
    is_hidden?: boolean;
    created_at: string;
    updated_at: string;
    like_count?: number;
    favorite_count?: number;
    share_count?: number;
    play_count?: number;
    liked?: boolean;
    favorited?: boolean;
}

export interface OperaListResponse {
    code: number;
    msg: string;
    data: {
        list: OperaListItem[];
        pagination: {
            total: number;
            page: number;
            page_size: number;
        }
    }
}

// 列表页的艺术家信息（简化版）
export interface ArtistListItem {
    artist_id: number;
    name: string;
    avatar?: string;
    bio?: string;
    created_at: string;
}

// 详情页的艺术家信息（完整版）
export interface ArtistDetail {
    artist_id: number;
    name: string;
    bio: string;
    avatar?: string;
    operas: SimpleOpera[];
    created_at: string;
    updated_at: string;
}

// 通用艺术家类型（向后兼容）
export interface Artist {
    artist_id: number;
    name: string;
    bio: string;
    avatar?: string;
    operas?: SimpleOpera[];
    created_at: string;
    updated_at: string;
}

export interface ArtistListResponse {
    list: ArtistListItem[];
    pagination: {
        total: number;
        page: number;
        page_size: number;
    }
}

// 搜索结果类型
export interface OperaSearchResult {
    opera_id: number;
    opera_title: string;
    artist_names: string[];
    artist_ids: number[];
}

export interface ArtistSearchResult {
    artist_id: number;
    name: string;
}

export interface User {
    user_id: number;
    username: string;
    phone: string;
    email?: string;
    sex: string;
    icon?: string;
    role: string;
    last_login_at?: string;
    last_login_ip?: string;
    created_at: string;
    updated_at: string;
}
export interface UpdateUserProfileRequest {
    username: string;
    phone: string;
    sex: string;
    email?: string;
    code?: string;
    icon?: string;
}
export interface UpdatePasswordRequest {
    old_password: string;
    new_password: string;
}

export interface Comment {
    comment_id: number;
    user_id: number;
    username: string;
    user_icon?: string;
    opera_id: number;
    parent_comment_id?: number;
    comment_text: string;
    like_count: number;
    liked: boolean;
    created_at: string;
    replies?: Comment[];
    showReplyInput?: boolean;
    replyText?: string;
    showReplyEmojiPicker?: boolean;
}

export interface NewsListItem {
    news_id: number;
    title: string;
    summary: string;
    cover?: string;
    source?: string;
    author?: string;
    is_published: boolean;
    published_at?: string;
    created_at: string;
    updated_at: string;
}

export interface NewsDetail extends NewsListItem {
    content: string;
}

export interface EducationBookListItem {
    book_id: number;
    title: string;
    description: string;
    pdf_path: string;
    pdf_url: string;
    cover_path: string;
    cover_url?: string;
    is_published: boolean;
    created_at: string;
    updated_at: string;
}

export interface EducationBookDetail extends EducationBookListItem {}
