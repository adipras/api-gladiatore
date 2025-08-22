import React, { useState } from 'react';
import ServiceList from './components/ServiceList';
import ServiceForm from './components/ServiceForm';
import EndpointDetails from './components/EndpointDetails';
import { Service } from './types';
import './index.css';

const App: React.FC = () => {
  const [showServiceForm, setShowServiceForm] = useState(false);
  const [editingServiceId, setEditingServiceId] = useState<number | null>(null);
  const [showEndpoints, setShowEndpoints] = useState(false);
  const [selectedService, setSelectedService] = useState<Service | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleNewService = () => {
    setEditingServiceId(null);
    setShowServiceForm(true);
  };

  const handleEditService = (serviceId: number) => {
    setEditingServiceId(serviceId);
    setShowServiceForm(true);
  };

  const handleViewDetails = (service: Service) => {
    setSelectedService(service);
    setShowEndpoints(true);
  };

  const handleServiceFormClose = () => {
    setShowServiceForm(false);
    setEditingServiceId(null);
  };

  const handleServiceFormSuccess = () => {
    setShowServiceForm(false);
    setEditingServiceId(null);
    setRefreshTrigger(prev => prev + 1);
  };

  const handleEndpointsClose = () => {
    setShowEndpoints(false);
    setSelectedService(null);
  };

  return (
    <div className="min-h-screen gradient-bg">
      {/* Header */}
      <header className="glass-effect border-b border-white/20">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="text-4xl animate-float">⚔️</div>
              <div>
                <h1 className="text-3xl font-bold text-white">API Gladiatore</h1>
                <p className="text-white/80 text-sm">Microservice Generation & Management Platform</p>
              </div>
            </div>
            <div className="text-white/60 text-sm">
              Generate powerful Go microservices from JSON
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 py-8">
        <ServiceList
          onEdit={handleEditService}
          onViewDetails={handleViewDetails}
          onNewService={handleNewService}
          refreshTrigger={refreshTrigger}
        />
      </main>

      {/* Service Form Modal */}
      <ServiceForm
        visible={showServiceForm}
        serviceId={editingServiceId}
        onClose={handleServiceFormClose}
        onSuccess={handleServiceFormSuccess}
      />

      {/* Endpoint Details Modal */}
      <EndpointDetails
        service={selectedService}
        visible={showEndpoints}
        onClose={handleEndpointsClose}
      />
    </div>
  );
};

export default App;
