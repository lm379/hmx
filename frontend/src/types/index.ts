export interface Opera {
    opera_id: number;
    opera_title: string;
    artists: { Name: string }[];
    release_date?: string;
    duration?: string;
    music_path?: string;
    video_path: string;
    srt_path?: string;
    description: string;
    avatar?: string;
    ai_summary: string;
    created_at: string;
    updated_at: string;
    like_count?: number;
    favorite_count?: number;
    share_count?: number;
    liked?: boolean;
    favorited?: boolean;
    play_count?: number;
}

export interface OperaListResponse {
    list: Opera[];
    pagination: {
        total: number;
        page: number;
        page_size: number;
    }
}

export interface Artist {
    artist_id: number;
    name: string;
    bio: string;
    avatar?: string;
    operas?: Opera[];
    created_at: string;
    updated_at: string;
}

export interface ArtistListResponse {
    list: Artist[];
    pagination: {
        total: number;
        page: number;
        page_size: number;
    }
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
}
export interface UpdatePasswordRequest {
    old_password: string;
    new_password: string;
}