import { apiClient } from '@/lib/api';

describe('apiClient', () => {
  beforeEach(() => {
    localStorage.clear();
    jest.clearAllMocks();
  });

  it('should add authorization header with token', async () => {
    localStorage.setItem('token', 'test-token-123');

    // Mock fetch
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ data: { id: 1, name: 'Test' } }),
      } as Response)
    );

    await apiClient.get('/test');

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/test'),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer test-token-123',
        }),
      })
    );
  });

  it('should throw error on non-ok response', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: false,
        status: 404,
        json: () => Promise.resolve({ error: { message: 'Not found' } }),
      } as Response)
    );

    await expect(apiClient.get('/notfound')).rejects.toThrow('Not found');
  });

  it('should handle network errors', async () => {
    global.fetch = jest.fn(() =>
      Promise.reject(new Error('Network error'))
    );

    await expect(apiClient.get('/test')).rejects.toThrow('Network error');
  });
});
