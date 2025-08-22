import axios from 'axios';
import { Service, ServiceConfig, Endpoint, GenerationResult } from '../types';

const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const serviceApi = {
  // Get all services
  async getServices(): Promise<{ services: Service[]; total: number }> {
    const response = await api.get('/services');
    return response.data;
  },

  // Get service by ID
  async getService(id: number): Promise<Service> {
    const response = await api.get(`/services/${id}`);
    return response.data;
  },

  // Get service configuration
  async getServiceConfig(id: number): Promise<ServiceConfig> {
    const response = await api.get(`/services/${id}/config`);
    return response.data;
  },

  // Create new service
  async createService(config: ServiceConfig): Promise<{ service: Service; result: GenerationResult }> {
    const response = await api.post('/services', config);
    return response.data;
  },

  // Update service
  async updateService(id: number, config: ServiceConfig): Promise<{ service: Service; result: GenerationResult }> {
    const response = await api.put(`/services/${id}`, config);
    return response.data;
  },

  // Toggle service status
  async toggleServiceStatus(id: number): Promise<{ message: string; status: boolean }> {
    const response = await api.patch(`/services/${id}/status`);
    return response.data;
  },

  // Delete service
  async deleteService(id: number): Promise<{ message: string }> {
    const response = await api.delete(`/services/${id}`);
    return response.data;
  },

  // Get service endpoints
  async getServiceEndpoints(serviceId: number): Promise<{ endpoints: Endpoint[]; total: number }> {
    const response = await api.get(`/services/${serviceId}/endpoints`);
    return response.data;
  },

  // Toggle endpoint status
  async toggleEndpointStatus(endpointId: number): Promise<{ message: string; status: boolean }> {
    const response = await api.patch(`/endpoints/${endpointId}/status`);
    return response.data;
  },

  // Get example configuration
  async getExample(): Promise<ServiceConfig> {
    const response = await api.get('/example');
    return response.data;
  },

  // Validate configuration
  async validateConfig(config: ServiceConfig): Promise<{ valid: boolean; error?: string; config?: ServiceConfig }> {
    const response = await api.post('/validate', config);
    return response.data;
  },
};