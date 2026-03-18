import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'

// Mock axios
vi.mock('axios', () => ({
    default: {
        get: vi.fn(),
        post: vi.fn(),
        interceptors: {
            request: { use: vi.fn() },
            response: { use: vi.fn() },
        },
        defaults: {
            headers: {
                common: {} as Record<string, string>,
            },
        },
    },
}))

// Mock localStorage
const localStorageMock = {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
})

describe('useAuthStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('should initialize with default values', () => {
        const store = useAuthStore()
        expect(store.isLoggedIn).toBe(false)
        expect(store.user).toBeNull()
        expect(store.loading).toBe(false)
    })

    it('should check login status when no token exists', async () => {
        localStorageMock.getItem.mockReturnValue(null)

        const store = useAuthStore()
        await store.checkLoginStatus()

        expect(store.isLoggedIn).toBe(false)
        expect(store.user).toBeNull()
    })

    it('should login and store tokens', async () => {
        const mockUser = {
            user_id: 1,
            username: 'testuser',
            email: 'test@example.com',
            role: 'User',
        }

        localStorageMock.getItem.mockReturnValue('test-access-token')

        const axios = await import('axios')
        vi.mocked(axios.default.get).mockResolvedValueOnce({
            data: { data: mockUser },
        })

        const store = useAuthStore()
        await store.login('test-access-token', 'test-refresh-token')

        expect(localStorageMock.setItem).toHaveBeenCalledWith('access_token', 'test-access-token')
        expect(localStorageMock.setItem).toHaveBeenCalledWith('refresh_token', 'test-refresh-token')
        expect(store.isLoggedIn).toBe(true)
        expect(store.user).toEqual(mockUser)
    })

    it('should logout and clear tokens', async () => {
        const axios = await import('axios')
        vi.mocked(axios.default.post).mockResolvedValueOnce({})

        const store = useAuthStore()
        store.isLoggedIn = true
        store.user = { 
            user_id: 1, 
            username: 'test', 
            phone: '13800138000',
            email: 'test@test.com', 
            sex: 'male',
            role: 'User',
            created_at: '2024-01-01',
            updated_at: '2024-01-01'
        }

        await store.logout()

        expect(localStorageMock.removeItem).toHaveBeenCalledWith('access_token')
        expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
        expect(store.isLoggedIn).toBe(false)
        expect(store.user).toBeNull()
    })

    it('should handle logout errors gracefully', async () => {
        const axios = await import('axios')
        vi.mocked(axios.default.post).mockRejectedValueOnce(new Error('Logout failed'))

        const store = useAuthStore()
        store.isLoggedIn = true

        await store.logout()

        // Should still clear local state even if API call fails
        expect(store.isLoggedIn).toBe(false)
        expect(localStorageMock.removeItem).toHaveBeenCalledWith('access_token')
        expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token')
    })

    it('should handle check login status errors', async () => {
        localStorageMock.getItem.mockReturnValue('test-token')

        const axios = await import('axios')
        vi.mocked(axios.default.get).mockRejectedValueOnce(new Error('Unauthorized'))

        const store = useAuthStore()
        await store.checkLoginStatus()

        expect(store.isLoggedIn).toBe(false)
        expect(store.user).toBeNull()
    })
})
