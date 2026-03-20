import axios, { AxiosInstance } from 'axios';
import {
  Tenant,
  Store,
  Product,
  Category,
  Collection,
  Tag,
  ProductImage,
  LoginRequest,
  LoginResponse,
  CreateTenantRequest,
  CreateStoreRequest,
  UpdateStoreRequest,
  CreateProductRequest,
  CreateTagRequest,
  CreateCategoryRequest,
  CreateCollectionRequest,
  PaginatedResponse,
  ProductFilters
} from '@/lib/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add request interceptor to include auth token
    this.client.interceptors.request.use((config) => {
      const token = this.getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token expired or invalid
          this.clearToken();
          window.location.href = '/auth/login';
        }
        return Promise.reject(error);
      }
    );
  }

  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  }

  setToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  }

  private clearToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  }

  // Auth endpoints
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>('/api/auth/tenant/login', data);
    this.setToken(response.data.token);
    return response.data;
  }

  async logout(): Promise<void> {
    this.clearToken();
  }

  async createTenant(data: CreateTenantRequest): Promise<Tenant> {
    const response = await this.client.post<Tenant>('/api/tenants', data);
    return response.data;
  }

  // Store endpoints
  async createStore(data: CreateStoreRequest): Promise<Store> {
    const response = await this.client.post<Store>('/api/stores', data);
    return response.data;
  }

  async getStores(): Promise<Store[]> {
    const response = await this.client.get<Store[]>('/api/stores');
    return response.data;
  }

  async getStore(id: string): Promise<Store> {
    const response = await this.client.get<Store>(`/api/stores/${id}`);
    return response.data;
  }

  async updateStore(id: string, data: UpdateStoreRequest): Promise<Store> {
    const response = await this.client.put<Store>(`/api/stores/${id}`, data);
    return response.data;
  }

  async publishStoreCustomization(id: string, useDraftLayout = true): Promise<Store> {
    const response = await this.client.post<Store>(`/api/stores/${id}/customization/publish`, {
      use_draft_layout: useDraftLayout,
    });
    return response.data;
  }

  async deleteStore(id: string): Promise<void> {
    await this.client.delete(`/api/stores/${id}`);
  }

  // Product endpoints
  async createProduct(storeId: string, data: CreateProductRequest): Promise<Product> {
    const response = await this.client.post<Product>(`/api/stores/${storeId}/products`, data);
    return response.data;
  }

  async getProducts(storeId: string, filters?: ProductFilters): Promise<PaginatedResponse<Product>> {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, value.toString());
        }
      });
    }

    const response = await this.client.get<PaginatedResponse<Product>>(
      `/api/stores/${storeId}/products?${params.toString()}`
    );
    return response.data;
  }

  async getProduct(storeId: string, productId: string): Promise<Product> {
    const response = await this.client.get<Product>(`/api/stores/${storeId}/products/${productId}`);
    return response.data;
  }

  async updateProduct(storeId: string, productId: string, data: Partial<CreateProductRequest>): Promise<Product> {
    const response = await this.client.put<Product>(`/api/stores/${storeId}/products/${productId}`, data);
    return response.data;
  }

  async deleteProduct(storeId: string, productId: string): Promise<void> {
    await this.client.delete(`/api/stores/${storeId}/products/${productId}`);
  }

  // Category endpoints
  async createCategory(storeId: string, data: CreateCategoryRequest): Promise<Category> {
    const response = await this.client.post<Category>(`/api/stores/${storeId}/categories`, data);
    return response.data;
  }

  async getCategories(storeId: string): Promise<Category[]> {
    const response = await this.client.get<Category[]>(`/api/stores/${storeId}/categories`);
    return response.data;
  }

  async getCategory(storeId: string, categoryId: string): Promise<Category> {
    const response = await this.client.get<Category>(`/api/stores/${storeId}/categories/${categoryId}`);
    return response.data;
  }

  async updateCategory(storeId: string, categoryId: string, data: Partial<CreateCategoryRequest>): Promise<Category> {
    const response = await this.client.put<Category>(`/api/stores/${storeId}/categories/${categoryId}`, data);
    return response.data;
  }

  async deleteCategory(storeId: string, categoryId: string): Promise<void> {
    await this.client.delete(`/api/stores/${storeId}/categories/${categoryId}`);
  }

  // Collection endpoints
  async createCollection(storeId: string, data: CreateCollectionRequest): Promise<Collection> {
    const response = await this.client.post<Collection>(`/api/stores/${storeId}/collections`, data);
    return response.data;
  }

  async getCollections(storeId: string): Promise<Collection[]> {
    const response = await this.client.get<Collection[]>(`/api/stores/${storeId}/collections`);
    return response.data;
  }

  async getCollection(storeId: string, collectionId: string): Promise<Collection> {
    const response = await this.client.get<Collection>(`/api/stores/${storeId}/collections/${collectionId}`);
    return response.data;
  }

  async updateCollection(storeId: string, collectionId: string, data: Partial<CreateCollectionRequest>): Promise<Collection> {
    const response = await this.client.put<Collection>(`/api/stores/${storeId}/collections/${collectionId}`, data);
    return response.data;
  }

  async deleteCollection(storeId: string, collectionId: string): Promise<void> {
    await this.client.delete(`/api/stores/${storeId}/collections/${collectionId}`);
  }

  // Tag endpoints
  async createTag(storeId: string, data: CreateTagRequest): Promise<Tag> {
    const response = await this.client.post<Tag>(`/api/stores/${storeId}/tags`, data);
    return response.data;
  }

  async getTags(storeId: string): Promise<Tag[]> {
    const response = await this.client.get<Tag[]>(`/api/stores/${storeId}/tags`);
    return response.data;
  }

  async getTag(storeId: string, tagId: string): Promise<Tag> {
    const response = await this.client.get<Tag>(`/api/stores/${storeId}/tags/${tagId}`);
    return response.data;
  }

  async updateTag(storeId: string, tagId: string, data: Partial<CreateTagRequest>): Promise<Tag> {
    const response = await this.client.put<Tag>(`/api/stores/${storeId}/tags/${tagId}`, data);
    return response.data;
  }

  async deleteTag(storeId: string, tagId: string): Promise<void> {
    await this.client.delete(`/api/stores/${storeId}/tags/${tagId}`);
  }

  // Product Image endpoints
  async getProductImages(storeId: string, productId: string): Promise<ProductImage[]> {
    const response = await this.client.get<ProductImage[]>(`/api/stores/${storeId}/products/${productId}/images`);
    return response.data;
  }

  async createProductImage(storeId: string, productId: string, data: { url: string; alt_text?: string; position: number; is_featured: boolean }): Promise<ProductImage> {
    const response = await this.client.post<ProductImage>(`/api/stores/${storeId}/products/${productId}/images`, data);
    return response.data;
  }

  async updateProductImage(storeId: string, productId: string, imageId: string, data: Partial<{ url: string; alt_text?: string; position: number; is_featured: boolean }>): Promise<ProductImage> {
    const response = await this.client.put<ProductImage>(`/api/stores/${storeId}/products/${productId}/images/${imageId}`, data);
    return response.data;
  }

  async deleteProductImage(storeId: string, productId: string, imageId: string): Promise<void> {
    await this.client.delete(`/api/stores/${storeId}/products/${productId}/images/${imageId}`);
  }

  async reorderProductImages(storeId: string, productId: string, data: { image_ids: string[] }): Promise<void> {
    await this.client.post(`/api/stores/${storeId}/products/${productId}/images/reorder`, data);
  }
}

export const apiClient = new ApiClient();