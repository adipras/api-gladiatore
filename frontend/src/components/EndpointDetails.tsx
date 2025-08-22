import React, { useEffect, useState } from 'react';
import { Service, Endpoint } from '../types';
import { serviceApi } from '../api/services';
import { 
  XMarkIcon, 
  PowerIcon,
  GlobeAltIcon,
  LockClosedIcon,
  LockOpenIcon
} from '@heroicons/react/outline';

interface EndpointDetailsProps {
  service: Service | null;
  visible: boolean;
  onClose: () => void;
}

const EndpointDetails: React.FC<EndpointDetailsProps> = ({ service, visible, onClose }) => {
  const [endpoints, setEndpoints] = useState<Endpoint[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (visible && service) {
      loadEndpoints();
    }
  }, [visible, service]);

  const loadEndpoints = async () => {
    if (!service) return;
    
    setLoading(true);
    try {
      const data = await serviceApi.getServiceEndpoints(service.id);
      setEndpoints(data.endpoints || []);
    } catch (error) {
      console.error('Failed to load endpoints:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleEndpointStatus = async (endpointId: number) => {
    try {
      await serviceApi.toggleEndpointStatus(endpointId);
      await loadEndpoints();
    } catch (error) {
      console.error('Failed to update endpoint status:', error);
    }
  };

  if (!visible || !service) return null;

  const getMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      GET: 'bg-blue-500',
      POST: 'bg-green-500',
      PUT: 'bg-amber-500',
      DELETE: 'bg-red-500',
      PATCH: 'bg-purple-500'
    };
    return colors[method] || 'bg-gray-500';
  };

  const getOperationIcon = (operation: string) => {
    const icons: Record<string, string> = {
      create: '➕',
      read: '👁',
      update: '✏️',
      delete: '🗑',
      list: '📋',
      custom: '⚡'
    };
    return icons[operation] || '📌';
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl max-w-5xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-green-600 to-teal-600 p-6 text-white">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <GlobeAltIcon className="h-8 w-8" />
              <div>
                <h2 className="text-2xl font-bold">{service.name} - Endpoints</h2>
                <p className="text-green-100 text-sm mt-1">
                  Service running on port {service.port}
                </p>
              </div>
            </div>
            <button
              onClick={onClose}
              className="p-2 hover:bg-white/20 rounded-lg transition-colors"
            >
              <XMarkIcon className="h-6 w-6" />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto" style={{ maxHeight: 'calc(90vh - 100px)' }}>
          {loading ? (
            <div className="flex justify-center items-center h-32">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600"></div>
            </div>
          ) : endpoints.length === 0 ? (
            <div className="text-center py-12">
              <GlobeAltIcon className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">No endpoints found for this service</p>
            </div>
          ) : (
            <div className="space-y-4">
              {endpoints.map((endpoint) => (
                <div
                  key={endpoint.id}
                  className="bg-white border border-gray-200 rounded-xl p-5 hover:shadow-lg transition-all duration-200"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-3">
                        <span className={`px-3 py-1 text-white text-xs font-bold rounded ${getMethodColor(endpoint.method)}`}>
                          {endpoint.method}
                        </span>
                        <code className="text-lg font-mono text-gray-800">{endpoint.path}</code>
                        <span className={`status-badge ${endpoint.status ? 'status-enabled' : 'status-disabled'}`}>
                          {endpoint.status ? 'ENABLED' : 'DISABLED'}
                        </span>
                      </div>
                      
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <span className="text-gray-500">Handler:</span>
                          <p className="font-semibold text-gray-700">{endpoint.handler}</p>
                        </div>
                        <div>
                          <span className="text-gray-500">Operation:</span>
                          <p className="font-semibold text-gray-700">
                            <span className="mr-1">{getOperationIcon(endpoint.operation)}</span>
                            {endpoint.operation}
                          </p>
                        </div>
                        <div>
                          <span className="text-gray-500">Table:</span>
                          <p className="font-semibold text-gray-700">{endpoint.table || '-'}</p>
                        </div>
                        <div>
                          <span className="text-gray-500">Auth:</span>
                          <p className="font-semibold text-gray-700 flex items-center">
                            {endpoint.auth ? (
                              <>
                                <LockClosedIcon className="h-4 w-4 text-green-600 mr-1" />
                                Required
                              </>
                            ) : (
                              <>
                                <LockOpenIcon className="h-4 w-4 text-gray-400 mr-1" />
                                None
                              </>
                            )}
                          </p>
                        </div>
                      </div>
                    </div>
                    
                    <button
                      onClick={() => handleToggleEndpointStatus(endpoint.id)}
                      className={`ml-4 p-2 ${
                        endpoint.status 
                          ? 'text-amber-600 hover:bg-amber-50' 
                          : 'text-green-600 hover:bg-green-50'
                      } rounded-lg transition-colors`}
                      title={endpoint.status ? 'Disable' : 'Enable'}
                    >
                      <PowerIcon className="h-5 w-5" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default EndpointDetails;