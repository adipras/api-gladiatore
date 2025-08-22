import React, { useState, useEffect } from 'react';
import { ServiceConfig } from '../types';
import { serviceApi } from '../api/services';
import { 
  XMarkIcon, 
  DocumentTextIcon, 
  CheckCircleIcon, 
  PlayIcon,
  CodeBracketIcon
} from '@heroicons/react/outline';

interface ServiceFormProps {
  visible: boolean;
  serviceId?: number | null;
  onClose: () => void;
  onSuccess: () => void;
}

const ServiceForm: React.FC<ServiceFormProps> = ({ visible, serviceId, onClose, onSuccess }) => {
  const [jsonContent, setJsonContent] = useState<string>('');
  const [evaluating, setEvaluating] = useState(false);
  const [validating, setValidating] = useState(false);
  const [validationResult, setValidationResult] = useState<{ valid: boolean; error?: string } | null>(null);
  const [evaluationResult, setEvaluationResult] = useState<any>(null);
  const [jsonError, setJsonError] = useState<string>('');

  useEffect(() => {
    if (visible) {
      if (serviceId) {
        loadServiceConfig();
      } else {
        loadExample();
      }
      setValidationResult(null);
      setEvaluationResult(null);
      setJsonError('');
    }
  }, [visible, serviceId]);

  const loadServiceConfig = async () => {
    if (!serviceId) return;
    
    try {
      const configData = await serviceApi.getServiceConfig(serviceId);
      setJsonContent(JSON.stringify(configData, null, 2));
    } catch (error) {
      console.error('Failed to load service configuration:', error);
    }
  };

  const loadExample = async () => {
    try {
      const example = await serviceApi.getExample();
      setJsonContent(JSON.stringify(example, null, 2));
    } catch (error) {
      console.error('Failed to load example configuration:', error);
    }
  };

  const validateJson = async () => {
    try {
      const configObj = JSON.parse(jsonContent);
      setJsonError('');
      
      setValidating(true);
      const result = await serviceApi.validateConfig(configObj);
      setValidationResult(result);
      
      return result.valid;
    } catch (error: any) {
      if (error instanceof SyntaxError) {
        setJsonError('Invalid JSON syntax');
        setValidationResult({ valid: false, error: 'Invalid JSON syntax' });
      }
      return false;
    } finally {
      setValidating(false);
    }
  };

  const evaluateJson = async () => {
    // First validate the JSON
    const isValid = await validateJson();
    if (!isValid) {
      return;
    }

    try {
      const configObj = JSON.parse(jsonContent) as ServiceConfig;
      setEvaluating(true);
      
      let result;
      if (serviceId) {
        result = await serviceApi.updateService(serviceId, configObj);
      } else {
        result = await serviceApi.createService(configObj);
      }
      
      setEvaluationResult(result);
      
      // Show success for 2 seconds before closing
      setTimeout(() => {
        onSuccess();
        onClose();
      }, 2000);
    } catch (error: any) {
      console.error('Evaluation failed:', error);
      setEvaluationResult({
        error: error.response?.data?.error || 'Evaluation failed'
      });
    } finally {
      setEvaluating(false);
    }
  };

  if (!visible) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl max-w-6xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-primary-600 to-secondary-600 p-6 text-white">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <CodeBracketIcon className="h-8 w-8" />
              <div>
                <h2 className="text-2xl font-bold">
                  {serviceId ? 'Edit Service Configuration' : 'New Service Configuration'}
                </h2>
                <p className="text-primary-100 text-sm mt-1">
                  Define your microservice configuration in JSON format
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

        <div className="flex h-[calc(90vh-180px)]">
          {/* JSON Editor Panel */}
          <div className="flex-1 p-6 border-r border-gray-200">
            <div className="mb-4 flex items-center justify-between">
              <label className="text-lg font-semibold text-gray-700">
                JSON Configuration
              </label>
              <button
                onClick={loadExample}
                className="text-sm text-primary-600 hover:text-primary-700 flex items-center space-x-1"
              >
                <DocumentTextIcon className="h-4 w-4" />
                <span>Load Example</span>
              </button>
            </div>
            
            <textarea
              value={jsonContent}
              onChange={(e) => {
                setJsonContent(e.target.value);
                setJsonError('');
                setValidationResult(null);
                setEvaluationResult(null);
              }}
              className="w-full h-[calc(100%-60px)] p-4 bg-gray-900 text-green-400 rounded-lg font-mono text-sm resize-none focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder={`{
  "service": {
    "name": "my-service",
    "description": "My microservice",
    "package": "github.com/example/my-service",
    "port": 8080
  },
  ...
}`}
              spellCheck={false}
            />
          </div>

          {/* Results Panel */}
          <div className="flex-1 p-6 bg-gray-50">
            <div className="mb-4">
              <h3 className="text-lg font-semibold text-gray-700 mb-4">Evaluation Results</h3>
              
              {/* Validation Status */}
              {validationResult && (
                <div className={`mb-4 p-4 rounded-lg ${
                  validationResult.valid 
                    ? 'bg-green-50 border border-green-200' 
                    : 'bg-red-50 border border-red-200'
                }`}>
                  <div className="flex items-start">
                    <CheckCircleIcon className={`h-5 w-5 mr-2 mt-0.5 ${
                      validationResult.valid ? 'text-green-600' : 'text-red-600'
                    }`} />
                    <div className="flex-1">
                      <p className={`font-semibold ${
                        validationResult.valid ? 'text-green-800' : 'text-red-800'
                      }`}>
                        {validationResult.valid ? '✓ Valid Configuration' : '✗ Invalid Configuration'}
                      </p>
                      {validationResult.error && (
                        <p className="text-sm text-red-600 mt-1">{validationResult.error}</p>
                      )}
                    </div>
                  </div>
                </div>
              )}

              {/* JSON Syntax Error */}
              {jsonError && (
                <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-red-800 font-semibold">Syntax Error</p>
                  <p className="text-sm text-red-600 mt-1">{jsonError}</p>
                </div>
              )}

              {/* Evaluation Results */}
              {evaluationResult && (
                <div className={`p-4 rounded-lg ${
                  evaluationResult.error 
                    ? 'bg-red-50 border border-red-200' 
                    : 'bg-blue-50 border border-blue-200'
                }`}>
                  {evaluationResult.error ? (
                    <div>
                      <p className="text-red-800 font-semibold">Evaluation Failed</p>
                      <p className="text-sm text-red-600 mt-1">{evaluationResult.error}</p>
                    </div>
                  ) : (
                    <div>
                      <p className="text-blue-800 font-semibold mb-2">
                        ✓ Service Generated Successfully!
                      </p>
                      {evaluationResult.result && (
                        <div className="space-y-2 text-sm">
                          <p><span className="font-semibold">Service:</span> {evaluationResult.result.service_name}</p>
                          <p><span className="font-semibold">Output:</span> {evaluationResult.result.output_path}</p>
                          <div>
                            <span className="font-semibold">Files Generated:</span>
                            <ul className="mt-1 ml-4 list-disc">
                              {evaluationResult.result.files?.map((file: string, idx: number) => (
                                <li key={idx} className="text-blue-600">{file}</li>
                              ))}
                            </ul>
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              )}

              {/* Instructions when no results */}
              {!validationResult && !evaluationResult && !jsonError && (
                <div className="bg-gray-100 rounded-lg p-6 text-center">
                  <CodeBracketIcon className="h-12 w-12 text-gray-400 mx-auto mb-3" />
                  <p className="text-gray-600">
                    Enter your service configuration JSON and click "Evaluate" to generate your microservice
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="flex justify-between items-center p-6 border-t border-gray-200 bg-gray-50">
          <button
            onClick={validateJson}
            disabled={validating || !jsonContent}
            className="btn btn-secondary flex items-center space-x-2"
          >
            <CheckCircleIcon className="h-5 w-5" />
            <span>{validating ? 'Validating...' : 'Validate JSON'}</span>
          </button>
          
          <div className="flex space-x-3">
            <button
              onClick={onClose}
              className="btn btn-secondary"
            >
              Cancel
            </button>
            <button
              onClick={evaluateJson}
              disabled={evaluating || !jsonContent}
              className="btn btn-primary flex items-center space-x-2"
            >
              <PlayIcon className="h-5 w-5" />
              <span>{evaluating ? 'Evaluating...' : 'Evaluate'}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ServiceForm;