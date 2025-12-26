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