import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import { getArtistById, getArtistsByIds } from './artist'

// Mock axios
vi.mock('axios')

describe('artist API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getArtistById', () => {
    it('should fetch artist by id', async () => {
      const mockArtist = {
        artist_id: 1,
        name: 'Test Artist',
        bio: 'Test Bio',
        created_at: '2024-01-01',
      }

      vi.mocked(axios.get).mockResolvedValueOnce({
        data: {
          code: 200,
          data: mockArtist,
        },
      })

      const result = await getArtistById(1)

      expect(axios.get).toHaveBeenCalledWith('/api/v1/artists/1')
      expect(result.code).toBe(200)
      expect(result.data).toEqual(mockArtist)
    })

    it('should handle errors', async () => {
      vi.mocked(axios.get).mockRejectedValueOnce(new Error('Network error'))

      await expect(getArtistById(1)).rejects.toThrow('Network error')
    })
  })

  describe('getArtistsByIds', () => {
    it('should fetch multiple artists', async () => {
      const mockArtists = [
        { artist_id: 1, name: 'Artist 1', bio: 'Bio 1', created_at: '2024-01-01' },
        { artist_id: 2, name: 'Artist 2', bio: 'Bio 2', created_at: '2024-01-02' },
      ]

      vi.mocked(axios.get)
        .mockResolvedValueOnce({
          data: { code: 200, data: mockArtists[0] },
        })
        .mockResolvedValueOnce({
          data: { code: 200, data: mockArtists[1] },
        })

      const result = await getArtistsByIds([1, 2])

      expect(result).toHaveLength(2)
      expect(result[0]!.name).toBe('Artist 1')
      expect(result[1]!.name).toBe('Artist 2')
    })

    it('should handle partial failures', async () => {
      vi.mocked(axios.get)
        .mockResolvedValueOnce({
          data: {
            code: 200,
            data: { artist_id: 1, name: 'Artist 1', bio: 'Bio 1', created_at: '2024-01-01' },
          },
        })
        .mockRejectedValueOnce(new Error('Not found'))

      const result = await getArtistsByIds([1, 2])

      expect(result).toHaveLength(1)
      expect(result[0]!.artist_id).toBe(1)
    })

    it('should handle empty array', async () => {
      const result = await getArtistsByIds([])
      expect(result).toHaveLength(0)
      expect(axios.get).not.toHaveBeenCalled()
    })

    it('should filter out failed requests', async () => {
      vi.mocked(axios.get)
        .mockResolvedValueOnce({
          data: { code: 404, data: null },
        })
        .mockResolvedValueOnce({
          data: {
            code: 200,
            data: { artist_id: 2, name: 'Artist 2', bio: 'Bio 2', created_at: '2024-01-02' },
          },
        })

      const result = await getArtistsByIds([1, 2])

      expect(result).toHaveLength(1)
      expect(result[0]!.artist_id).toBe(2)
    })
  })
})
