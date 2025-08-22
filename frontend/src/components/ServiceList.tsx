import React, { useEffect, useState } from 'react';
import { Service } from '../types';
import { serviceApi } from '../api/services';
import { 
  PencilIcon, 
  EyeIcon, 
  TrashIcon, 
  PowerIcon,
  PlusIcon,
  ServerIcon
} from '@heroicons/react/outline';

interface ServiceListProps {
  onEdit: (serviceId: number) => void;
  onViewDetails: (service: Service) => void;
  onNewService: () => void;
  refreshTrigger?: number;
}

const ServiceList: React.FC<ServiceListProps> = ({ onEdit, onViewDetails, onNewService, refreshTrigger }) => {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadServices();
  }, [refreshTrigger]);

  const loadServices = async () => {
    setLoading(true);
    try {
      const data = await serviceApi.getServices();
      setServices(data.services || []);
    } catch (error) {
      console.error('Failed to load services:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleStatus = async (serviceId: number) => {
    try {
      await serviceApi.toggleServiceStatus(serviceId);
      await loadServices();
    } catch (error) {
      console.error('Failed to update service status:', error);
    }
  };

  const handleDelete = async (serviceId: number, serviceName: string) => {
    if (!window.confirm(`Are you sure you want to delete "${serviceName}"?`)) {
      return;
    }
    
    try {
      await serviceApi.deleteService(serviceId);
      await loadServices();
    } catch (error) {
      console.error('Failed to delete service:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold text-gray-800">Services</h2>
        <button
          onClick={onNewService}
          className="btn btn-primary flex items-center space-x-2"
        >
          <PlusIcon className="h-5 w-5" />
          <span>New Service</span>
        </button>
      </div>

      {services.length === 0 ? (
        <div className="card text-center py-12">
          <ServerIcon className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-gray-700 mb-2">No Services Yet</h3>
          <p className="text-gray-500 mb-6">Create your first microservice to get started</p>
          <button
            onClick={onNewService}
            className="btn btn-primary inline-flex items-center space-x-2"
          >
            <PlusIcon className="h-5 w-5" />
            <span>Create First Service</span>
          </button>
        </div>
      ) : (
        <div className="grid gap-4">
          {services.map((service) => (
            <div
              key={service.id}
              className="card hover:shadow-2xl transition-all duration-300 border border-gray-100"
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center space-x-3 mb-2">
                    <h3 className="text-xl font-bold text-gray-800">{service.name}</h3>
                    <span className={`status-badge ${service.status ? 'status-enabled' : 'status-disabled'}`}>
                      {service.status ? 'ENABLED' : 'DISABLED'}
                    </span>
                  </div>
                  <p className="text-gray-600 mb-3">{service.description || 'No description'}</p>
                  <div className="flex space-x-6 text-sm text-gray-500">
                    <div>
                      <span className="font-semibold">Port:</span> {service.port}
                    </div>
                    <div>
                      <span className="font-semibold">Package:</span> {service.package}
                    </div>
                    <div>
                      <span className="font-semibold">Created:</span>{' '}
                      {new Date(service.created_at).toLocaleDateString()}
                    </div>
                  </div>
                </div>
                
                <div className="flex space-x-2 ml-4">
                  <button
                    onClick={() => onEdit(service.id)}
                    className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <PencilIcon className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => onViewDetails(service)}
                    className="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors"
                    title="View Details"
                  >
                    <EyeIcon className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => handleToggleStatus(service.id)}
                    className={`p-2 ${
                      service.status ? 'text-amber-600 hover:bg-amber-50' : 'text-green-600 hover:bg-green-50'
                    } rounded-lg transition-colors`}
                    title={service.status ? 'Disable' : 'Enable'}
                  >
                    <PowerIcon className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(service.id, service.name)}
                    className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="Delete"
                  >
                    <TrashIcon className="h-5 w-5" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ServiceList;